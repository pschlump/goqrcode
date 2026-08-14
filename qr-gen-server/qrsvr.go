package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	svg "github.com/ajstarks/svgo"
	"github.com/boombuler/barcode/qr"
	"github.com/pschlump/ReadConfig"
	"github.com/pschlump/goqrcode"
	params "github.com/pschlump/goqrcode/qr-gen-server/params"
	"github.com/pschlump/goqrsvg"

	"github.com/pschlump/dbgo"
)

var Cfg = flag.String("cfg", "cfg.json", "config file for this call")
var HostPort = flag.String("hostport", ":2003", "Host/Port to listen on")
var DbFlagParam = flag.String("db_flag", "", "Additional Debug Flags")
var Version = flag.Bool("version", false, "Report version of code and exit")
var Comment = flag.String("comment", "", "Unused comment for ps.")
var CdTo = flag.String("CdTo", ".", "Change directory to before running server.")

// GlobalConfigData holds the server configuration read from the --cfg file.
type GlobalConfigData struct {
}

var gCfg GlobalConfigData
var DbOn map[string]bool

func init() {
	DbOn = make(map[string]bool)
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
			fmt.Fprintf(os.Stderr, "Unable to change to %s directory, error:%s\n", *CdTo, err)
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
		fmt.Fprintf(os.Stderr, "Unable to read configuration: %s error %s\n", *Cfg, err)
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
	// fmt PNG, SVG, text
	pp, err := params.ParseParams(www, req, "GET", "POST")
	log.Printf("%s: %s, error:%s, at:%s\n", req.Method, dbgo.SVarI(pp), err, dbgo.LF())
	if err != nil {
		return
	}
	genqr_bl(www, req, pp)
}

func genqr_bl(www http.ResponseWriter, req *http.Request, pp params.ApiDataType) {

	switch pp.Fmt {
	case "svg":

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
		qrCode, err := qr.Encode(pp.Url, red, qr.Auto)
		if err != nil {
			log.Printf("Error: %s\n", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Error: %s\n", err)
			return
		}

		if DbOn["echo-params-1"] {
			fmt.Printf("at: %s Params: %+v\n", dbgo.LF(), pp)
		}

		// Write QR code to SVG
		qs := goqrsvg.NewQrSVG(qrCode, 5)
		qs.StartQrSVG(s)
		if err := qs.WriteQrSVG(s); err != nil {
			log.Printf("Error: %s\n", err)
			return
		}

		s.End()

	case "png":

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
		q, err := goqrcode.New(pp.Url, redundancy)
		if err != nil {
			log.Printf("Error: %s\n", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Error: %s\n", err)
			return
		}

		// Swap colors
		if pp.Invert {
			q.ForegroundColor, q.BackgroundColor = q.BackgroundColor, q.ForegroundColor
		}

		// Output QR Code as a PNG
		png, err := q.PNG(pp.Size)
		if err != nil {
			log.Printf("Error: %s\n", err)
			www.WriteHeader(http.StatusInternalServerError) // 500
			fmt.Fprintf(www, "Error: %s\n", err)
			return
		}

		if _, err := www.Write(png); err != nil {
			log.Printf("Error writing response: %s\n", err)
		}

	case "text":

		www.Header().Set("Content-Type", "text/plain; charset=utf-8")

		redundancy := goqrcode.Highest
		switch pp.Redundancy {
		case "H":
			redundancy = goqrcode.Highest
		case "M":
			redundancy = goqrcode.Medium
		case "L":
			redundancy = goqrcode.Low
		}

		q, err := goqrcode.New(pp.Url, redundancy)
		if err != nil {
			log.Printf("Error: %s\n", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Error: %s\n", err)
			return
		}

		if _, err := fmt.Fprint(www, q.ToString(pp.Invert)); err != nil {
			log.Printf("Error writing response: %s\n", err)
		}

	default:
		log.Printf("Error: unsupported format %q\n", pp.Fmt)
		www.WriteHeader(http.StatusBadRequest) // 400
		fmt.Fprintf(www, "Error: unsupported format %q, expected png, svg or text\n", pp.Fmt)
	}
}
