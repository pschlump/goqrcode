package ApiQRGen

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/pschlump/req_param/reqlib"
)

type ApiTestDataType struct {
	Url    string
	Fmt    string
	Invert bool
	Size   int
}

func ParsePOSTParams(www http.ResponseWriter, req *http.Request) (rv ApiTestDataType, methodPost bool, err error) {

	if req.Method == "POST" {
		methodPost = true
		var s string
		var b bool
		var n64 int64

		req.ParseForm()

		// --------------------------------------------------------------------
		// Parameter: url
		// Method: POST
		// --------------------------------------------------------------------
		s = req.Form.Get("url")

		if "" != "" {
			s = ""
		} else if false {
			log.Println("Url Param 'url' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'url' paramter\n")
			err = fmt.Errorf("Missing 'url' paramter")
			return
		} else {
			s = ""
		}

		rv.Url = s

		// --------------------------------------------------------------------
		// Parameter: fmt
		// Method: POST
		// --------------------------------------------------------------------
		s = req.Form.Get("fmt")

		if "png" != "" {
			s = "png"
		} else if false {
			log.Println("Url Param 'fmt' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'fmt' paramter\n")
			err = fmt.Errorf("Missing 'fmt' paramter")
			return
		} else {
			s = ""
		}

		rv.Fmt = s

		// --------------------------------------------------------------------
		// Parameter: invert
		// Method: POST
		// --------------------------------------------------------------------
		s = req.Form.Get("invert")

		if "false" != "" {
			s = "false"
		} else if false {
			log.Println("Url Param 'invert' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'invert' paramter\n")
			err = fmt.Errorf("Missing 'invert' paramter")
			return
		} else {
			s = ""
		}

		b, err = reqlib.ParseBool(s)
		if err != nil {
			log.Println("Url Param 'invert' is invalid boolean: %s", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Url Param 'invert' is invalid boolean, %s\n", err)
			err = fmt.Errorf("Url Param 'invert' is invalid boolean, %s", err)
			return
		}

		rv.Invert = b

		// --------------------------------------------------------------------
		// Parameter: size
		// Method: POST
		// --------------------------------------------------------------------
		s = req.Form.Get("size")

		if "256" != "" {
			s = "256"
		} else if false {
			log.Println("Url Param 'size' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'size' paramter\n")
			err = fmt.Errorf("Missing 'size' paramter")
			return
		} else {
			s = ""
		}

		n64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Println("Url Param 'size' is invalid integer: %s", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Url Param 'size' is invalid integer, %s\n", err)
			err = fmt.Errorf("Url Param 'size' is invalid integer, %s", err)
			return
		}
		rv.Size = int(n64)

	}
	return
}

func ParseGETParams(www http.ResponseWriter, req *http.Request) (rv ApiTestDataType, methodGet bool, err error) {

	if req.Method == "GET" {
		methodGet = true
		params := req.URL.Query()
		var sA []string
		var s string
		var ok, b bool
		var n64 int64

		// --------------------------------------------------------------------
		// Parameter: invert
		// Method: GET
		// --------------------------------------------------------------------
		sA, ok = params["invert"]
		if ok && len(sA) >= 1 {
			s = sA[0]
		} else if "false" != "" {
			s = "false"
		} else if false {
			log.Println("Url Param 'invert' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'invert' paramter\n")
			err = fmt.Errorf("Missing 'invert' paramter")
			return
		} else {
			s = ""
		}

		b, err = reqlib.ParseBool(s)
		if err != nil {
			log.Println("Url Param 'invert' is invalid boolean: %s", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Url Param 'invert' is invalid boolean, %s\n", err)
			err = fmt.Errorf("Url Param 'invert' is invalid boolean, %s", err)
			return
		}

		rv.Invert = b

		// --------------------------------------------------------------------
		// Parameter: size
		// Method: GET
		// --------------------------------------------------------------------
		sA, ok = params["size"]
		if ok && len(sA) >= 1 {
			s = sA[0]
		} else if "256" != "" {
			s = "256"
		} else if false {
			log.Println("Url Param 'size' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'size' paramter\n")
			err = fmt.Errorf("Missing 'size' paramter")
			return
		} else {
			s = ""
		}

		n64, err = strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Println("Url Param 'size' is invalid integer: %s", err)
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Url Param 'size' is invalid integer, %s\n", err)
			err = fmt.Errorf("Url Param 'size' is invalid integer, %s", err)
			return
		}
		rv.Size = int(n64)

		// --------------------------------------------------------------------
		// Parameter: url
		// Method: GET
		// --------------------------------------------------------------------
		sA, ok = params["url"]
		if ok && len(sA) >= 1 {
			s = sA[0]
		} else if "" != "" {
			s = ""
		} else if false {
			log.Println("Url Param 'url' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'url' paramter\n")
			err = fmt.Errorf("Missing 'url' paramter")
			return
		} else {
			s = ""
		}

		rv.Url = s

		// --------------------------------------------------------------------
		// Parameter: fmt
		// Method: GET
		// --------------------------------------------------------------------
		sA, ok = params["fmt"]
		if ok && len(sA) >= 1 {
			s = sA[0]
		} else if "png" != "" {
			s = "png"
		} else if false {
			log.Println("Url Param 'fmt' is missing")
			www.WriteHeader(http.StatusBadRequest) // 400
			fmt.Fprintf(www, "Missing 'fmt' paramter\n")
			err = fmt.Errorf("Missing 'fmt' paramter")
			return
		} else {
			s = ""
		}

		rv.Fmt = s

	}
	return
}
