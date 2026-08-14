// Copyright 2026 Philip Schlump

package goqrcode

import (
	"bytes"
	"image/color"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

func TestToString(t *testing.T) {
	q, err := New("hello", Medium)
	if err != nil {
		t.Fatalf("New: %s", err)
	}

	art := q.ToString(false)
	lines := strings.Split(strings.TrimSuffix(art, "\n"), "\n")

	// The bitmap is square: one line per row, two characters per module.
	bitmap := q.Bitmap()
	if len(lines) != len(bitmap) {
		t.Fatalf("ToString produced %d lines, expected %d", len(lines), len(bitmap))
	}
	for i, line := range lines {
		if got := utf8.RuneCountInString(line); got != 2*len(bitmap[i]) {
			t.Fatalf("line %d is %d characters, expected %d", i, got, 2*len(bitmap[i]))
		}
	}

	// Inverting swaps the block characters for spaces.
	inverse := q.ToString(true)
	if inverse == art {
		t.Fatal("ToString(true) identical to ToString(false), expected inverted output")
	}
}

func TestToSmallString(t *testing.T) {
	q, err := New("hello", Medium)
	if err != nil {
		t.Fatalf("New: %s", err)
	}

	art := q.ToSmallString(false)
	lines := strings.Split(strings.TrimSuffix(art, "\n"), "\n")

	// Each output line covers two bitmap rows, so the count is the bitmap size
	// divided by two, rounded up.
	bitmap := q.Bitmap()
	expectedLines := (len(bitmap) + 1) / 2
	if len(lines) != expectedLines {
		t.Fatalf("ToSmallString produced %d lines, expected %d", len(lines), expectedLines)
	}

	// Each output column covers one module (half the width of ToString). The
	// block characters are multi-byte, so count runes, not bytes.
	for i, line := range lines {
		if got := utf8.RuneCountInString(line); got != len(bitmap[0]) {
			t.Fatalf("line %d is %d characters, expected %d", i, got, len(bitmap[0]))
		}
	}
}

func TestEncodeTooLong(t *testing.T) {
	_, err := Encode(strings.Repeat("#", 2954), Low, 256)
	if err == nil {
		t.Fatal("encoding 2954 bytes succeeded, expected failure")
	}
}

// WriteColorFile with content too long for a QR Code must return an error,
// not panic (regression test: the color fields used to be set before the
// error check, dereferencing a nil QRCode).
func TestWriteColorFileTooLong(t *testing.T) {
	err := WriteColorFile(strings.Repeat("#", 2954), Low, 256,
		color.White, color.Black, t.TempDir()+"/qr.png")
	if err == nil {
		t.Fatal("WriteColorFile with too-long content succeeded, expected failure")
	}
}

func TestImageVariableSize(t *testing.T) {
	q, err := New("https://example.org", Medium)
	if err != nil {
		t.Fatalf("New: %s", err)
	}

	// A negative size requests a variable sized image: -5 means 5px per
	// module, so the result must be exactly 5 * symbol size pixels square,
	// with no extra padding.
	img := q.Image(-5)
	bounds := img.Bounds()
	if bounds.Dx() != bounds.Dy() {
		t.Fatalf("image is %dx%d, expected square", bounds.Dx(), bounds.Dy())
	}
	if expected := 5 * len(q.Bitmap()); bounds.Dx() != expected {
		t.Fatalf("image is %dpx wide, expected %dpx", bounds.Dx(), expected)
	}
}

func TestWriteToWriter(t *testing.T) {
	q, err := New("https://example.org", Medium)
	if err != nil {
		t.Fatalf("New: %s", err)
	}

	var buf bytes.Buffer
	if err := q.Write(256, &buf); err != nil {
		t.Fatalf("Write: %s", err)
	}

	png, err := q.PNG(256)
	if err != nil {
		t.Fatalf("PNG: %s", err)
	}

	if !bytes.Equal(buf.Bytes(), png) {
		t.Fatalf("Write produced %d bytes, PNG produced %d bytes, expected identical",
			buf.Len(), len(png))
	}
}

func TestWriteFileRoundTrip(t *testing.T) {
	filename := t.TempDir() + "/qr.png"
	if err := WriteFile("https://example.org", Medium, 256, filename); err != nil {
		t.Fatalf("WriteFile: %s", err)
	}

	png, err := Encode("https://example.org", Medium, 256)
	if err != nil {
		t.Fatalf("Encode: %s", err)
	}

	onDisk, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading %s: %s", filename, err)
	}

	if !bytes.Equal(onDisk, png) {
		t.Fatalf("file differs from Encode output")
	}
}
