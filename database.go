package vc

import (
	"encoding/gob"
	"errors"
	"io/fs"
	"os"
)

// hash to Script map
type Database map[string]Script

// load from file
func LoadDatabase(fn string) (d Database, err error) {
	f, err := os.Open(fn)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			d = make(Database)

			err = nil
		}

		return
	}

	dec := gob.NewDecoder(f)

	err = dec.Decode(&d)

	return
}

func (d Database) WriteDatabase(fn string) (err error) {
	f, err := os.Create(fn)
	if err != nil {
		return
	}

	enc := gob.NewEncoder(f)

	return enc.Encode(d)
}
