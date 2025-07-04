package main

type ignoreTimeoutsGroup int

const (
	ignoreTimeoutsGroupNone ignoreTimeoutsGroup = iota
	ignoreTimeoutsGroupAll
	ignoreTimeoutsGroupExternal
)
