package filemark

import (
	"io"
)

// SectionReaders colletion of readers that knows their limits
type SectionReaders []*io.SectionReader

// Split a reader into readers devoted to parts that ends with delimiter
func Split(r io.ReaderAt, delim string, size int64, pieces int) SectionReaders {
	var mrr SectionReaders
	mk := New(r, delim, size)
	zz := mk.Marks(pieces)
	for i := 1; i < len(zz); i++ {
		// set a to previous z
		a := zz[i-1]
		z := zz[i]
		mr := io.NewSectionReader(r, a, z-a)
		mrr = append(mrr, mr)
	}
	return mrr
}
