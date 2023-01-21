package main

import (
	"errors"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy"
)

func TestPageResultFormatterFormatEmptyResult(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{"http://foo.com", nil, nil, d},
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
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, d},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), d},
				},
				d,
			},
		),
	)
}

func TestPageResultFormatterFormatSuccessLinkResultsVerbosely(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
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

func TestPageResultFormatterFormatErrorLinkResultsVerbosely(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, d},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), d},
				},
				d,
			},
		),
	)
}

func TestPageResultFormatterSortSuccessLinkResults(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true, true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200, d},
					{"http://bar.com", 200, d},
				},
				nil,
				d,
			},
		),
	)
}

func TestPageResultFormatterSortErrorLinkResults(t *testing.T) {
	d, _ := time.ParseDuration("1s")
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false, true).Format(
			&pageResult{
				"http://foo.com",
				nil,
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500"), d},
					{"http://bar.com", errors.New("500"), d},
				},
				d,
			},
		),
	)
}
