package filemark

import (
	"bufio"
	"io"
	"log"
)

// MarkReader is a reader who knows his limits
type MarkReader struct {
	f io.ReadWriteSeeker
	a int64
	z int64
	r io.Reader
}

// NewMarkReader return a reader that reads between offsets
func NewMarkReader(f io.ReadWriteSeeker, a int64, z int64) *MarkReader {
	a, err := f.Seek(a, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}
	lr := &io.LimitedReader{
		R: f,
		N: z - a,
	}
	return &MarkReader{f, a, z, bufio.NewReader(lr)}
}

// Size returns size that can be read from this reader, distance between two consecutive marks
func (pr *MarkReader) Size() int {
	return int(pr.z - pr.a)
}

func (pr *MarkReader) Read(p []byte) (n int, err error) {
	crt, err := pr.f.Seek(0, io.SeekCurrent)
	if crt >= pr.z {
		return 0, io.EOF
	}
	if err != nil {
		return
	}
	return pr.r.Read(p)
}
