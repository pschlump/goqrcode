package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"github.com/pschlump/godebug"
	"github.com/pschlump/goqrcode"
	params "github.com/pschlump/goqrcode/qr-gen-server/params"
	"github.com/pschlump/goqrsvg"

	"github.com/pschlump/ReadConfig"
	ymux "github.com/pschlump/gintools/data"
)

// ymux "git.q8s.co/pschlump/piserver/ymux"

// TODO
// var TLS_crt = flag.String("tls_crt", "", "TLS Signed Publick Key")
// var TLS_key = flag.String("tls_key", "", "TLS Signed Private Key")
// var ChkTables = flag.Bool("chk-tables", false, "Chack table structre and exit")
var Cfg = flag.String("cfg", "cfg.json", "config file for this call")
var HostPort = flag.String("hostport", ":2003", "Host/Port to listen on")
var DbFlagParam = flag.String("db_flag", "", "Additional Debug Flags")
var Version = flag.Bool("version", false, "Report version of code and exit")
var Comment = flag.String("comment", "", "Unused comment for ps.")
var CdTo = flag.String("CdTo", ".", "Change directory to before running server.")

type GlobalConfigData struct {
	ymux.BaseConfigType
}

var gCfg GlobalConfigData
var DB *sql.DB
var DbOn map[string]bool
var logger *log.Logger
var isTLS bool
var curDir string

func init() {
	isTLS = false
	DbOn = make(map[string]bool)
	logger = log.New(os.Stdout, "", 0)
}

func main() {

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "qrsvr2 : Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Printf("Extra arguments are not supported [%s]\n", fns)
		os.Exit(1)
	}

	if *CdTo != "" {
		err := os.Chdir(*CdTo)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Uable to chagne to %s directory, error:%s\n", *CdTo, err)
			os.Exit(1)
		}
	}

	if *Version {
		fmt.Printf("Version (Git Commit): %s\n", GitCommit)
		os.Exit(0)
	}

	if Cfg == nil {
		fmt.Printf("--cfg is a required parameter\n")
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Read in Configuration
	// ------------------------------------------------------------------------------
	err := ReadConfig.ReadFile(*Cfg, &gCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read confguration: %s error %s\n", *Cfg, err)
		os.Exit(1)
	}

	// ------------------------------------------------------------------------------
	// Debug Flag Processing
	// ------------------------------------------------------------------------------
	// xyzzy - fix this - put back in
	// ymux.DebugFlagProcess(DbFlagParam, DbOn, &(gCfg.BaseConfigType))

	// ------------------------------------------------------------------------------
	// ------------------------------------------------------------------------------
	http.Handle("/", http.HandlerFunc(genqr))
	err = http.ListenAndServe(*HostPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func genqr(www http.ResponseWriter, req *http.Request) {

	// size, url
	// redundancy H, M, L
	// fmt PNG, SVG, SVG-fragment
	pp, err := params.ParseParams(www, req, "GET", "POST")
	fmt.Printf("%s: %s, error:%s, at:%s\n", req.Method, godebug.SVarI(pp), err, godebug.LF())
	if err != nil {
		return
	}
	genqr_bl(www, req, pp)
}

func genqr_bl(www http.ResponseWriter, req *http.Request, pp params.ApiDataType) {

	if pp.Fmt == "svg" {

		s := svg.New(www)
		www.Header().Set("Content-Type", "image/svg+xml")

		red := qr.M
		switch pp.Redundancy {
		case "H":
			red = qr.H
		case "M":
			red = qr.M
		case "L":
			red = qr.L
		}
		// Create the QR Code in SVG
		// qrCode, _ := qr.Encode("Hello World", qr.M, qr.Auto)
		qrCode, err := qr.Encode(pp.Url, red, qr.Auto)
		if err != nil {
			log.Printf("Error: %s\n", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Error: %s\n", err)
			err = fmt.Errorf("Error: %s", err)
			return
		}

		if DbOn["echo-params-1"] {
			fmt.Printf("at: %s Params: %+v\n", godebug.LF(), pp)
		}

		// Write QR code to SVG
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(s)
		qs.WriteQrSVG(s)

		s.End()

	} else if pp.Fmt == "png" {

		www.Header().Set("Content-Type", "image/png")

		redundancy := goqrcode.Highest
		switch pp.Redundancy {
		case "H":
			redundancy = goqrcode.Highest
		case "M":
			redundancy = goqrcode.Medium
		case "L":
			redundancy = goqrcode.Low
		}

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

	} else if pp.Fmt == "text" {

		// xyzzy - TODO

	}
}
