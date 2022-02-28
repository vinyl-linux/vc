package vc

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"
)

var (
	sh = "/bin/sh"
)

type dummyHttpClient struct {
	err    bool
	status int
}

func (d dummyHttpClient) Get(url string) (resp *http.Response, err error) {
	if d.err {
		err = fmt.Errorf("an error")

		return
	}

	resp = &http.Response{
		StatusCode: d.status,
		Status:     "",
		Body:       io.NopCloser(strings.NewReader("echo 123")),
	}

	return
}

func TestScript_Run(t *testing.T) {
	for _, test := range []struct {
		name         string
		script       string
		expectError  bool
		expectLog    []string
		expectStatus string
	}{
		{"successful script", "testdata/scripts/success.sh", false, []string{"+ echo hello world", "hello world", ""}, "completed successfully"},
		{"failing script", "testdata/scripts/fail.sh", true, []string{"uh-oh", ""}, "task failed"},
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

			if test.expectStatus != s.RunningState.String() {
				t.Errorf("expected %q, received %q", test.expectStatus, s.RunningState.String())
			}
		})
	}
}

func TestScript_Run_Download(t *testing.T) {
	url := "https://gist.github.com/jspc/af5adde4977610bbb9eee7efcaf14dee/raw/986f6a3f6be4f07ff0fd00f2c80bc25c3e658b9a/vc-test-script.sh"

	for _, test := range []struct {
		name        string
		client      httpClient
		expectError bool
	}{
		{"happy path", dummyHttpClient{status: 200}, false},
		{"http client errors are returned upwards", dummyHttpClient{err: true}, true},
		{"http 404s error out", dummyHttpClient{status: 404}, true},
	} {
		t.Run(test.name, func(t *testing.T) {
			client = test.client

			s := Script{
				Interpreter: &sh,
				URL:         &url,
			}

			err := s.Run()
			if err == nil && test.expectError {
				t.Error("expected error, received none")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", s.Error)
			}

		})
	}
}
