package gotify

type adapter interface {
	send(subject string, message ...Line) error
}
