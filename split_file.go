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
func SplitFile(f *os.File, delim string, size int64, pieces int, step int) []error {
	var err error
	mrr := Split(f, delim, size, pieces)
	lenr := len(mrr)
	if step == 0 {
		step = lenr
	}
	// jumps from mark to mark reading between
	fail := make(chan error, lenr)
	done := make(chan bool, lenr)

	// launch workers on separate gouroutines
	// WARN: loop can raise to many file handlers
	go func() {
		var wg sync.WaitGroup
		// as a number of steps
		wall := lenr / step
		// leftover
		extra := lenr % step
		// fit perfectly as steps accumulation
		wall *= step
		// add an entire step
		if extra != 0 {
			wall += step
		}
		for xstep := 0; xstep < wall; xstep += step {
			a := xstep
			z := xstep + step
			if z > lenr {
				z = lenr
			}
			sliced := mrr[a:z]
			wg.Add(len(sliced))
			// i 0..step
			for i, mr := range sliced {
				i := i
				mr := mr
				// worker here
				go func() {
					defer wg.Done()
					// name part file
					nm := path.Base(f.Name())
					nm = fmt.Sprintf("%d.%s", i+xstep, nm)
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
			// wait for step to to be executes
			wg.Wait()
		}
		close(fail)
	}()

	var errs []error
	// lisen here for workers signals
	for j := 0; j < lenr; j++ {
		select {
		case <-done:
		case err = <-fail:
			errs = append(errs, err)
		}
	}
	close(done)

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
