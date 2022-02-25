package vc

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cnf/structhash"
)

var (
	ErrMissingInterpreter = fmt.Errorf("interpreter is either empty or unset")
	ErrUrlOrBodyOnly      = fmt.Errorf("only one of url or body may be set at once")
	ErrUrlAndBodyUnset    = fmt.Errorf("either script url or body must be set")
)

type Script struct {
	Interpreter *string `toml:"interpreter"`
	URL         *string `toml:"url"`
	Body        *string `toml:"body"`

	RunningState ScriptRunningState `toml:"-" hash:"-"`
	Logs         []string           `toml:"-" hash:"-"`
	Error        string             `toml:"-" hash:"-"`
	RunAt        time.Time          `toml:"-" hash:"-"`
}

func (s Script) Validate() error {
	if s.Interpreter == nil {
		return ErrMissingInterpreter
	}

	if s.URL != nil && s.Body != nil {
		return ErrUrlOrBodyOnly
	}

	if s.URL == nil && s.Body == nil {
		return ErrUrlAndBodyUnset
	}

	return nil
}

type Config struct {
	Name        string  `toml:"name"`
	Description string  `toml:"description"`
	Script      *Script `toml:"script"`

	Sum string `toml:"-"`
}

func LoadConfig(fn string) (c Config, err error) {
	_, err = toml.DecodeFile(fn, &c)
	if err != nil {
		return
	}

	err = c.Script.Validate()
	if err != nil {
		return
	}

	c.Sum, err = structhash.Hash(c.Script, 1)

	return
}
