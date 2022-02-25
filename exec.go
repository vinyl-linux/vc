package vc

import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

type ScriptRunningState uint8

const (
	ScriptRunningState_NotStarted ScriptRunningState = iota
	ScriptRunningState_Started
	ScriptRunningState_Succes
	ScriptRunningState_Failure
)

func (s *Script) Run() {
	var (
		data   []byte
		fn     string
		output []byte
	)

	if s.URL != nil {
		data, s.Error = download(*s.URL)
		if s.Error != nil {
			return
		}
	} else {
		data = []byte(*s.Body)
	}

	fn, s.Error = writedata(data)
	if s.Error != nil {
		return
	}

	c := exec.Command(*s.Interpreter, fn)
	output, s.Error = c.CombinedOutput()

	s.Logs = strings.Split(string(output), "\n")

	return
}

func download(url string) (data []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = fmt.Errorf("%s returned %s", url, resp.Status)

		return
	}

	defer resp.Body.Close()

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
	defer f.Close()

	_, err = f.Write(data)

	return
}
