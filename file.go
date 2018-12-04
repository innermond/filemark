package filemark

import (
	"os"
)

// File gives a file handler
func File(fn string) (*os.File, os.FileInfo, error) {
	f, err := os.Open(fn)
	if err != nil {
		return nil, nil, err
	}
	fi, err := f.Stat()
	if err != nil {
		return nil, nil, err
	}
	return f, fi, nil
}
