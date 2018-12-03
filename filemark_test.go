package filemark

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func file(s string) (*os.File, os.FileInfo, error) {
	fn := filepath.Join("test_files", s)
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

var marksOnEmptyFile = []int{0, 1, 2, 3, 4, 5}

// TestEmptyFile an empty source
func TestEmptyFile(t *testing.T) {
	t.Skip()
	f, fi, err := file("empty_file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "", fi.Size())

	t.Run("empty file always returns one mark at zero pos", func(t *testing.T) {
		for _, n := range marksOnEmptyFile {
			t.Log("val: ", n)
			k := mk.Marks(n)
			if len(k) != 2 {
				t.Errorf("expected two marks from %v extremities 0 as start and 0 as final", k)
			}
		}
	})
}

var emptySeparatorCases = map[int][]int64{
	//0: []int64{0, 20},
	//1: []int64{0, 20},
	2: []int64{0, 12, 20},
	3: []int64{0, 8, 16, 20},
}

// TestEmptySeparator using a known file with empty separator
func TestEmptySeparator(t *testing.T) {
	f, fi, err := file("newline")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "", fi.Size())

	t.Run("Marks", func(t *testing.T) {
		for n, expected := range emptySeparatorCases {
			t.Log(n, "PartSize got ", mk.PartSize(n))
			got := mk.Marks(n)
			t.Log(fmt.Sprintf("marks for %d got", n), got, expected)
		}
	})
}
