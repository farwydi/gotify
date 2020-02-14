package gotify

type Line []interface{}
type TextElement []byte
type B TextElement
type CODE TextElement

func C(text ...interface{}) (c Line) {
	c = make(Line, len(text))
	for i, tx := range text {
		c[i] = tx
	}
	return c
}

type Formatter interface {
	Format(text ...Line) []byte
}
