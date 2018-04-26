package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const rootURL = "http://localhost:8080"
const existentURL = "http://localhost:8080/foo"
const nonExistentURL = "http://localhost:8080/bar"
const erroneousURL = "http://localhost:8080/erroneous"
const fragmentURL = "http://localhost:8080/fragment"
const tagsURL = "http://localhost:8080/tags"
const nonExistentIDURL = "http://localhost:8080/#non-existent-id"
const baseURL = "http://localhost:8080/base"
const invalidBaseURL = "http://localhost:8080/invalid-base"
const redirectURL = "http://localhost:8080/redirect"
const invalidRedirectURL = "http://localhost:8080/invalid-redirect"

type handler struct{}

func (handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
		w.Write([]byte(htmlWithBody(`<a href="/foo" />`)))
	case "/foo":
		w.Write([]byte(htmlWithBody(`<a href="/" />`)))
	case "/erroneous":
		w.Write([]byte(htmlWithBody(`
			<a href=":" />
			<a href="mailto:me@right.here" />
			<a href="/bar" />
			<a href="#foo" />
		`)))
	case "/fragment":
		w.Write([]byte(htmlWithBody(`<a id="foo" href="#foo" />`)))
	case "/tags":
		// TODO: Test <frame> tag.
		w.Write([]byte(htmlWithBody(`
			<a href="/" />
			<iframe src="/"></iframe>
			<img src="/foo.png" />
			<link href="/" />
			<script src="/foo.js"></script>
			<source src="/foo.png" />
			<track src="/foo.vtt" />
		`)))
	case "/foo.js":
		w.Write(nil)
	case "/foo.png":
		w.Write(nil)
	case "/foo.vtt":
		w.Write(nil)
	case "/base":
		w.Write([]byte(`
			<html>
				<head>
					<base href="/parent/" />
				</head>
				<body>
					<a href="child" />
				</body>
			</html>
		`))
	case "/invalid-base":
		w.Write([]byte(`
			<html>
				<head>
					<base href=":" />
				</head>
				<body>
					<a href="child" />
				</body>
			</html>
		`))
	case "/parent/child":
		w.Write(nil)
	case "/redirect":
		w.Header().Add("Location", "/")
		w.WriteHeader(300)
	case "/invalid-redirect":
		w.WriteHeader(300)
	default:
		w.WriteHeader(404)
	}
}

func TestMain(m *testing.M) {
	go http.ListenAndServe(":8080", handler{})
	time.Sleep(time.Millisecond)

	os.Exit(m.Run())
}

func htmlWithBody(b string) string {
	return fmt.Sprintf(`<html><body>%v</body></html>`, b)
}
