package filemark

import (
	"fmt"
	"io"
	"os"
	"sync"
)

// SplitFile a file into parts that ends with delimiter
func SplitFile(fn string, delim string, pieces int) error {
	var err error
	// open file
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	mrr := Split(f, delim, pieces)
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