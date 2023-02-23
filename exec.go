package vc

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	client httpClient = http.DefaultClient
)

type httpClient interface {
	Get(string) (*http.Response, error)
}

type ScriptRunningState uint8

func (s ScriptRunningState) String() string {
	switch s {
	case ScriptRunningState_NotStarted:
		return "not started"
	case ScriptRunningState_Started:
		return "started"
	case ScriptRunningState_Success:
		return "completed successfully"
	case ScriptRunningState_Failure:
		return "task failed"
	}

	return "unknown"
}

const (
	ScriptRunningState_NotStarted ScriptRunningState = iota
	ScriptRunningState_Started
	ScriptRunningState_Success
	ScriptRunningState_Failure
)

func (s *Script) Run() (err error) {
	var (
		data   []byte
		fn     string
		output []byte
	)

	defer func() {
		if err != nil {
			s.Error = err.Error()
		}
	}()

	if s.URL != nil {
		data, err = download(*s.URL)
		if err != nil {
			return
		}
	} else {
		data = []byte(*s.Body)
	}

	fn, err = writedata(data)
	if err != nil {
		return
	}

	s.RunAt = time.Now()
	s.RunningState = ScriptRunningState_Started

	c := exec.Command(*s.Interpreter, fn) // #nosec G204
	output, err = c.CombinedOutput()

	s.Logs = strings.Split(string(output), "\n")

	switch err {
	case nil:
		s.RunningState = ScriptRunningState_Success

	default:
		s.RunningState = ScriptRunningState_Failure
	}

	return
}

func download(url string) (data []byte, err error) {
	resp, err := client.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s returned %s", url, resp.Status)

		return
	}

	defer resp.Body.Close() // #nosec: G307

	b := bytes.Buffer{}
	_, err = b.ReadFrom(resp.Body)

	data = b.Bytes()

	return
}

func writedata(data []byte) (fn string, err error) {
	f, err := os.CreateTemp("", "")
	if err != nil {
		return
	}

	fn = f.Name()
	defer f.Close() // #nosec G307

	_, err = f.Write(data)

	return
}
