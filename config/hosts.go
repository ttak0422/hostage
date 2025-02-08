package config

import (
	"fmt"
	"os"
	"strings"

	_ "embed"

	"github.com/philandstuff/dhall-golang/v6"
)

//go:embed hosts.dhall
var hostsDhall string

type HostEntry struct {
	IP      string
	Aliases []string
}

type HostGroup struct {
	Name    string
	Entries []HostEntry
}

type Config struct {
	Groups []HostGroup
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := dhall.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func (c *Config) GetGroupKeys() []string {
	keys := make([]string, len(c.Groups))
	for i, group := range c.Groups {
		keys[i] = group.Name
	}
	return keys
}

func (c *Config) GetGroupByName(name string) (*HostGroup, error) {
	for _, group := range c.Groups {
		if group.Name == name {
			return &group, nil
		}
	}
	return nil, fmt.Errorf("group not found: %s", name)
}

func (hg *HostGroup) ToFormatedText() string {
	var builder strings.Builder

	builder.WriteString("##\n")
	builder.WriteString("# configured by hostage.\n")
	builder.WriteString("#\n")
	builder.WriteString(fmt.Sprintf("# group: %s\n", hg.Name))
	builder.WriteString("##\n")

	for _, entry := range hg.Entries {
		builder.WriteString(fmt.Sprintf("%s %s\n", entry.IP, strings.Join(entry.Aliases, " ")))
	}

	return builder.String()
}

func SampleDhall() string {
	return hostsDhall
}
