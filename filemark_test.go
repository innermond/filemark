package filemark

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

// num of parts
var numPiecesEmptyFile = []int{0, 1, 2, 3, 4, 5}

// TestEmptyFile an empty source
func TestZeroFile(t *testing.T) {
	f, fi, err := File("test_files/empty_file")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "", fi.Size())

	t.Run("always returns extremities 0 as start and 0 (size) as final", func(t *testing.T) {
		for _, n := range numPiecesEmptyFile {
			k := mk.Marks(n)
			t.Log()
			if len(k) != 2 {
				t.Errorf("actual two marks from %v extremities 0 as start and 0 as final", k)
			}
		}
	})
}

// key = num of parts
// value = marks actual to be returned
var marksEmptySeparator = map[int][]int64{
	0: []int64{0, 10},
	1: []int64{0, 10},
	2: []int64{0, 5, 10},
	3: []int64{0, 4, 8, 10},
}

// TestEmptySeparator using a known file with empty separator
func TestEmptySeparator(t *testing.T) {
	f, fi, err := File("test_files/empty_separator")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "", fi.Size())

	t.Run("Marks", func(t *testing.T) {
		for n, actual := range marksEmptySeparator {
			got := mk.Marks(n)
			assert.Equal(t, got, actual, fmt.Sprintf("got %v instead %v", got, actual))
		}
	})
}

// key = num of parts
// value = marks actual to be returned
var marksNewlineSeparator = map[int][]int64{
	0: []int64{0, 20},
	1: []int64{0, 20},
	2: []int64{0, 10, 20},
	3: []int64{0, 8, 16, 20},
}

// TestEmptySeparator using a known file with empty separator
func TestNewlineSeparator(t *testing.T) {
	f, fi, err := File("test_files/newline")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "\n", fi.Size())

	t.Run("Marks", func(t *testing.T) {
		for n, actual := range marksNewlineSeparator {
			got := mk.Marks(n)
			assert.Equal(t, got, actual, fmt.Sprintf("%d = got %v instead %v", n, got, actual))
		}
	})
}

<<<<<<< HEAD
var emptySeparatorCases = map[int][]int64{
	//0: []int64{20},
	//1: []int64{20},
	2: []int64{12, 20},
	//3: []int64{8, 16, 20},
=======
// key = num of parts
// value = marks actual to be returned
var marksFancySeparator = map[int][]int64{
	0: []int64{0, 19},
	1: []int64{0, 19},
	2: []int64{0, 13, 19},
	3: []int64{0, 9, 17, 19},
>>>>>>> 4334fc90f01ee25784834d45e023ba08acce2382
}

// TestEmptySeparator using a known file with empty separator
func TestFancySeparator(t *testing.T) {
	f, fi, err := File("test_files/fancy_separator")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "xx", fi.Size())
	t.Run("Marks", func(t *testing.T) {
		for n, actual := range marksFancySeparator {
			ps := mk.PartSize(n)
			t.Log(n, ps)
			got := mk.Marks(n)
			assert.Equal(t, got, actual, fmt.Sprintf("%d = got %v instead %v", n, got, actual))
		}
	})
}
