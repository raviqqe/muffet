package main

type page struct {
	url  string
	body []byte
}

func newPage(u string, b []byte) page {
	return page{u, b}
}

func (p page) URL() string {
	return p.url
}

func (p page) Body() []byte {
	return p.body
}
