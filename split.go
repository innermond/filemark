package filemark

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"
)

// Split a file into parts that ends with delimiter
func Split(fn string, delim string, pieces int) error {
	var err error
	// open file
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	mk := New(f, delim)
	zz := mk.Marks(pieces)
	// jumps from mark to mark reading between
	var wg sync.WaitGroup
	wg.Add(len(zz))
	fail := make(chan error, 1)
	done := make(chan bool, 1)
	for i, z := range zz {
		go func(z int64, i int) {
			defer wg.Done()
			// set a to previous z
			a := int64(0)
			if i > 0 {
				a = zz[i-1]
			}
			// create a part file
			nm := fmt.Sprintf("%d.%s.%s-%s", i, f.Name(), strconv.FormatInt(a, 10), strconv.FormatInt(z, 10))
			fp, err := os.Create(nm)
			if err != nil {
				fail <- err
				return
			}
			defer fp.Close()

			// fill part file with mark reader content
			mr := io.NewSectionReader(f, a, z-a)
			_, err = io.Copy(fp, mr)
			if err != nil {
				fail <- err
				return
			}
		}(z, i)
	}

	go func() {
		wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case err = <-fail:
		close(fail)
		return err

	}
}
