package ApiQRGen

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParseGETParams(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected ApiDataType
	}{
		{
			name: "defaults",
			url:  "http://example.com/",
			expected: ApiDataType{
				Fmt: "png", Invert: false, Redundancy: "M", Size: 256,
			},
		},
		{
			name: "all parameters",
			url:  "http://example.com/?url=abc.xyz&fmt=svg&invert=true&redundancy=H&size=512",
			expected: ApiDataType{
				Url: "abc.xyz", Fmt: "svg", Invert: true, Redundancy: "H", Size: 512,
			},
		},
		{
			name: "fmt text",
			url:  "http://example.com/?url=abc.xyz&fmt=text",
			expected: ApiDataType{
				Url: "abc.xyz", Fmt: "text", Invert: false, Redundancy: "M", Size: 256,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			recorder := httptest.NewRecorder()

			pp, err := ParseParams(recorder, req, "GET", "POST")
			if err != nil {
				t.Fatalf("ParseParams: %s", err)
			}
			if pp != test.expected {
				t.Fatalf("got %+v, expected %+v", pp, test.expected)
			}
		})
	}
}

func TestParseGETParamsInvalid(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{"bad invert", "http://example.com/?invert=not-a-bool"},
		{"bad size", "http://example.com/?size=not-a-number"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", test.url, nil)
			recorder := httptest.NewRecorder()

			_, err := ParseParams(recorder, req, "GET", "POST")
			if err == nil {
				t.Fatal("ParseParams succeeded, expected failure")
			}
			if recorder.Code != http.StatusBadRequest {
				t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestParsePOSTParams(t *testing.T) {
	req := httptest.NewRequest("POST", "http://example.com/",
		strings.NewReader("url=abc.xyz&fmt=svg&invert=1&redundancy=L&size=128"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	recorder := httptest.NewRecorder()

	pp, err := ParseParams(recorder, req, "GET", "POST")
	if err != nil {
		t.Fatalf("ParseParams: %s", err)
	}

	expected := ApiDataType{
		Url: "abc.xyz", Fmt: "svg", Invert: true, Redundancy: "L", Size: 128,
	}
	if pp != expected {
		t.Fatalf("got %+v, expected %+v", pp, expected)
	}
}

func TestParseParamsInvalidMethod(t *testing.T) {
	req := httptest.NewRequest("DELETE", "http://example.com/", nil)
	recorder := httptest.NewRecorder()

	_, err := ParseParams(recorder, req, "GET", "POST")
	if err == nil {
		t.Fatal("ParseParams succeeded, expected failure")
	}
	if recorder.Code != http.StatusMethodNotAllowed {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusMethodNotAllowed)
	}
}

func TestParseBool(t *testing.T) {
	trueCases := []string{"t", "T", "yes", "YES", "1", "true", "TRUE", "on", "ON"}
	for _, s := range trueCases {
		b, err := ParseBool(s)
		if err != nil {
			t.Fatalf("ParseBool(%q): %s", s, err)
		}
		if !b {
			t.Fatalf("ParseBool(%q) = false, expected true", s)
		}
	}

	falseCases := []string{"f", "F", "no", "NO", "0", "false", "FALSE", "off", "OFF"}
	for _, s := range falseCases {
		b, err := ParseBool(s)
		if err != nil {
			t.Fatalf("ParseBool(%q): %s", s, err)
		}
		if b {
			t.Fatalf("ParseBool(%q) = true, expected false", s)
		}
	}

	if _, err := ParseBool("not-a-bool"); err == nil {
		t.Fatal(`ParseBool("not-a-bool") succeeded, expected failure`)
	}
}
