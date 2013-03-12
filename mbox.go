// vim: ts=8 ai noexpandtab

// The mbox package provides a method of reading legacy MBOX-format e-mail files.
package mbox


import (
	"fmt"
	"io"
)


// MboxReader objects represent streams of e-mail messages that happen to be in MBOX format.
// The MboxReader object also implements the io.ReaderAt interface.
type MboxReader struct {
	r io.ReaderAt
}


// CreateMboxReader decorates an io.ReaderAt instance with an mbox parser.
// It will produce an error if the file doesn't appear to be an mbox-formatted file.
// It determines this by verifying the first five characters of the file matches "From " (note the space).
// Observe, however, that CreateMboxReader() succeeding does not imply that it actually is a correctly formatted mbox file.
func CreateMboxReader(s io.ReaderAt) (m *MboxReader, err error) {
	bs := make([]byte, 6)
	m = nil

	// All MBOX files must begin with "From ", and at least one character for a sender address.
	n, err := s.ReadAt(bs, 0)
	if err != nil {
		return
	}
	if n < 6 {
		err = fmt.Errorf("Mbox file is too short to be valid")
		return
	}
	if string(bs[0:5]) != "From " {
		err = fmt.Errorf("First line of mbox file does not begin with \"From \" and a sender address")
		return
	}

	m = &MboxReader{
		r: s,
	}
	return
}

