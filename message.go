// vim: ts=8 noexpandtab ai

// The message package provides message-related definitions.
package mbox

import "io"

// A Message represents a single message in the file.
type Message struct {
	mbox		*MboxReader
	headers		map[string][]string
	SendingAddress	string
}

// A bodyReader implements an io.Reader, confined to the current message to
// which this instance is bound.
type bodyReader struct {
	msg		*Message
	mbox		*MboxReader
	where		int
	srcErr		error
}


// The Headers method provides raw access to the headers of a message.
func (m *Message) Headers() map[string][]string {
	return m.headers
}

// BodyReader() provides an io.Reader compatible object that will read the body
// of the message.  It will return io.EOF if you attempt to read beyond the end
// of the message.
func (m *Message) BodyReader() io.Reader {
	br := &bodyReader{
		msg: m,
		mbox: m.mbox,
	}

	return br
}


func (r *bodyReader) Read(bs []byte) (n int, err error) {
	if r.srcErr != nil {
		return 0, r.srcErr
	}

	if (len(r.mbox.prefetch) > 5) && (string(r.mbox.prefetch[0:5]) == "From ") {
		return 0, io.EOF
	}

	n = copy(bs, r.mbox.prefetch[r.where:])
	r.where = r.where + n
	if r.where >= len(r.mbox.prefetch) {
		r.where = 0
		r.srcErr = r.mbox.nextLine()
	}
	return
}

