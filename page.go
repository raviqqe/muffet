package main

import (
	"net/url"
)

type page interface {
	URL() *url.URL
	Fragments() map[string]struct{}
	Links() map[string]error
}
