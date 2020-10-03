package main

type color string

const (
	auto   color = "auto"
	always       = "always"
	never        = "never"
)

func isColorEnabled(color color, terminal bool) bool {
	return color == always || terminal && color == auto
}
