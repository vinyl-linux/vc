package vc

import (
	"testing"
)

func TestConfig_Finalise(t *testing.T) {
	finaliserNothingCommand = "date"
	finaliserRebootCommand = "echo reboot"
	finaliserShutdownCommand = ""

	for _, test := range []struct {
		name        string
		finaliser   Finaliser
		expectError bool
	}{
		{"nothing command", Finaliser_Nothing, false},
		{"reboot command", Finaliser_Reboot, false},
		{"shutdown command", Finaliser_Shutdown, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			c := Config{
				Finaliser: test.finaliser,
			}

			err := c.Finalise()
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}
