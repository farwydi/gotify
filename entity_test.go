package gotify

import (
	"testing"
)

func BenchmarkRangeLoop(t *testing.B) {
	f := func(text ...interface{}) (c Line) {
		c = make(Line, len(text))
		for i, tx := range text {
			c[i] = tx
		}
		return c
	}
	for i := 0; i < t.N; i++ {
		f("x", 1, 5, "5")
	}
}

func BenchmarkCopy(t *testing.B) {
	f := func(text ...interface{}) (c Line) {
		c = make(Line, len(text))
		copy(c, text)
		return c
	}
	for i := 0; i < t.N; i++ {
		f("x", 1, 5, "5")
	}
}
