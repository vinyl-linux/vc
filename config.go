package vc

import (
	"fmt"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/cnf/structhash"
)

type Finaliser uint8

const (
	Finaliser_Nothing Finaliser = iota
	Finaliser_Reboot
	Finaliser_Shutdown
)

var (
	ErrMissingScript      = fmt.Errorf("missing script section")
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

func (f *Finaliser) UnmarshalText(text []byte) (err error) {
	t := string(text)

	switch t {
	case "", "default", "nothing":
		*f = Finaliser_Nothing
	case "reboot":
		*f = Finaliser_Reboot
	case "shutdown":
		*f = Finaliser_Shutdown
	default:
		err = fmt.Errorf("invalid type %q; must be in set (%q,%q,%q)",
			t, "nothing", "reboot", "shutdown")
	}

	return
}

type Config struct {
	Name        string    `toml:"name"`
	Description string    `toml:"description"`
	Script      *Script   `toml:"script"`
	Finaliser   Finaliser `toml:"finaliser"`

	Sum string `toml:"-"`
}

func LoadConfig(fn string) (c Config, err error) {
	_, err = toml.DecodeFile(fn, &c)
	if err != nil {
		return
	}

	if c.Script == nil {
		err = ErrMissingScript

		return
	}

	err = c.Script.Validate()
	if err != nil {
		return
	}

	c.Sum, err = structhash.Hash(c.Script, 1)

	return
}
