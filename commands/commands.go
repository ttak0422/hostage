package commands

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/urfave/cli/v3"

	"github.com/ttak0422/hostage/config"
)

const (
	HOSTS_PATH       = "/etc/hosts"
	CONFIG_FILE_NAME = "config.dhall"
)

func getConfig() (*config.Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	cfgPath := filepath.Join(homeDir, ".config", "hostage", CONFIG_FILE_NAME)
	return config.LoadConfig(cfgPath)
}

func ShowConfigKey(ctx context.Context, cmd *cli.Command) error {
	cfg, err := getConfig()
	if err != nil {
		return err
	}

	for _, key := range cfg.GetGroupKeys() {
		fmt.Println(key)
	}

	return nil
}

func ShowConfig(ctx context.Context, cmd *cli.Command) error {
	if err := ensureDir(HOSTS_PATH); err != nil {
		return err
	}

	key := cmd.Args().Get(0)
	if key == "" {
		return errors.New("key is required")
	}

	cfg, err := getConfig()
	if err != nil {
		return err
	}

	hg, err := cfg.GetGroupByName(key)
	if err != nil {
		return err
	}

	fmt.Println(hg.ToFormatedText())

	return nil
}

func SetupConfig(ctx context.Context, cmd *cli.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(homeDir, ".config", "hostage", CONFIG_FILE_NAME)

	if err = ensureDir(cfgPath); err != nil {
		return err
	}

	_, err = os.Stat(cfgPath)
	if !os.IsNotExist(err) {
		return errors.New("config file already exists")
	}

	cfg := config.SampleDhall()
	if err := os.WriteFile(cfgPath, []byte(cfg), 0o644); err != nil {
		return err
	}

	fmt.Println("config file created at", cfgPath)

	return nil
}

func OpenWithEditor(ctx context.Context, cmd *cli.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgPath := filepath.Join(homeDir, ".config", "hostage", CONFIG_FILE_NAME)

	_, err = os.Stat(cfgPath)
	if os.IsNotExist(err) {
		return errors.New("config file not found")
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		return errors.New("EDITOR is not set")
	}

	parts := strings.Fields(editor)
	args := append(parts[1:], cfgPath)

	command := exec.Command(parts[0], args...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	return command.Run()
}

func ensureDir(path string) error {
	dir := filepath.Dir(path)
	return os.MkdirAll(dir, os.ModePerm)
}
