package filemark

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// SplitFile a file into parts that ends with delimiter
func SplitFile(f *os.File, delim string, size int64, pieces int) []error {
	var err error
	mrr := Split(f, delim, size, pieces)
	lenr := len(mrr)
	// jumps from mark to mark reading between
	var wg sync.WaitGroup
	wg.Add(lenr)
	fail := make(chan error, lenr)
	done := make(chan bool, lenr)

	// launch workers on separate gouroutines
	for i, mr := range mrr {
		i := i
		mr := mr
		// worker here
		go func() {
			defer wg.Done()
			// name part file
			nm := path.Base(f.Name())
			nm = fmt.Sprintf("%d.%s", i, nm)
			// create a part file
			nm = filepath.Clean(nm)
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
			done <- true
		}()
	}

	// wait for workers to finish their job
	go func() {
		wg.Wait()
		close(done)
	}()

	// lisen here for workers signals
	var errs []error
	for ; lenr > 0; lenr-- {
		select {
		case <-done:
		case err = <-fail:
			errs = append(errs, err)
		}
	}
	close(fail)

	return errs
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
