package filemark

import (
	"io"
	"log"
	"math"
)

const zero64 = int64(0)

// Filemark provides needed structure
type Filemark struct {
	f     io.ReadSeeker
	delim string
	err   error
}

// New constructs a Filemark pointer
func New(f io.ReadSeeker, d string) *Filemark {
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
	// reset pointer to begining
	f.Seek(0, io.SeekStart)

	partsize, filesize := mk.PartSize(numberOfParts), mk.Size()
	for {
		// 1. move pointer
		_, mk.err = f.Seek(partsize, io.SeekCurrent)
		if mk.err != nil {
			break
		}
		// because f.Seek seems to not be able to signal EOF
		// check for especially for io.EOF by trying ro read
		bk := make([]byte, 1)
		_, err := f.Read(bk)
		if err == io.EOF {
			mk.err = io.EOF
			break
		}
		// 2. find delim
		z := mk.findelim()
		/*bk = make([]byte, 1)
		log.Print(z)
		fAt := f.(*os.File)
		n, err := fAt.ReadAt(bk, z)
		log.Print(n, err, string(bk))*/
		// error or be paranoid and watch position to not go beyond file size
		if mk.Err() != nil || z >= filesize {
			break
		}
		marks = append(marks, z)
	}
	marks = append(marks, filesize)
	// reset pointer to begining
	f.Seek(0, io.SeekStart)

	return marks
}

func (mk *Filemark) findelim() int64 {
	f := mk.f
	ldelim := len(mk.delim)
	z := zero64
	m := make([]byte, ldelim)
	mc := 0
	if ldelim == 0 {
		z, mk.err = f.Seek(0, io.SeekCurrent)
		return z
	}

	for {
		mc, mk.err = f.Read(m)
		if mk.err != nil {
			return zero64
		}
		z, mk.err = f.Seek(0, io.SeekCurrent)
		m = m[:mc]
		log.Print(z, mk.err, m)
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
	// make number of parts to be natural
	if n < 1 {
		n = 1
	}
	sz := mk.Size()
	eaten := int64(len(mk.delim) * n)
	return int64(math.Ceil(float64(sz-eaten) / float64(n)))
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
