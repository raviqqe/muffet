package main

import (
	"errors"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

func TestPageResultFormatterFormatEmptyResult(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false).Format(
			&pageResult{"http://foo.com", nil, nil},
		),
	)
}

func TestPageResultFormatterFormatSuccessLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200},
				},
				nil,
			},
		),
	)
}

func TestPageResultFormatterFormatErrorLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500")},
				},
			},
		),
	)
}

func TestPageResultFormatterFormatSuccessLinkResultsVerbosely(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200},
				},
				nil,
			},
		),
	)
}

func TestPageResultFormatterFormatErrorLinkResultsVerbosely(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200},
				},
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500")},
				},
			},
		),
	)
}

func TestPageResultFormatterSortSuccessLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(true).Format(
			&pageResult{
				"http://foo.com",
				[]*successLinkResult{
					{"http://foo.com", 200},
					{"http://bar.com", 200},
				},
				nil,
			},
		),
	)
}

func TestPageResultFormatterSortormatErrorLinkResults(t *testing.T) {
	cupaloy.SnapshotT(t,
		newPageResultFormatter(false).Format(
			&pageResult{
				"http://foo.com",
				nil,
				[]*errorLinkResult{
					{"http://foo.com", errors.New("500")},
					{"http://bar.com", errors.New("500")},
				},
			},
		),
	)
}
