package vc

import (
	"os"
	"reflect"
	"testing"
)

func TestScript_Run(t *testing.T) {
	sh := "/bin/sh"

	for _, test := range []struct {
		name        string
		script      string
		expectError bool
		expectLog   []string
	}{
		{"successful script", "testdata/scripts/success.sh", false, []string{"+ echo 'hello world!'", "hello world!", ""}},
		{"failing script", "testdata/scripts/fail.sh", true, []string{"uh-oh", ""}},
	} {
		t.Run(test.name, func(t *testing.T) {
			body, _ := os.ReadFile(test.script)
			bodyStr := string(body)

			s := Script{
				Interpreter: &sh,
				Body:        &bodyStr,
			}

			s.Run()

			if s.Error == "" && test.expectError {
				t.Error("expected error, received none")
			} else if s.Error != "" && !test.expectError {
				t.Errorf("unexpected error: %+v", s.Error)
			}

			if !reflect.DeepEqual(test.expectLog, s.Logs) {
				t.Errorf("expected %#v, received %#v", test.expectLog, s.Logs)
			}
		})
	}
}
