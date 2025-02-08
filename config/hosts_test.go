package config

import (
	"reflect"
	"testing"

	"github.com/philandstuff/dhall-golang/v6"
)

func TestEmbed(t *testing.T) {
	exp := `let Entries
    : Type
    = { IP : Text, Aliases : List Text }

let Group
    : Type
    = { Name : Text, Entries : List Entries }

let Config
    : Type
    = { Groups : List Group }

let config
    : Config
    = { Groups =
        [ { Name = "sample"
          , Entries =
            [ { IP = "127.0.0.1", Aliases = [ "localhost" ] }
            , { IP = "255.255.255.255", Aliases = [ "broadcasthost" ] }
            , { IP = "::1", Aliases = [ "localhost" ] }
            ]
          }
        ]
      }

in  config
`

	act := SampleDhall()

	if act != exp {
		t.Errorf("expected:\n%s\nactual:\n%s\n", act, exp)
	}
}

func TestUnmarshal(t *testing.T) {
	src := SampleDhall()
	exp := Config{
		Groups: []HostGroup{
			{
				Name: "sample",
				Entries: []HostEntry{
					{
						IP:      "127.0.0.1",
						Aliases: []string{"localhost"},
					},
					{
						IP:      "255.255.255.255",
						Aliases: []string{"broadcasthost"},
					},
					{
						IP:      "::1",
						Aliases: []string{"localhost"},
					},
				},
			},
		},
	}

	var act Config
	if err := dhall.Unmarshal([]byte(src), &act); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(act, exp) {
		t.Errorf("expected:\n%v\nactual:\n%v\n", act, exp)
	}
}

func TestGetGroupKeys(t *testing.T) {
	tests := []struct {
		desc   string
		config Config
		exp    []string
	}{
		{
			desc:   "empty",
			config: Config{Groups: []HostGroup{}},
			exp:    []string{},
		},
		{
			desc: "simple",
			config: Config{
				Groups: []HostGroup{
					{
						Name: "foo",
						Entries: []HostEntry{
							{
								IP:      "192.168.1.1",
								Aliases: []string{"hoge"},
							},
						},
					},
				},
			},
			exp: []string{"foo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			t.Parallel()
			act := tt.config.GetGroupKeys()
			if !reflect.DeepEqual(act, tt.exp) {
				t.Error("group keys are not equal")
			}
		})
	}
}

func TestGetGroupByName(t *testing.T) {
	config := Config{
		Groups: []HostGroup{
			{
				Name: "foo",
				Entries: []HostEntry{
					{
						IP:      "192.168.1.1",
						Aliases: []string{"hoge"},
					},
				},
			},
		},
	}
	exp := HostGroup{
		Name: "foo",
		Entries: []HostEntry{
			{
				IP:      "192.168.1.1",
				Aliases: []string{"hoge"},
			},
		},
	}

	act, err := config.GetGroupByName("foo")
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(*act, exp) {
		t.Errorf("expected: %v\nactual %v\n", act, exp)
	}
}

func TestToFormatedText(t *testing.T) {
	src := SampleDhall()
	var config Config
	if err := dhall.Unmarshal([]byte(src), &config); err != nil {
		t.Fatal(err)
	}

	exp := `##
# configured by hostage.
#
# group: sample
##
127.0.0.1 localhost
255.255.255.255 broadcasthost
::1 localhost
`
	act := config.Groups[0].ToFormatedText()

	if act != exp {
		t.Errorf("expected: %s\nactual %s\n", act, exp)
	}
}
