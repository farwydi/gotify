package gotify

import (
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	testSubject := "Hello"
	testMessage := []Line{
		C("Hello world"),
		C(B("bold"), " text"),
	}

	client, err := NewClient(
		func() (adapter Adapter, err error) {
			m := &MockAdapter{}
			m.
				On("Send", testSubject, testMessage).
				Return(nil).
				Run(func(args mock.Arguments) {
					subject := args.Get(0).(string)
					require.Equal(t, testSubject, subject)

					message := args.Get(1).([]Line)
					require.Equal(t, testMessage, message)
				})
			return m, nil
		},
	)
	require.NoError(t, err)

	errs := client.Send(testSubject, testMessage...)
	for _, err := range errs {
		require.NoError(t, err)
	}
}
