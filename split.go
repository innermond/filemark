package filemark

import (
	"io"
)

// ReadAtSeeker can do
type ReadAtSeeker interface {
	io.ReadSeeker
	io.ReaderAt
}

// SectionReaders colletion of readers that knows their limits
type SectionReaders []*io.SectionReader

// Split a reader into readers devoted to parts that ends with delimiter
func Split(r ReadAtSeeker, delim string, pieces int) SectionReaders {
	var mrr SectionReaders
	mk := New(r, delim)
	zz := mk.Marks(pieces)
	for i, z := range zz {
		// set a to previous z
		a := int64(0)
		if i > 0 {
			a = zz[i-1]
		}
		mr := io.NewSectionReader(r, a, z-a)
		mrr = append(mrr, mr)
	}
	return mrr
}
