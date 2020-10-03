package main

type color string

const (
	auto   color = "auto"
	always color = "always"
	never  color = "never"
)

func isColorEnabled(c color, terminal bool) bool {
	return c == always || terminal && c == auto
}
