package main

import (
	"fmt"
	"log"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"github.com/pschlump/godebug"
	"github.com/pschlump/goqrcode"
	params "github.com/pschlump/goqrcode/qr-gen-server/params"
	"github.com/pschlump/goqrsvg"
)

func main() {
	http.Handle("/", http.HandlerFunc(genqr))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func genqr(www http.ResponseWriter, req *http.Request) {

	// size, url
	// redundancy H, M, L
	// fmt PNG, SVG, SVG-fragment
	pp, err := params.ParseParams(www, req, "GET", "POST")
	fmt.Printf("%s: %s\n", req.Method, godebug.SVarI(pp))
	if err != nil {
		return
	}

	// xyzzy - if pp.Fmt == "svg"
	// xyzzy - if pp.Fmt == "png"

	if pp.Fmt == "svg" {

		s := svg.New(www)
		www.Header().Set("Content-Type", "image/svg+xml")

		// Create the QR Code in SVG
		// qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)
		qrCode, err := qr.Encode(pp.Url, qr.M, qr.Auto)
		if err != nil {
			fmt.Printf("error: %s\n", err)
			// error processing -- TODO
		}

		if db1 {
			fmt.Printf("at: %s Params: %+v\n", godebug.LF(), pp)
		}

		// Write QR code to SVG
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(s)
		qs.WriteQrSVG(s)

		s.End()

	} else if pp.Fmt == "png" {

		redundancy := goqrcode.Highest // xyzzy TODO

		// Generate the QR code in internal format
		var err error
		var q *goqrcode.QRCode
		q, err = goqrcode.New(pp.Url, redundancy)
		goqrcode.CheckError(err)

		// Swap colors
		if pp.Invert {
			q.ForegroundColor, q.BackgroundColor = q.BackgroundColor, q.ForegroundColor
		}

		// Output QR Code as a PNG
		var png []byte
		png, err = q.PNG(pp.Size)
		goqrcode.CheckError(err)

		www.Write(png)

	}

}

var db1 = true
