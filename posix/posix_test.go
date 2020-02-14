package posix

import (
	"github.com/farwydi/gotify"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPosix(t *testing.T) {
	init := NewPosixAdapter("1", []string{"1"}, "1")
	adapter, err := init()
	require.NoError(t, err)

	message := []gotify.Line{
		gotify.C("Hello world"),
		gotify.C(gotify.B("bold"), " text"),
		gotify.C(gotify.CODE("code here")),
	}

	require.Equal(t,
		[]byte("Hello world</br><b>bold</b> text</br><pre>code here</pre></br>"),
		adapter.Format(message),
	)
}
