package vc

import (
	"reflect"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	var (
		emptyConfig = Config{
			Name:        "My simple vc script",
			Description: "Create some default users and config",
			Script:      nil,
			Finaliser:   Finaliser_Nothing,
		}

		nothingConfig  = emptyConfig
		rebootConfig   = emptyConfig
		shutdownConfig = emptyConfig
	)

	nothingConfig.Finaliser = Finaliser_Nothing
	nothingConfig.Sum = "v1_7be38c4c1866b8c70c2b12b7ec130d8a"

	rebootConfig.Finaliser = Finaliser_Reboot
	rebootConfig.Sum = "v1_7be38c4c1866b8c70c2b12b7ec130d8a"

	shutdownConfig.Finaliser = Finaliser_Shutdown
	shutdownConfig.Sum = "v1_7be38c4c1866b8c70c2b12b7ec130d8a"

	for _, test := range []struct {
		name        string
		fn          string
		expect      Config
		expectError bool
	}{
		{"Happy path", "testdata/simple-vc.toml", Config{
			Name:        "My simple vc script",
			Description: "Create some default users and config",
			Script:      nil,
			Finaliser:   Finaliser_Nothing,
			Sum:         "v1_7be38c4c1866b8c70c2b12b7ec130d8a",
		}, false},

		// Specific cases, failures, errors, etc.
		{"Unknown file fails", "testdata/script-configs/nonsuch.toml", Config{}, true},
		{"Empty config fails", "testdata/script-configs/empty.toml", Config{}, true},
		{"Missing script section fails", "testdata/script-configs/no-script-section.toml", emptyConfig, true},
		{"Missing script interpreter fails", "testdata/script-configs/no-interpreter.toml", emptyConfig, true},
		{"Both url and body fails", "testdata/script-configs/both-url-and-body.toml", emptyConfig, true},
		{"Neither url or body fails", "testdata/script-configs/neither-url-or-body.toml", emptyConfig, true},
		{"Empty finaliser defaults to nothing", "testdata/script-configs/empty-finaliser.toml", nothingConfig, false},
		{"Default finaliser defaults to nothing", "testdata/script-configs/default-finaliser.toml", nothingConfig, false},
		{"Nothing finaliser works accordingly", "testdata/script-configs/nothing-finaliser.toml", nothingConfig, false},
		{"Reboot finaliser works accordingly", "testdata/script-configs/reboot-finaliser.toml", rebootConfig, false},
		{"Shutdown finaliser works accordingly", "testdata/script-configs/shutdown-finaliser.toml", shutdownConfig, false},
		{"Dodgy finaliser fails", "testdata/script-configs/dodgy-finaliser.toml", emptyConfig, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			c, err := LoadConfig(test.fn)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}

			// Avoid validating that this potentially large
			// script has stuff in it
			c.Script = nil

			if !reflect.DeepEqual(test.expect, c) {
				t.Errorf("expected:\n%#v\nreceived:\n%#v", test.expect, c)
			}
		})
	}
}
