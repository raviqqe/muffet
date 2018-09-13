package main

import (
	"errors"
	"net/url"

	"github.com/temoto/robotstxt"
	"github.com/valyala/fasthttp"
	"github.com/yterajima/go-sitemap"
)

type urlInspector struct {
	hostname     string
	includedURLs map[string]struct{}
	robotsTxt    *robotstxt.RobotsData
}

func newURLInspector(c *fasthttp.Client, s string, r, sm bool) (urlInspector, error) {
	u, err := url.Parse(s)

	if err != nil {
		return urlInspector{}, err
	}

	rd := (*robotstxt.RobotsData)(nil)

	if r {
		u.Path = "robots.txt"
		c, bs, err := c.Get(nil, u.String())

		if err != nil {
			return urlInspector{}, err
		} else if c != 200 {
			return urlInspector{}, errors.New("robots.txt not found")
		}

		rd, err = robotstxt.FromBytes(bs)

		if err != nil {
			return urlInspector{}, err
		}
	}

	us := map[string]struct{}{}

	if sm {
		u.Path = "sitemap.xml"

		sitemap.SetFetch(func(s string, _ interface{}) ([]byte, error) {
			_, bs, err := c.Get(nil, s)
			return bs, err
		})

		m, err := sitemap.Get(u.String(), nil)

		if err != nil {
			return urlInspector{}, err
		}

		for _, u := range m.URL {
			us[u.Loc] = struct{}{}
		}
	}

	return urlInspector{u.Hostname(), us, rd}, nil
}

func (i urlInspector) Inspect(u *url.URL) bool {
	if len(i.includedURLs) != 0 {
		if _, ok := i.includedURLs[u.String()]; !ok {
			return false
		}
	}

	if i.robotsTxt != nil && !i.robotsTxt.TestAgent(u.Path, "muffet") {
		return false
	}

	return u.Hostname() == i.hostname
}
