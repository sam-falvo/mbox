// The message package provides message-related definitions.
package mbox

// A Message represents a single message in the file.
type Message struct {
	headers		map[string][]string
	SendingAddress	string
}


// The Headers method provides raw access to the headers of a message.
func (m *Message) Headers() map[string][]string {
	return m.headers
}


