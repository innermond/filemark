package filemark

import (
	"io"
	"math"
	"os"
)

const zero64 = int64(0)

// Filemark provides needed structure
type Filemark struct {
	f     *os.File
	delim string
	err   error
}

// New constructs a Filemark pointer
func New(f *os.File, d string) *Filemark {
	mk := &Filemark{f, d, nil}
	return mk
}

// Marks find where are positioned marks
// trying to respect number of parts and
// taking into account the delimiter
func (mk *Filemark) Marks(numberOfParts int) []int64 {
	zero := []int64{0}
	marks := []int64{}

	// previous error prevents any continuation
	if mk.Err() != nil {
		return zero
	}
	// proxy variable
	f := mk.f
	psz := mk.PartSize(numberOfParts)
	for {
		// move pointer
		_, mk.err = f.Seek(psz, io.SeekCurrent)
		if mk.Err() != nil {
			break
		}
		// find delim
		z := mk.findelim()
		if mk.Err() != nil {
			break
		}
		marks = append(marks, z)
	}
	marks = append(marks, mk.Size())
	return marks
}

func (mk *Filemark) findelim() int64 {
	f := mk.f
	m := make([]byte, len(mk.delim))
	mc := 0
	z := zero64
	for {
		mc, mk.err = f.Read(m)
		if mk.err != nil {
			return zero64
		}
		m = m[:mc]
		if string(m) == mk.delim {
			break
		}
	}
	z, mk.err = f.Seek(0, io.SeekCurrent)
	return z
}

// Err returns the mistake that breaked "the happy flow"
func (mk *Filemark) Err() error {
	return mk.err
}

// PartSize calculates single part length of n parts
func (mk *Filemark) PartSize(n int) int64 {
	sz := mk.Size()
	return int64(math.Ceil(float64(sz) / float64(n)))
}

// Size calculate just bytes number of a file
func (mk *Filemark) Size() int64 {
	f := mk.f
	// keep original position
	orig, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		mk.err = err
		return orig
	}
	// move to the file's end to get size
	size, err := f.Seek(0, io.SeekEnd)
	// reset to original position
	orig, origErr := f.Seek(orig, io.SeekStart)
	if origErr != nil {
		mk.err = origErr
		return orig
	}
	if err != nil {
		mk.err = err
		return size
	}

	return size
}
