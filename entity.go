package gotify

type Line []interface{}
type TextElement []byte
type B TextElement
type CODE TextElement

// Concatenation line
func C(text ...interface{}) (c Line) {
	c = make(Line, len(text))
	copy(c, text)
	return c
}
