package utils

import (
	"context"
	"io"
	"log"

	"github.com/a-h/templ"
)

// Writes an unescaped string to a templ template (string must be from a
// trusted source)
func Unsafe(s string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) (err error) {
		_, err = io.WriteString(w, s)
		if err != nil {
			log.Fatal("Failed to write unescaped string to templ template:", err)
		}
		return
	})
}
