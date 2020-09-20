package main

import (
	"errors"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
	"github.com/yterajima/go-sitemap"
)

type urlInspector struct {
	hostname     string
	includedURLs map[string]struct{}
	robotsTxt    *robotstxt.RobotsData
}

func newURLInspector(c httpClient, s string, useRobotsTxt, useSitemap bool) (urlInspector, error) {
	u, err := url.Parse(s)
	if err != nil {
		return urlInspector{}, err
	}

	rd := (*robotstxt.RobotsData)(nil)

	if useRobotsTxt {
		u.Path = "robots.txt"
		r, err := c.Get(u, nil, time.Duration(0))

		if err != nil {
			return urlInspector{}, err
		} else if r.StatusCode() != 200 {
			return urlInspector{}, errors.New("robots.txt not found")
		}

		rd, err = robotstxt.FromBytes(r.Body())

		if err != nil {
			return urlInspector{}, err
		}
	}

	us := map[string]struct{}{}

	if useSitemap {
		u.Path = "sitemap.xml"

		sitemap.SetFetch(func(s string, _ interface{}) ([]byte, error) {
			u, err := url.Parse(s)

			if err != nil {
				return nil, err
			}

			r, err := c.Get(u, nil, time.Duration(0))

			if err != nil {
				return nil, err
			} else if r.StatusCode() != 200 {
				return nil, errors.New("sitemap not found")
			}

			return r.Body(), err
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
