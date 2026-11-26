package main

import (
	"errors"
	"strings"
	"testing"
)

func TestNewCSVPageResult(t *testing.T) {
	r := &pageResult{
		URL: "http://foo.com",
		SuccessLinkResults: []*successLinkResult{
			{URL: "http://foo.com/success", StatusCode: 200},
		},
		ErrorLinkResults: []*errorLinkResult{
			{URL: "http://foo.com/error", Error: errors.New("404")},
		},
	}

	t.Run("verbose mode", func(t *testing.T) {
		result := newCSVPageResult(r, true)

		if result.URL != "http://foo.com" {
			t.Errorf("expected URL to be 'http://foo.com', got '%s'", result.URL)
		}

		if len(result.Links) != 2 {
			t.Errorf("expected 2 links, got %d", len(result.Links))
		}

		successLink := result.Links[0]
		if successLink.URL != "http://foo.com/success" {
			t.Errorf("expected success URL to be 'http://foo.com/success', got '%s'", successLink.URL)
		}
		if successLink.Status != "200" {
			t.Errorf("expected status to be '200', got '%s'", successLink.Status)
		}

		errorLink := result.Links[1]
		if errorLink.URL != "http://foo.com/error" {
			t.Errorf("expected error URL to be 'http://foo.com/error', got '%s'", errorLink.URL)
		}
		if errorLink.Status != "404" {
			t.Errorf("expected status to be '404', got '%s'", errorLink.Status)
		}
	})

	t.Run("non-verbose mode", func(t *testing.T) {
		result := newCSVPageResult(r, false)

		if len(result.Links) != 1 {
			t.Errorf("expected 1 link, got %d", len(result.Links))
		}

		errorLink := result.Links[0]
		if errorLink.URL != "http://foo.com/error" {
			t.Errorf("expected error URL to be 'http://foo.com/error', got '%s'", errorLink.URL)
		}
	})
}

func TestCSVPageResultString(t *testing.T) {
	result := &csvPageResult{
		URL: "http://foo.com",
		Links: []csvLinkResult{
			{URL: "http://foo.com/success", Status: "200"},
			{URL: "http://foo.com/error", Status: "404"},
		},
	}

	output := result.String()
	lines := strings.Split(strings.TrimSpace(output), "\n")

	if len(lines) != 3 {
		t.Errorf("expected 3 lines (header + 2 data), got %d", len(lines))
	}

	expectedHeader := `"Page URL","Link URL",Status`
	if lines[0] != expectedHeader {
		t.Errorf("expected header '%s', got '%s'", expectedHeader, lines[0])
	}

	expectedRow1 := `"http://foo.com","http://foo.com/success",200`
	if lines[1] != expectedRow1 {
		t.Errorf("expected first row '%s', got '%s'", expectedRow1, lines[1])
	}

	expectedRow2 := `"http://foo.com","http://foo.com/error",404`
	if lines[2] != expectedRow2 {
		t.Errorf("expected second row '%s', got '%s'", expectedRow2, lines[2])
	}
}
