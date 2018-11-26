package filemark

import (
	"fmt"
	"io"
	"os"
	"strconv"
)

// Split a file into parts that ends with delimiter
func Split(fn string, delim string, pieces int) error {
	// open file
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	mk := New(f, delim)
	zz := mk.Marks(pieces)
	// jumps from mark to mark reading between
	a := int64(0)
	for i, z := range zz {
		// create a part file
		nm := fmt.Sprintf("%d.%s.%s-%s", i, f.Name(), strconv.FormatInt(a, 10), strconv.FormatInt(z, 10))
		fp, err := os.Create(nm)
		if err != nil {
			return err
		}
		defer fp.Close()

		// fill part file with mark reader content
		mr := io.NewSectionReader(f, a, z-a)
		_, err = io.Copy(fp, mr)
		if err != nil {
			return err
		}
		a = z

	}
	return nil
}
