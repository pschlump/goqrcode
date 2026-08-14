package main

// Copyright 2014 Tom Harwood
// Copyright 2019-2026 Philip Schlump

import (
	"flag"
	"fmt"
	"os"
	"strings"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"github.com/pschlump/filelib"
	"github.com/pschlump/goqrcode"
	"github.com/pschlump/goqrsvg"
)

func main() {

	outFile := flag.String("o", "", "out PNG file prefix, empty for stdout")
	size := flag.Int("s", 256, "image size (pixel)")
	textArt := flag.Bool("txt", false, "print as text-art on stdout")
	svgArt := flag.Bool("svg", false, "print as SVG on stdout")
	negative := flag.Bool("i", false, "invert black and white")
	level := flag.String("l", "h", "Level of error redundancy, h|m|l for high, medium, low")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `qrcode -- QR Code encoder in Go
https://github.com/pschlump/goqrcode

Flags:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Usage: After the options the "string" that is passed is encoded into the qr-code.

Example:
       qrcode -o id_1234467721.png "http://www.mysite.com/qr-display/1234467721"

       qrcode -o xyz.svg -svg "http://www.mysite.com/qr-display/1234467721"

Options:
	-o File			Output file name (will add .png if you do not specify it)
	-s ImageSize	Pixel size of image, 256 is the default.
	-txt			Paint as text art in output.
	-svg			Paint as svg in output.
	-i				Invert color - swap black and white
	-l Level		Level is h, m, l for high, medium, low levels of qr error correction

`)
	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		goqrcode.CheckError(fmt.Errorf("no text to encode into image"))
	}

	out := os.Stdout
	if *outFile != "" {
		var err error
		out, err = filelib.Fopen(*outFile, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", *outFile, err)
			os.Exit(1)
		}
		defer func() {
			if err := out.Close(); err != nil {
				fmt.Fprintf(os.Stderr, "Unable to close %s: %s\n", *outFile, err)
			}
		}()
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
		if _, err := fmt.Fprint(out, art); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to write output: %s\n", err)
			os.Exit(1)
		}
		return
	}

	if *svgArt {
		genqr(content, *negative, out)
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

	if _, err := out.Write(png); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write output: %s\n", err)
		os.Exit(1)
	}
}

func genqr(uu string, neg bool, fp *os.File) {
	s := svg.New(fp)

	// Create the QR Code in SVG
	qrCode, err := qr.Encode(uu, qr.M, qr.Auto)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to encode %q: %s\n", uu, err)
		os.Exit(1)
	}

	// Write QR code to SVG
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	if neg {
		err = qs.WriteQrSVGInverse(s)
	} else {
		err = qs.WriteQrSVG(s)
	}
	goqrcode.CheckError(err)

	s.End()
}
