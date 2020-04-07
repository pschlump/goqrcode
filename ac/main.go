package main

import (
	"encoding/base64"
	"fmt"

	"gitlab.com/pschlump/goqrcode/crypt"
)

func main() {

	data := []byte("http://www.2c-why.com/")
	passphrase := "pointer stew"
	ec := crypt.EncryptWithPassphrase(data, passphrase)
	sEnc := base64.StdEncoding.EncodeToString(ec)
	fmt.Printf("%s\n", sEnc)
}
