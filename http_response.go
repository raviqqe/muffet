package main

type httpResponse interface {
	URL() string
	StatusCode() int
	Header(string) string
	Headers(string) []string
	Body() ([]byte, error)
}
