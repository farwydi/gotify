package gotify

type Adapter interface {
	Format(text []Line) []byte
	Send(subject string, message ...Line) error
	SendRaw(p []byte) error
}
