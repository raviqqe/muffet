package main

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"
)

const (
	rootURL            = "http://localhost:8080"
	existentURL        = "http://localhost:8080/foo"
	nonExistentURL     = "http://localhost:8080/bar"
	erroneousURL       = "http://localhost:8080/erroneous"
	fragmentURL        = "http://localhost:8080/fragment"
	nonExistentIDURL   = "http://localhost:8080/#non-existent-id"
	baseURL            = "http://localhost:8080/base"
	invalidBaseURL     = "http://localhost:8080/invalid-base"
	redirectURL        = "http://localhost:8080/redirect"
	invalidRedirectURL = "http://localhost:8080/invalid-redirect"
	missingSitemapURL  = "http://localhost:8081"
)

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
	case "/redirect":
		w.Header().Add("Location", "/")
		w.WriteHeader(300)
	case "/invalid-redirect":
		w.WriteHeader(300)
	case "/sitemap.xml":
		w.Write([]byte(fmt.Sprintf(`
			<?xml version="1.0" encoding="UTF-8"?>
			<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
				<url>
					<loc>%v</loc>
					<lastmod>896-5-4</lastmod>
					<changefreq>monthly</changefreq>
					<priority>0.8</priority>
				</url>
				<url>
					<loc>%v</loc>
				</url>
			</urlset>
		`, rootURL, existentURL)))
	default:
		w.WriteHeader(404)
	}
}

type noSitemapHandler struct{}

func (noSitemapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "", "/":
	default:
		w.WriteHeader(404)
	}
}

func TestMain(m *testing.M) {
	go http.ListenAndServe(":8080", handler{})
	go http.ListenAndServe(":8081", noSitemapHandler{})

	time.Sleep(time.Millisecond)

	os.Exit(m.Run())
}

func htmlWithBody(b string) string {
	return fmt.Sprintf(`<html><body>%v</body></html>`, b)
}
