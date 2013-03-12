// vim: ts=8 noexpandtab ai

package mbox


import "strings"
import "testing"


const mboxWith1Message = `From foo@bar.com
Subject: Hello world

Test message
`

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
		t.Error(msg)
	}
}


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

// Given a valid mbox file
// When I try to open the file
// Then I expect no error and a valid MboxReader instance.
func TestMboxFile10(t *testing.T) {
	var mr *MboxReader

	expectNoError(t, mboxWith1Message, "Mbox file with one valid message should not yield an error.", &mr)

	if mr == nil {
		t.Error("Returned MboxReader is nil for some reason")
	}
}

