package main

import (
	"encoding/base64"
	"fmt"
	"os"

	"github.com/pschlump/goqrcode/crypt"
)

func main() {

	data := []byte("http://www.2c-why.com/")
	passphrase := "pointer stew"
	ec, err := crypt.EncryptWithPassphrase(data, passphrase)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to encrypt: %s\n", err)
		os.Exit(1)
	}
	sEnc := base64.StdEncoding.EncodeToString(ec)
	fmt.Printf("%s\n", sEnc)
}
