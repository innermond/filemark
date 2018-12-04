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

// N = num of parts
// M = marks actual to be returned
var marksEmptySeparator = []struct {
	N int
	M []int64
}{
	{0, []int64{0, 10}},
	{1, []int64{0, 10}},
	{2, []int64{0, 5, 10}},
	{3, []int64{0, 4, 8, 10}},
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
		for _, tc := range marksEmptySeparator {
			got := mk.Marks(tc.N)
			assert.Equal(t, got, tc.M, fmt.Sprintf("got %v instead %v", got, tc.M))
		}
	})
}

// N = num of parts
// M = marks actual to be returned
var marksNewlineSeparator = []struct {
	N int
	M []int64
}{
	{0, []int64{0, 20}},
	{1, []int64{0, 20}},
	{2, []int64{0, 10, 20}},
	{3, []int64{0, 8, 16, 20}},
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
		for _, tc := range marksNewlineSeparator {
			got := mk.Marks(tc.N)
			assert.Equal(t, got, tc.M, fmt.Sprintf("%d = got %v instead %v", tc.N, got, tc.M))
		}
	})
}

// N = num of parts
// M = marks actual to be returned
var marksFancySeparator = []struct {
	N int
	M []int64
}{
	{0, []int64{0, 19}},
	{1, []int64{0, 19}},
	{2, []int64{0, 13, 19}},
	{3, []int64{0, 9, 17, 19}},
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
		for _, tc := range marksFancySeparator {
			got := mk.Marks(tc.N)
			assert.Equal(t, got, tc.M, fmt.Sprintf("%d = got %v instead %v", tc.N, got, tc.M))
		}
	})
}
