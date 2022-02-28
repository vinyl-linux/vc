package vc

import (
	"testing"
)

func TestLoadDatabase(t *testing.T) {
	for _, test := range []struct {
		name        string
		fn          string
		expectError bool
	}{
		{"Missing file returns empty database", "testdata/nonsuch.db", false},
		{"Non-db file returns error", "testdata/dodgy.db", true},
		{"Valid database is loaded properly", "testdata/database", false},
	} {
		t.Run(test.name, func(t *testing.T) {
			_, err := LoadDatabase(test.fn)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}

func TestWriteDatabase(t *testing.T) {
	validDb, _ := LoadDatabase("testdata/database")

	for _, test := range []struct {
		name        string
		fn          string
		db          Database
		expectError bool
	}{
		{"Non-writable file errors", "testdata/this/path/cant/be/written/to.db", Database{}, true},
		{"Valid database can be written to a valid path", "testdata/database", validDb, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := test.db.WriteDatabase(test.fn)
			if err == nil && test.expectError {
				t.Errorf("expected error")
			} else if err != nil && !test.expectError {
				t.Errorf("unexpected error: %+v", err)
			}
		})
	}
}
