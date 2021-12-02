package main

import (
	"errors"
	"fmt"
	"net/url"
)

type redirectHttpClient struct {
	client          httpClient
	maxRedirections int
}

func newRedirectHttpClient(c httpClient, maxRedirections int) httpClient {
	return &redirectHttpClient{c, maxRedirections}
}

func (c *redirectHttpClient) Get(u *url.URL) (httpResponse, error) {
	i := 0

	for {
		r, err := c.client.Get(u)
		if err != nil && i > 0 {
			return nil, fmt.Errorf("%w (following redirect %v)", err, u.String())
		} else if err != nil {
			return nil, err
		}

		switch r.StatusCode() / 100 {
		case 2:
			return r, nil
		case 3:
			i++

			if i > c.maxRedirections {
				return nil, errors.New("too many redirections")
			}

			s := r.Header("Location")

			if len(s) == 0 {
				return nil, errors.New("location header not found")
			}

			u, err = url.Parse(s)

			if err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("%v", r.StatusCode())
		}
	}
}
