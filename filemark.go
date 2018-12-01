package filemark

import (
	"io"
	"math"
)

const zero64 = int64(0)

// Filemark provides needed structure
type Filemark struct {
	f     io.ReaderAt
	delim string
	size  int64
	err   error
}

// New constructs a Filemark pointer
func New(f io.ReaderAt, d string, sz int64) *Filemark {
	mk := &Filemark{f, d, sz, nil}
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

	ldelim := int64(len(mk.delim))
	partsize, filesize := mk.PartSize(numberOfParts), mk.size
	pos := partsize

eof:
	for {
		buf := make([]byte, ldelim)
		n, err := mk.f.ReadAt(buf, pos)
		if err != nil || int64(n) != ldelim {
			if err == io.EOF {
				break eof
			}
			mk.err = err
			return zero
		}
		for {
			pos += ldelim
			if string(buf) == mk.delim {
				marks = append(marks, pos)
				break
			}
			n, err = mk.f.ReadAt(buf, pos)
			if err != nil {
				if err == io.EOF {
					break eof
				}
				mk.err = err
				return zero
			}
		}
		pos += partsize
	}
	marks = append(marks, filesize)
	marks = append([]int64{0}, marks...)
	return marks
}

// Err returns the mistake that breaked "the happy flow"
func (mk *Filemark) Err() error {
	return mk.err
}

// PartSize calculates single part length of n parts
func (mk *Filemark) PartSize(n int) int64 {
	// make number of parts to be natural
	if n < 1 {
		n = 1
	}
	eaten := int64(len(mk.delim) * n)
	return int64(math.Ceil(float64(mk.size-eaten) / float64(n)))
}
