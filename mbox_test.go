// vim: ts=8 noexpandtab ai

package mbox

import "io"
import "strings"
import "testing"

const mboxWith1Message = `From foo@bar.com
Subject: Hello world

Test message
`

const mboxWith3Messages = `From foo@bar.com
Subject: Hello world

Test message



From foo@bar.com
From: Foo S. Ball <foo@bar.com>
To: Anyone W. Cares <anyone@bar.com>
Subject: You're all fired!

Haha, just joking.
I wasn't really trying to be a jerk.
It's just that it's April fools, and all.
From foo@bar.com
From: Foo S. Ball <foo@bar.com>
To: Loraine <amiga@bar.com>
Subject: Stella rules!

Old flames never die out.  They just smolder and smoke until you leave the room.
BTW, thanks for the Boing beach ball.

`

const mboxWithMessageNoHeaders = `From foo@bar.com

Test message
`

const mboxWithMessageNoAttribute = `From foo@bar.com
 continuation-line

Test message
`

const mboxWithMessageKeyMissing = `From foo@bar.com
: value-line

Test message
`

const mboxWithMessageHeaderWithContinuation = `From foo@bar.com
Subject: Hello
 world

Test message
`

const mboxWithMessage3Headers = `From foo@bar.com
From: foo@bar.com
To: user1@bar.com
 user2@bar.com
 user3@bar.com
 user4@bar.com
 user5@bar.com
Subject: Hello world

Greetings and hallucinations!
`

/* *** Test Utilities *** */

// in() returns true only if a string (needle) is found in an array of strings
// (haystack).
func in(needle string, haystack []string) (found bool) {
	found = false
	n := strings.TrimSpace(needle)
	for _, straw := range haystack {
		if strings.TrimSpace(straw) == n {
			found = true
		}
	}
	return
}

/* *** Setups of various kinds *** */

// withOpenMboxReader sets up a test.  It creates a MboxReader on a given
// string source.  If successful, it invokes the specified test, which then
// performs whatever checks it sees fit.
func withOpenMboxReader(t *testing.T, procname, source string, test func(mr *MboxReader)) {
	stringReader := strings.NewReader(source)
	mr, err := CreateMboxReader(stringReader)
	if err != nil {
		t.Error(procname, ": ", err)
		return
	}
	test(mr)
}

// withReadMessage sets up a test.  It creates a MboxReader on a
// known-good mbox file, then reads the first message in the mbox file.  The
// test then performs whichever checks it likes on the provided message.
func withReadMessage(t *testing.T, procname string, test func(msg *Message)) {
	withOpenMboxReader(t, procname, mboxWith1Message, func(mr *MboxReader) {
		msg, err := mr.ReadMessage()
		if err != nil {
			t.Error(procname, ": ", err)
			return
		}
		test(msg)
	})
}

// expectError() performs a basic sanity check for opening a new MboxReader object.
// This procedure checks for the absence of an error, and fails the test if found.
func expectError(t *testing.T, s string, msg string) {
	stringReader := strings.NewReader(s)
	_, err := CreateMboxReader(stringReader)
	if err == nil {
		t.Error(msg)
	}
}

// expectError() performs a basic sanity check for opening a new MboxReader object.
// This procedure checks for the existence of an error, and fails the test if found.
func expectNoError(t *testing.T, s string, msg string, pmr **MboxReader) {
	var err error

	stringReader := strings.NewReader(s)
	*pmr, err = CreateMboxReader(stringReader)
	if err != nil {
		t.Error(msg, ":", err)
	}
}

/* *** Test Cases *** */

// Given a corrupted mbox file with a missing From header on the first line
// When I try to open the file
// Then I expect an error.
func TestMalformedMboxFile10(t *testing.T) {
	expectError(
		t,
		"\nFrom foo\n",
		"Mbox files must start with \"From \"",
	)
}

// Given a corrupted mbox file with a From header on the first line
//  AND no sender address
// When I try to open the file
// Then I expect an error.
func TestMalformedMboxFile20(t *testing.T) {
	expectError(t, "From ", "Mbox files that are too short must produce an error")
}

// Given a corrupted mbox file with a valid size but an improperly spaced From line
// When I try to open the file
// Then I expect an error.
func TestMalformedMboxFile30(t *testing.T) {
	expectError(t, " From ", "Leading whitespace on the From line must produce an error")
}

