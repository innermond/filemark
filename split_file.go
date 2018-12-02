package filemark

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// SplitFile a file into parts that ends with delimiter
func SplitFile(f *os.File, delim string, pieces int) error {
	size, err := Size(f)
	if err != nil {
		return err
	}
	mrr := Split(f, delim, size, pieces)
	// jumps from mark to mark reading between
	var wg sync.WaitGroup
	wg.Add(len(mrr))
	fail := make(chan error, 1)
	done := make(chan bool, 1)
	for i, mr := range mrr {
		i := i
		mr := mr
		go func() {
			defer wg.Done()
			// create a part file
			nm := fmt.Sprintf("%d.%s", i, f.Name())
			fp, err := os.Create(nm)
			if err != nil {
				fail <- err
				return
			}
			defer fp.Close()

			// fill part file with mark reader content
			_, err = io.Copy(fp, mr)
			if err != nil {
				fail <- err
				return
			}
		}()
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

// Size calculate just bytes number of a file
func Size(f *os.File) (int64, error) {
	// keep original position
	orig, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return orig, err
	}
	// move to the file's end to get size
	size, err := f.Seek(0, io.SeekEnd)
	// reset to original position
	orig, origErr := f.Seek(orig, io.SeekStart)
	if origErr != nil {
		return orig, err
	}
	if err != nil {
		return size, err
	}

	return size, nil
}
