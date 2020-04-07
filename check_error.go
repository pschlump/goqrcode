package goqrcode

import (
	"fmt"
	"os"

	"github.com/pschlump/godebug"
)

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s from: %s\n", err, godebug.LF(2))
		os.Exit(1)
	}
}