// Given a corrupted mbox file with a valid size but an otherwise empty From line
//  AND I successfully open the file
// When I try to read the first message
// Then I expect an error.
func TestMalformedMboxFile40(t *testing.T) {
	expectError(t, "From   \t\t  \t\t", "Sender address cannot be whitespace")
}

// Given a valid mbox file
// When I try to open the file
// Then I expect no error and a valid MboxReader instance.
func TestOkMboxFile10(t *testing.T) {
	var mr *MboxReader

	expectNoError(t, mboxWith1Message, "Mbox file with one valid message should not yield an error.", &mr)

	if mr == nil {
		t.Error("Returned MboxReader is nil for some reason")
	}
}

// Given a valid mbox file
//  AND I successfully open the file
// When I try to read from the file
// Then I expect no error and a valid message instance.
func TestOkMboxFile20(t *testing.T) {
	withReadMessage(t, "TestOkMboxFile20", func(msg *Message) {
		if msg == nil {
			t.Error("Message instance is nil despite lack of error")
		}
	})
}

// Given a valid mbox file
//  AND I successfully open the file
// When I read from the file
// Then I expect a message with correct sending address.
func TestOkMboxFile30(t *testing.T) {
	withReadMessage(t, "TestOkMboxFile30", func(msg *Message) {
		if msg.SendingAddress != "foo@bar.com" {
			t.Error("TestOkMboxFile30: Expected valid sending address")
			return
		}
	})
}

// Given a valid mbox file
//  AND I successfully open the file
// When I read from the file
// Then I expect a message with correct headers.
func TestOkMboxFile40(t *testing.T) {
	withReadMessage(t, "TestOkMboxFile40", func(msg *Message) {
		hs := msg.Headers()
		if hs["Subject"][0] != "Hello world" {
			t.Error("TestOkMboxFile40: Subject isn't Hello World")
			return
		}
	})
}

// Given an invalid mbox file with zero headers
//  AND I successfully open the file
// When I read from the file
// Then I expect an error.
func TestMalformedMboxFile50(t *testing.T) {
	withOpenMboxReader(t, "TestMalformedMboxFile50", mboxWithMessageNoHeaders, func(mr *MboxReader) {
		_, err := mr.ReadMessage()
		if err == nil {
			t.Error("TestMalformedMboxFile50: Error expected for message with no headers")
			return
		}
	})
}

// Given an invalid mbox file with corrupted headers
//  AND I successfully open the file
// When I read from the file
// Then I expect an error.
func TestMalformedMboxFile60(t *testing.T) {
	withOpenMboxReader(t, "TestMalformedMboxFile60", mboxWithMessageNoAttribute, func(mr *MboxReader) {
		_, err := mr.ReadMessage()
		if err == nil {
			t.Error("TestMalformedMboxFile60: Error expected for missing 'key: value' syntax")
			return
		}
	})
}

// Given an invalid mbox file with a malformed key/value pair
//  AND I successfully open the file
// When I read from the file
// Then I expect an error.
func TestMalformedMboxFile70(t *testing.T) {
	withOpenMboxReader(t, "TestMalformedMboxFile70", mboxWithMessageKeyMissing, func(mr *MboxReader) {
		_, err := mr.ReadMessage()
		if err == nil {
			t.Error("TestMalformedMboxFile70: Error expected for missing 'key: value' syntax")
			return
		}
	})
}

// Given a valid mbox file with a key/value pair with at least one continuation line
// When I read the file
// Then I expect a key and a value of two strings.
func TestOkMboxFile50(t *testing.T) {
	withOpenMboxReader(t, "TestOkMboxFile50", mboxWithMessageHeaderWithContinuation, func(mr *MboxReader) {
		msg, err := mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile50: ", err)
			return
		}
		hs := msg.Headers()
		if len(hs) != 1 {
			t.Error("Only one header provided in source mbox content")
			return
		}
		if len(hs["Subject"]) != 2 {
			t.Error("One key/value line and one continuation line should give us two lines total")
			return
		}
		s := hs["Subject"]
		if s[0] != "Hello" {
			t.Error("String extraction seems to have failed for attribute line")
			return
		}
		if s[1] != " world" {
			t.Error("String extraction seems to have failed for continuation: ")
			return
		}
	})
}

