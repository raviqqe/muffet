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

func (c *redirectHttpClient) Get(u *url.URL, headers map[string]string) (httpResponse, error) {
	i := 0

	for {
		r, err := c.client.Get(u, headers)
		if err != nil {
			return nil, c.formatError(err, i, u)
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
				return nil, errors.New("location header not set")
			}

			u, err = u.Parse(s)

			if err != nil {
				return nil, err
			}
		default:
			return nil, c.formatError(fmt.Errorf("%v", r.StatusCode()), i, u)
		}
	}
}

func (*redirectHttpClient) formatError(err error, redirections int, u *url.URL) error {
	if redirections == 0 {
		return err
	}

	return fmt.Errorf("%w (following redirect %v)", err, u.String())
}
