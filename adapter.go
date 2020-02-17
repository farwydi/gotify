package gotify

type Adapter interface {
	Format(text []Line) []byte
	Send(subject string, message ...Line) error
	SendWithoutFormatting(p []byte) error
}
