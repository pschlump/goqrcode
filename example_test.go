// go-qrcode
// Copyright 2014 Tom Harwood
// Copyright 2018-2026 Philip Schlump

package goqrcode_test

import (
	"fmt"
	"image/color"
	"log"
	"os"

	goqrcode "github.com/pschlump/goqrcode"
)

// Encode a URL and get the PNG bytes.
func ExampleEncode() {
	png, err := goqrcode.Encode("https://example.org", goqrcode.Medium, 256)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("PNG is %d bytes long\n", len(png))
}

// Write a QR Code directly to a file.
func ExampleWriteFile() {
	if err := goqrcode.WriteFile("https://example.org", goqrcode.Medium, 256,
		"/tmp/qr.png"); err != nil {
		log.Fatal(err)
	}
}

// Write a QR Code with custom colors (white on black) to a file.
func ExampleWriteColorFile() {
	if err := goqrcode.WriteColorFile("https://example.org", goqrcode.Medium, 256,
		color.Black, color.White, "/tmp/qr.png"); err != nil {
		log.Fatal(err)
	}
}

// Draw a QR Code to any io.Writer, here a file.
func ExampleQRCode_Write() {
	q, err := goqrcode.New("https://example.org", goqrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}

	f, err := os.Create("/tmp/qr.png")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	if err := q.Write(256, f); err != nil {
		log.Fatal(err)
	}
}

// Print a QR Code as text art in a terminal.
func ExampleQRCode_ToString() {
	q, err := goqrcode.New("https://example.org", goqrcode.Medium)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(q.ToString(false))
}

// The four error recovery levels. Higher levels produce more robust but
// physically larger QR Codes.
func ExampleRecoveryLevel() {
	for _, level := range []goqrcode.RecoveryLevel{
		goqrcode.Low, goqrcode.Medium, goqrcode.High, goqrcode.Highest,
	} {
		q, err := goqrcode.New("https://example.org", level)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("level %d: version %d, %d modules wide\n",
			level, q.VersionNumber, len(q.Bitmap()))
	}
}
