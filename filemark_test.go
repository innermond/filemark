package filemark

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

func file(s string) (*os.File, error) {
	fn := filepath.Join("test_files", s)
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	return f, nil
}

var marksOnEmptyFile = []int{0, 1, 2, 3, 4, 5}

// TestEmptyFile an empty source
func _TestEmptyFile(t *testing.T) {
	f, err := file("empty_file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "")

	t.Run("empty file has size 0", func(t *testing.T) {
		got := mk.Size()
		if got != zero64 {
			t.Errorf("expected %d but got %d", zero64, got)
		}
	})

	t.Run("empty file always returns one mark at zero pos", func(t *testing.T) {
		for _, n := range marksOnEmptyFile {
			k := mk.Marks(n)
			if len(k) != 1 {
				t.Errorf("expected one mark from %v", k)
			}
			got := k[0]
			if got != zero64 {
				t.Errorf("expected zero mark on a nil file but got %d", got)
			}
		}
	})
}

var emptySeparatorCases = map[int][]int64{
	//0: []int64{20},
	//1: []int64{20},
	2: []int64{12, 20},
	3: []int64{8, 16, 20},
}

// TestEmptySeparator using a known file with empty separator
func TestEmptySeparator(t *testing.T) {
	f, err := file("newline")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "\n")

	t.Run("Marks", func(t *testing.T) {
		for n, expected := range emptySeparatorCases {
			t.Log(n, "PartSize got ", mk.PartSize(n))
			got := mk.Marks(n)
			t.Log(fmt.Sprintf("marks for %d got", n), got, expected)
		}
	})
}
