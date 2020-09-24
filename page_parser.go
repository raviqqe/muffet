package main

import (
	"bytes"
	"net/url"

	"github.com/yhat/scrape"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

type pageParser struct {
	linkFinder linkFinder
}

func newPageParser(f linkFinder) *pageParser {
	return &pageParser{f}
}

func (p pageParser) Parse(rawURL string, body []byte) (*page, error) {
	n, err := html.Parse(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	u.Fragment = ""

	frs := map[string]struct{}{}

	scrape.FindAllNested(n, func(n *html.Node) bool {
		for _, a := range []string{"id", "name"} {
			if s := scrape.Attr(n, a); s != "" {
				frs[s] = struct{}{}
			}
		}

		return false
	})

	base := u

	if n, ok := scrape.Find(n, func(n *html.Node) bool {
		return n.DataAtom == atom.Base
	}); ok {
		u, err := url.Parse(scrape.Attr(n, "href"))
		if err != nil {
			return nil, err
		}

		base = base.ResolveReference(u)
	}

	return newPage(u, frs, p.linkFinder.Find(n, base)), nil
}
