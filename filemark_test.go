package filemark

import (
	"log"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/afero"
)

var (
	fs afero.Fs
)

func init() {
	fs = afero.NewMemMapFs()
	err := fs.Mkdir("t", 0777)
	if err != nil {
		log.Fatal(err)
	}
}

func file(s string) (afero.File, error) {
	fn := strconv.FormatInt(time.Now().UnixNano(), 10)
	fn = filepath.Join("t", fn)
	f, err := fs.Create(fn)
	if err != nil {
		return nil, err
	}
	_, err = f.WriteString(s)
	if err != nil {
		return nil, err
	}
	return f, nil
}

var marksOnEmptyFile = []int{0, 1, 2, 3, 4, 5}

// TestEmptyFile an empty source
func TestEmptyFile(t *testing.T) {
	f, err := file("")
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

var cobai = `0
1
2
3
4
5
6
7
8
9`

// TestCobaiFile a known, predictable file
func TestEmptySeparator(t *testing.T) {
	f, err := file(cobai)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	mk := New(f, "\n")

	t.Run("PartSize", func(t *testing.T) {
		n := 5
		k := mk.PartSize(n)
		t.Log("PartSize got ", k)
	})

	t.Run("", func(t *testing.T) {
		n := 2
		k := mk.Marks(n)
		t.Log("marks got ", k)
	})
}
