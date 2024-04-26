package utils

import (
	"context"
	"io"
	"log"
	"net/url"

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

func UrlHostname(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		log.Fatalf("Failed to parse URL '%v' for it's hostname: %v", urlStr, err)
	}
	if !u.IsAbs() {
		log.Fatalf("Failed to find absolute URL whilst attempting to extract hostname for URL '%v': %v", urlStr, err)
	}
	return u.Hostname()
}
