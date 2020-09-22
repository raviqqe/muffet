package main

type httpResponse interface {
	URL() string
	StatusCode() int
	Header(string) string
	Body() []byte
}
