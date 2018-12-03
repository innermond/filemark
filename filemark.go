package filemark

import (
	"fmt"
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

	// previous error prevents any continuation
	if mk.Err() != nil {
		return zero
	}

	partsize, size := mk.PartSize(numberOfParts), mk.size
	pos := partsize

	marks := []int64{}
	ldelim := int64(len(mk.delim))
	isDelimEmpty := ldelim == zero64

	lbuf := ldelim
	// in order to read buf must be at least 1 byte
	if isDelimEmpty {
		lbuf = int64(1)
	}
	buf := make([]byte, lbuf)

	var err error

eof:
	for {
		// find a delimiter advancing by delimiter length
		for {
			pos += lbuf
			fmt.Println(ldelim, partsize, isDelimEmpty, buf, mk.delim, string(buf) == mk.delim)
			if string(buf) == mk.delim {
				marks = append(marks, pos)
				break
			}
			_, err = mk.f.ReadAt(buf, pos)
			if err != nil {
				if err == io.EOF {
					break eof
				}
				mk.err = err
				return zero
			}
		}
		// increment position
		pos += partsize
		// reset buffer
		buf = nil
	}
	// add extremities 0 and size
	marks = append([]int64{0}, marks...)
	marks = append(marks, size)
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
