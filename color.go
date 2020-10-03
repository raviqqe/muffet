package main

type color string

const (
	auto   color = "auto"
	always       = "always"
	never        = "never"
)

func isColorEnabled(c color, terminal bool) bool {
	return c == always || terminal && c == auto
}
