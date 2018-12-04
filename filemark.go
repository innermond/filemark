package filemark

import (
	"io"
	"math"
)

const (
	zero64 = int64(0)
	one64  = int64(1)
)

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

	var (
		n   int
		err error
	)

	var found bool
	// loop to parts
eof:
	for {
		// loop to find a delimiter around a part
		for {
			n, err = mk.f.ReadAt(buf, pos)
			if err != nil && int64(n) != lbuf {
				if err == io.EOF {
					break eof
				}
				mk.err = err
				return zero
			}
			found = isDelimEmpty || string(buf) == mk.delim
			if found {
				// go right after the new found delimiter
				pos += lbuf
				// abort when pos goes out of size limit
				if pos >= size {
					break eof
				}
				marks = append(marks, pos)
				break
			}
			if !isDelimEmpty {
				// advance one byte more searching for delimiter
				pos += one64
			}
		}
		// increment position
		pos += partsize
		if pos >= size {
			break eof
		}
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
	leaten := len(mk.delim)
	if leaten == 0 {
		leaten = 1
	}
	eaten := int64(leaten * n)
	delta := mk.size - eaten
	if delta <= 0 {
		return zero64
	}
	return int64(math.Ceil(float64(delta) / float64(n)))
}
