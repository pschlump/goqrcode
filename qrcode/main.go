package main

// Copyright 2014 Tom Harwood
// Copyright 2019 Philip Schlump

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
	svgArt := flag.Bool("svg", false, "print as text-art on stdout")
	negative := flag.Bool("i", false, "invert black and white")
	level := flag.String("l", "h", "Level of errro redundancey, h|m|l for high, medium, low")
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

Options:
	-o File			Output file name (will add .png if you do not specify it)
	-s ImageSize	Pixel size of image, 256 is the default.
	-txt				Paint as text art in output.
	-svg				Paint as svg in output.
	-i				Inert color - swap black and white
	-l Level		Level is h, m, l for high, medium, low levesl of qr error correction
	
`)

	}
	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		goqrcode.CheckError(fmt.Errorf("Error: no text to encode into image"))
	}

	out := os.Stdout
	if *outFile != "" {
		var err error
		out, err = filelib.Fopen(*outFile, "w")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to open %s for output: %s\n", *outFile, err)
			os.Exit(1)
		}
		defer out.Close()
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
		// TODO if *negative {
		art := q.ToString(*negative)
		fmt.Fprintf(out, art)
		return
	}

	if *svgArt {
		// TODO -------------------------------------------------------------------------------
		// TODO if *negative {
		genqr(content, out)
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

	out.Write(png)
}

func genqr(uu string, fp *os.File) {
	s := svg.New(fp)

	// Create the QR Code in SVG
	qrCode, _ := qr.Encode(uu, qr.M, qr.Auto)

	// Write QR code to SVG
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()
}
