package gotify

import (
	"github.com/stretchr/testify/mock"
)

type MockAdapter struct {
	mock.Mock
}

func (ad *MockAdapter) Format(text []Line) []byte {
	args := ad.Called(text)
	return args.Get(0).([]byte)
}

func (ad *MockAdapter) Send(subject string, message ...Line) error {
	args := ad.Called(subject, message)
	return args.Error(0)
}

func (ad *MockAdapter) SendWithoutFormatting(p []byte) error {
	args := ad.Called(p)
	return args.Error(0)
}
