package filemark

import (
	"bufio"
	"io"
	"log"
	"os"
)

// PartReader is a reader who knows his limits
type PartReader struct {
	f *os.File
	a int64
	z int64
	r io.Reader
}

// NewPartReader return a reader that reads between offsets
func NewPartReader(f *os.File, a int64, z int64) *PartReader {
	a, err := f.Seek(a, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	lr := &io.LimitedReader{
		R: f,
		N: z - a,
	}
	return &PartReader{f, a, z, bufio.NewReader(lr)}
}

func (pr *PartReader) Read(p []byte) (n int, err error) {
	crt, err := pr.f.Seek(0, io.SeekCurrent)
	if crt >= pr.z {
		return 0, io.EOF
	}
	if err != nil {
		return
	}
	return pr.r.Read(p)
}
