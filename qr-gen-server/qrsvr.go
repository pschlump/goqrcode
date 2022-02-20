package main

import (
	"fmt"
	"log"
	"net/http"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"github.com/pschlump/godebug"
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

	s := svg.New(www)
	www.Header().Set("Content-Type", "image/svg+xml")

	// Create the QR Code in SVG
	// qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)
	qrCode, _ := qr.Encode(pp.Url, qr.M, qr.Auto)

	// Write QR code to SVG
	qs := goqrsvg.NewQrSVG(qrCode, 5)
	qs.StartQrSVG(s)
	qs.WriteQrSVG(s)

	s.End()
}
