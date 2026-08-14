package crypt

import (
	"fmt"
	"os"
)

// EncryptFile encrypts data with passphrase and writes the result to filename.
func EncryptFile(filename string, data []byte, passphrase string) error {
	ciphertext, err := EncryptWithPassphrase(data, passphrase)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, ciphertext, 0600)
}

// DecryptFile reads filename and decrypts its contents with passphrase.
func DecryptFile(filename string, passphrase string) ([]byte, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", filename, err)
	}
	return DecryptWithPassphrase(data, passphrase)
}
