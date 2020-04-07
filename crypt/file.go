package crypt

import (
	"io/ioutil"
	"os"
)

//sDec, _ := b64.StdEncoding.DecodeString(sEnc)
//  fmt.Println(string(sDec))
// fmt.Println()

func EncryptFile(filename string, data []byte, passphrase string) {
	f, _ := os.Create(filename)
	defer f.Close()
	f.Write(EncryptWithPassphrase(data, passphrase))
}

func DecryptFile(filename string, passphrase string) []byte {
	data, _ := ioutil.ReadFile(filename)
	return DecryptWithPasphrase(data, passphrase)
}

// func EncryptWithPassphrase(data []byte, passphrase string) []byte {
// func DecryptWithPasphrase(data []byte, passphrase string) []byte {
