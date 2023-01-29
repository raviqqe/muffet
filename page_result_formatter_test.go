package main

import (
	"errors"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
)

func TestPageResultFormatterFormatEmptyResult(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{"http://foo.com", nil, nil, 0},
		),
	)
}

func TestPageResultFormatterFormatSuccessLinkResults(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, d},
				},
				nil,
				d,
			},
		),
	)
}

func TestPageResultFormatterFormatErrorLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, 0},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), 0},
				},
				0,
			},
		),
	)
}

func TestPageResultFormatterFormatSuccessLinkResultsVerbosely(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, 0},
				},
				nil,
				0,
			},
		),
	)
}

func TestPageResultFormatterFormatErrorLinkResultsVerbosely(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, 0},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), 0},
				},
				0,
			},
		),
	)
}

func TestPageResultFormatterSortSuccessLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, 0},
					{"http://bar.com", 200, 0},
				},
				nil,
				0,
			},
		),
	)
}

func TestPageResultFormatterSortErrorLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{
				"http://foo.com",
				nil,
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), 0},
					{"http://bar.com", errors.New("500"), 0},
				},
				0,
			},
		),
	)
}
