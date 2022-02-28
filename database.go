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
	f, err := os.Open(fn) // #nosec G304
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
	f, err := os.Create(fn) // #nosec G304
	if err != nil {
		return
	}

	defer f.Close() // #nosec G307

	enc := gob.NewEncoder(f)

	return enc.Encode(d)
}
