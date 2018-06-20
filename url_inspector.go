package main

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/temoto/robotstxt"
	"github.com/yterajima/go-sitemap"
)

type urlInspector struct {
	hostname     string
	includedURLs map[string]struct{}
	robotsTxt    *robotstxt.RobotsData
}

func newURLInspector(s string, r, sm bool) (urlInspector, error) {
	u, err := urlParse(s)

	if err != nil {
		return urlInspector{}, err
	}

	rd := (*robotstxt.RobotsData)(nil)

	if r {
		u.Path = "robots.txt"
		res, err := http.Get(u.String())

		if err != nil {
			return urlInspector{}, err
		} else if res.StatusCode != 200 {
			return urlInspector{}, errors.New("robots.txt not found")
		}

		rd, err = robotstxt.FromResponse(res)

		if err != nil {
			return urlInspector{}, err
		}
	}

	us := map[string]struct{}{}

	if sm {
		u.Path = "sitemap.xml"
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

// unescape and remove any embedded tabs and CR/LF characters
func urlParse(s string) (*url.URL, error) {
	s, err := url.PathUnescape(s)
	if err != nil {
		return nil, err
	}
	s = strings.Replace(s, "\t", "", -1)
	s = strings.Replace(s, "\r", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return url.Parse(s)
}
