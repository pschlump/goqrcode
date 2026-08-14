package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	params "github.com/pschlump/goqrcode/qr-gen-server/params"
)

func serveOne(t *testing.T, pp params.ApiDataType) *httptest.ResponseRecorder {
	t.Helper()

	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/", nil)
	genqr_bl(recorder, req, pp)
	return recorder
}

func TestGenqrPNG(t *testing.T) {
	recorder := serveOne(t, params.ApiDataType{
		Url: "abc.xyz", Fmt: "png", Redundancy: "M", Size: 256,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "image/png" {
		t.Fatalf("got Content-Type %q, expected %q", got, "image/png")
	}
	// PNG magic number.
	if !bytes.HasPrefix(recorder.Body.Bytes(), []byte{0x89, 'P', 'N', 'G'}) {
		t.Fatalf("response does not start with the PNG magic number")
	}
}

func TestGenqrSVG(t *testing.T) {
	recorder := serveOne(t, params.ApiDataType{
		Url: "abc.xyz", Fmt: "svg", Redundancy: "M", Size: 256,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "image/svg+xml" {
		t.Fatalf("got Content-Type %q, expected %q", got, "image/svg+xml")
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte("<svg")) {
		t.Fatalf("response does not contain an <svg> element")
	}
}

func TestGenqrText(t *testing.T) {
	recorder := serveOne(t, params.ApiDataType{
		Url: "abc.xyz", Fmt: "text", Redundancy: "M", Size: 256,
	})

	if recorder.Code != http.StatusOK {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusOK)
	}
	if got := recorder.Header().Get("Content-Type"); got != "text/plain; charset=utf-8" {
		t.Fatalf("got Content-Type %q, expected %q", got, "text/plain; charset=utf-8")
	}
	if recorder.Body.Len() == 0 {
		t.Fatal("text response is empty")
	}
}

func TestGenqrUnsupportedFormat(t *testing.T) {
	recorder := serveOne(t, params.ApiDataType{
		Url: "abc.xyz", Fmt: "gif", Redundancy: "M", Size: 256,
	})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGenqrContentTooLong(t *testing.T) {
	recorder := serveOne(t, params.ApiDataType{
		Url: string(bytes.Repeat([]byte("#"), 4000)), Fmt: "png", Redundancy: "M", Size: 256,
	})

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusBadRequest)
	}
}

func TestGenqrHandlerEndToEnd(t *testing.T) {
	recorder := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "http://example.com/?url=abc.xyz&fmt=text&redundancy=L", nil)

	genqr(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("got status %d, expected %d", recorder.Code, http.StatusOK)
	}
	if recorder.Body.Len() == 0 {
		t.Fatal("response is empty")
	}
}
