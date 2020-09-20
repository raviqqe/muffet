package main

import (
	"errors"
	"net/url"
	"time"

	"github.com/temoto/robotstxt"
	"github.com/yterajima/go-sitemap"
)

type urlValidator struct {
	hostname     string
	includedURLs map[string]struct{}
	robotsTxt    *robotstxt.RobotsData
}

func newURLValidator(c httpClient, s string, useRobotsTxt, useSitemap bool) (urlValidator, error) {
	u, err := url.Parse(s)
	if err != nil {
		return urlValidator{}, err
	}

	rd := (*robotstxt.RobotsData)(nil)

	if useRobotsTxt {
		u.Path = "robots.txt"
		r, err := c.Get(u, nil, time.Duration(0))

		if err != nil {
			return urlValidator{}, err
		} else if r.StatusCode() != 200 {
			return urlValidator{}, errors.New("robots.txt not found")
		}

		rd, err = robotstxt.FromBytes(r.Body())

		if err != nil {
			return urlValidator{}, err
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
			return urlValidator{}, err
		}

		for _, u := range m.URL {
			us[u.Loc] = struct{}{}
		}
	}

	return urlValidator{u.Hostname(), us, rd}, nil
}

func (i urlValidator) Validate(u *url.URL) bool {
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
