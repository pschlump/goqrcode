package main

// Copyright 2014 Tom Harwood
// Copyright 2019 Philip Schlump

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gitlab.com/pschlump/goqrcode"
)

func main() {

	outFile := flag.String("o", "", "out PNG file prefix, empty for stdout")
	size := flag.Int("s", 256, "image size (pixel)")
	textArt := flag.Bool("t", false, "print as text-art on stdout")
	negative := flag.Bool("i", false, "invert black and white")
	level := flag.String("l", "h", "Level of errro redundancey, h|m|l for high, medium, low")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `qrcode -- QR Code encoder in Go
https://gitlab.com/pschlump/goqrcode

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Usage: After the options the "string" that is passed is encoded into the qr-code.

Example:
       qrcode -o id_1234467721.png "http://www.mysite.com/qr-display/1234467721" 

Options:
	-o File			Output file name (will add .png if you do not specify it)
	-s ImageSize	Pixel size of image, 256 is the default.
	-t				Paint as text art on stdout.
	-i				Inert color - swap black and white
	-l Level		Level is h, m, l for high, medium, low levesl of qr error correction
	
`)
		// xyzzy - add SVG
		// xyzzy - add in image merge, paint one image over QR (for PNG, SVG) - must be able to position
		// xyzzy - add set of border size on QR
		// xyzzy - add of forground / background color

	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		goqrcode.CheckError(fmt.Errorf("Error: no text to encode into image"))
	}

	content := strings.Join(flag.Args(), " ")

	redundancy := goqrcode.Highest
	switch *level {
	case "h", "high", "H":
	case "m", "medium", "M":
		redundancy = goqrcode.Medium
	case "l", "low", "L":
		redundancy = goqrcode.Low
	default:
		flag.Usage()
		os.Exit(1)
	}

	// Generate the QR code in internal format
	var err error
	var q *goqrcode.QRCode
	q, err = goqrcode.New(content, redundancy)
	goqrcode.CheckError(err)

	if *textArt {
		// Output the QR Code as a string
		art := q.ToString(*negative)
		fmt.Println(art)
		return
	}

	// Swap colors
	if *negative {
		q.ForegroundColor, q.BackgroundColor = q.BackgroundColor, q.ForegroundColor
	}

	// Output QR Code as a PNG
	var png []byte
	png, err = q.PNG(*size)
	goqrcode.CheckError(err)

	if *outFile == "" {
		os.Stdout.Write(png)
	} else {
		var fh *os.File
		s := *outFile
		if !strings.HasSuffix(s, ".png") {
			s += ".png"
		}
		fh, err = os.Create(s)
		goqrcode.CheckError(err)
		defer fh.Close()
		fh.Write(png)
	}
}
