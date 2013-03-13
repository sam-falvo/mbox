// vim: ts=8 ai noexpandtab

// The mbox package provides a means of reading legacy MBOX-format e-mail files.
package mbox


import (
	"bufio"
	"fmt"
	"io"
	"strings"
)


// MboxReader objects represent streams of e-mail messages that happen to be in MBOX format.
// TODO(sfalvo): Rename this to MsgStream.
type MboxReader struct {
	prefetch	[]byte
	prefetchLength	int
	currentLine	int
	r		*bufio.Reader
}

// The ReadMessage method parses the input for another complete message.  A message consists of a From
// header, at least one header, followed by a collection of lines of text corresponding to the body of
// the message.
func (m *MboxReader) ReadMessage() (msg *Message, err error) {
	msg = &Message{
		headers: make(map[string][]string, 0),
	}

	msg.SendingAddress, err = m.parseFrom()
	if err != nil {
		msg = nil
		return
	}

	return
}

// errorf provides an error object whose string also includes the line-number.
// TODO(sfalvo): This prevents the user from testing for specific error responses.
// Find a better way of exposing the line on which an error occurs.
func (m *MboxReader) errorf(format string, args ...interface{}) error {
	s := fmt.Sprintf(format, args...)
	return fmt.Errorf("%d:%s", m.currentLine, s)
}

// parseFrom will succeed only if the current line of the mbox file is a properly
// formed "From " separator.  It will extract the sending e-mail address from this
// line.  If this line doesn't exist, it yields an error instead.
func (m *MboxReader) parseFrom() (who string, err error) {
	who, err = extractSendingAddress(m)
	if err == nil {
		err = m.nextLine()
	}
	return
}

func extractSendingAddress(m *MboxReader) (who string, err error) {
	if string(m.prefetch[0:5]) != "From " {
		return "", m.errorf("Mbox file not properly framed; 'From ' expected")
	}
	if m.prefetchLength < 6 {
		return "", m.errorf("Sender address expected")
	}
	who = strings.TrimSpace(string(m.prefetch[5:]))
	if who == "" {
		return "", m.errorf("Sender address cannot be whitespace")
	}
	return
}

// CreateMboxReader decorates an io.Reader instance with an mbox parser.
// It will produce an error if the file doesn't appear to be an mbox-formatted file.
// It determines this by verifying the first five characters of the file matches "From " (note the space).
// Observe, however, that CreateMboxReader() succeeding does not imply that it actually is a correctly formatted mbox file.
func CreateMboxReader(s io.Reader) (m *MboxReader, err error) {
	m = &MboxReader {
		prefetch: make([]byte, 1000),
		r: bufio.NewReader(s),
	}

	err = m.nextLine()
	if err != nil {
		m = nil
		return
	}

	_, err = extractSendingAddress(m)
	return
}

// nextLine retrieves the next logical line from the mbox file.  The caller
// should be concerned with one of three cases:
//
// - A successful read yields no error.
// - Attempting to read past the end of the input stream yields io.EOF.
// - All other errors are reported as necessary.
func (m *MboxReader) nextLine() error {
	slice, err := m.r.ReadSlice('\n')
	if (err != nil) && (err != io.EOF) {
		return err
	}
	copy(m.prefetch, slice)
	m.prefetch = m.prefetch[0:len(slice)]
	m.prefetchLength = len(m.prefetch)
	m.currentLine++
	return nil
}

