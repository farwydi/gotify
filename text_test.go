package gotify

import (
	"bytes"
	"testing"
)

func TestTelegram(t *testing.T) {
	f := telegramAdapter{}

	message := []Line{
		C("Hello *world*"),
		C(15),
		C(B("bold"), " text"),
		C(CODE("code here")),
	}

	ft := f.Format(message)
	if !bytes.Equal(ft, []byte("Hello \\*world\\*\n15\n*bold* text\n`code here`\n")) {
		t.Fail()
	}
}

func TestPosix(t *testing.T) {
	f := posixAdapter{}

	message := []Line{
		C("Hello world"),
		C(B("bold"), " text"),
		C(CODE("code here")),
	}

	if !bytes.Equal(f.Format(message), []byte("Hello world</br><b>bold</b> text</br><pre>code here</pre></br>")) {
		t.Fail()
	}
}
