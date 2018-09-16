package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"net/url"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type page struct {
	url   *url.URL
	ids   map[string]struct{}
	links map[string]error
}

func newPage(s string, n *html.Node, sc scraper) (*page, error) {
	u, err := url.Parse(s)

	if err != nil {
		return nil, err
	}

	u.Fragment = ""
	u.RawQuery = ""

	ids := map[string]struct{}{}

	scrape.FindAllNested(n, func(n *html.Node) bool {
		if s := scrape.Attr(n, "id"); s != "" {
			ids[s] = struct{}{}
		}

		return false
	})

	b := u

	if n, ok := scrape.Find(n, func(n *html.Node) bool {
		return n.DataAtom == atom.Base
	}); ok {
		u, err := url.Parse(scrape.Attr(n, "href"))

		if err != nil {
			return nil, err
		}

		b = b.ResolveReference(u)
	}

	return &page{u, ids, sc.Scrape(n, b)}, nil
}

func (p page) URL() *url.URL {
	return p.url
}

func (p page) IDs() map[string]struct{} {
	return p.ids
}

func (p page) Links() map[string]error {
	return p.links
}

type encodablePage struct {
	URL   *url.URL
	IDs   map[string]struct{}
	Links map[string]string
}

func (p page) MarshalBinary() ([]byte, error) {
	b := bytes.NewBuffer(nil)
	ls := make(map[string]string, len(p.links))

	for k, v := range p.links {
		if v == nil {
			ls[k] = ""
			continue
		}

		ls[k] = v.Error()
	}

	if err := gob.NewEncoder(b).Encode(encodablePage{p.url, p.ids, ls}); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (p *page) UnmarshalBinary(bs []byte) error {
	q := encodablePage{}

	if err := gob.NewDecoder(bytes.NewBuffer(bs)).Decode(&q); err != nil {
		return err
	}

	ls := make(map[string]error, len(q.Links))

	for k, v := range q.Links {
		if v == "" {
			ls[k] = nil
			continue
		}

		ls[k] = errors.New(v)
	}

	*p = page{q.URL, q.IDs, ls}

	return nil
}