// Given a valid mbox with a message using three headers
// When I read the message
// Then I expect to see all three headers.
func TestOkMboxFile60(t *testing.T) {
	withOpenMboxReader(t, "TestOkMboxFile60", mboxWithMessage3Headers, func(mr *MboxReader) {
		msg, err := mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile60: ", err)
			return
		}
		hs := msg.Headers()
		if len(hs) != 3 {
			t.Error("Expected three headers")
			return
		}
		if len(hs["Subject"]) != 1 {
			t.Error("Subject should have one value string")
			return
		}
		if len(hs["To"]) != 5 {
			t.Error("Multiple recipients should be listed")
			return
		}
		if len(hs["From"]) != 1 {
			t.Error("From: header should have one value")
			return
		}
		ff := hs["From"][0]
		tt := hs["To"]
		ss := hs["Subject"][0]
		if ff != "foo@bar.com" {
			t.Error("From: header has wrong value")
			return
		}
		for _, a := range []string{"user1@bar.com", "user2@bar.com", "user3@bar.com", "user4@bar.com", "user5@bar.com"} {
			if !in(a, tt) {
				t.Error("To header values missing an expected value: ", a)
				return
			}
		}
		if ss != "Hello world" {
			t.Error("Subject heading is wrong")
			return
		}
	})
}

// Given a valid mbox with a message with a body
// When I read the message
// Then I expect to access an io.Reader that lets me read in the body.
func TestOkMboxFile70(t *testing.T) {
	withOpenMboxReader(t, "TestOkMboxFile70", mboxWith1Message, func(mr *MboxReader) {
		msg, err := mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile70: ", err)
			return
		}
		br := msg.BodyReader()
		bs := make([]byte, 128)
		n, err := br.Read(bs)
		if err != nil {
			t.Error("TestOkMboxFile70: ", err)
			return
		}
		if n < 13 {
			t.Error("Expected 13 characters, got ", n)
			return
		}
		bs = bs[0:n]
		if string(bs) != "Test message\n" {
			t.Error("Expected Test message, but got ", string(bs))
			return
		}
		n, err = br.Read(bs)
		if err != io.EOF {
			t.Error("Expected io.EOF; got ", err, n, string(bs))
		}
	})
}

// Given a valid mbox with three messages
// When I read the messages,
// Then I expect to see each message in turn.
func TestOkMboxFile80(t *testing.T) {
	withOpenMboxReader(t, "TestOkMboxFile80", mboxWith3Messages, func(mr *MboxReader) {
		msg1, err := mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile80: ", err)
			return
		}
		msg2, err := mr.ReadMessage()
		if err == nil {
			t.Error("Expected error here because we haven't finished reading the body of msg1 yet")
			return
		}
		br := msg1.BodyReader()
		bs := make([]byte, 1000)
		err = nil
		for err == nil {
			bs = bs[:cap(bs)]
			n, err := br.Read(bs)
			bs = bs[:n]

			if err == io.EOF {
				break
			} else if err != nil {
				t.Error("TestOkMboxFile80: ", err)
				return
			}
		}

		msg2, err = mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile80: ", err)
			return
		}
		br = msg2.BodyReader()
		bs = make([]byte, 1000)
		err = nil
		for err == nil {
			bs = bs[:cap(bs)]
			n, err := br.Read(bs)
			bs = bs[:n]

			if err == io.EOF {
				break
			} else if err != nil {
				t.Error("TestOkMboxFile80: ", err)
				return
			}
		}

		msg3, err := mr.ReadMessage()
		if err != nil {
			t.Error("TestOkMboxFile80: ", err)
			return
		}
		br = msg3.BodyReader()
		bs = make([]byte, 1000)
		err = nil
		for err == nil {
			bs = bs[:cap(bs)]
			n, err := br.Read(bs)
			bs = bs[:n]

			if err == io.EOF {
				break
			} else if err != nil {
				t.Error("TestOkMboxFile80: ", err)
				return
			}
		}

		_, err = mr.ReadMessage()
		if err != io.EOF {
			t.Error("EOF expected after reading all messages; getting ", err)
			return
		}
	})
}
