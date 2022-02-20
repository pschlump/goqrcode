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

// params "github.com/pschlump/goqrcode/qr-gen-server/params"

func main() {
	http.Handle("/", http.HandlerFunc(genqr))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func genqr(www http.ResponseWriter, req *http.Request) {

	var method bool
	var err error
	var pp params.ApiTestDataType
	// size, url
	// redundancy H, M, L
	// fmt PNG, SVG, SVG-fragment

	if req.Method == "GET" {
		pp, method, err = params.ParseGETParams(www, req) //  (rv ApiTestDataType, methodGet bool, err error) {
		_, _ = method, err
		fmt.Printf("GET: %s\n", godebug.SVarI(pp))
	} else if req.Method == "POST" {
		pp, method, err = params.ParsePOSTParams(www, req) //  (rv ApiTestDataType, methodGet bool, err error) {
		_, _ = method, err
		fmt.Printf("POST: %s\n", godebug.SVarI(pp))
	} else {
		www.WriteHeader(http.StatusMethodNotAllowed) // 405
		fmt.Fprintf(www, "Invalid Method")
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
