package main

type pageParser interface {
	// Returned pages can be nil to indicate unrecognized file formats.
	Parse(string, string, []byte) (page, error)
}
