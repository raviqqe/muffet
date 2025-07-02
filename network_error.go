package main

type networkErrorGroup int

const (
	networkErrorGroupNone networkErrorGroup = iota
	networkErrorGroupAll
	networkErrorGroupExternal
)
