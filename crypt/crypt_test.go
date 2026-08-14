package crypt

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestEncryptDecryptRoundTrip(t *testing.T) {
	tests := []struct {
		name       string
		data       []byte
		passphrase string
	}{
		{"url", []byte("http://www.2c-why.com/"), "pointer stew"},
		{"empty data", []byte{}, "secret"},
		{"binary", []byte{0x00, 0x01, 0xfe, 0xff}, "p"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ciphertext, err := EncryptWithPassphrase(test.data, test.passphrase)
			if err != nil {
				t.Fatalf("EncryptWithPassphrase: %s", err)
			}

			if bytes.Equal(ciphertext, test.data) && len(test.data) > 0 {
				t.Fatal("ciphertext identical to plaintext")
			}

			plaintext, err := DecryptWithPassphrase(ciphertext, test.passphrase)
			if err != nil {
				t.Fatalf("DecryptWithPassphrase: %s", err)
			}

			if !bytes.Equal(plaintext, test.data) {
				t.Fatalf("got %q, expected %q", plaintext, test.data)
			}
		})
	}
}

func TestDecryptWrongPassphrase(t *testing.T) {
	ciphertext, err := EncryptWithPassphrase([]byte("data"), "correct horse battery staple")
	if err != nil {
		t.Fatalf("EncryptWithPassphrase: %s", err)
	}

	if _, err := DecryptWithPassphrase(ciphertext, "wrong passphrase"); err == nil {
		t.Fatal("decryption with wrong passphrase succeeded, expected failure")
	}
}

func TestDecryptShortData(t *testing.T) {
	if _, err := DecryptWithPassphrase([]byte{0x01}, "passphrase"); err == nil {
		t.Fatal("decryption of too-short data succeeded, expected failure")
	}
}

func TestDecryptWithPasphraseDeprecatedAlias(t *testing.T) {
	ciphertext, err := EncryptWithPassphrase([]byte("data"), "passphrase")
	if err != nil {
		t.Fatalf("EncryptWithPassphrase: %s", err)
	}

	plaintext, err := DecryptWithPasphrase(ciphertext, "passphrase")
	if err != nil {
		t.Fatalf("DecryptWithPasphrase: %s", err)
	}

	if string(plaintext) != "data" {
		t.Fatalf("got %q, expected %q", plaintext, "data")
	}
}

func TestEncryptDecryptFile(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "secret.dat")
	data := []byte("http://www.2c-why.com/")

	if err := EncryptFile(filename, data, "pointer stew"); err != nil {
		t.Fatalf("EncryptFile: %s", err)
	}

	onDisk, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading %s: %s", filename, err)
	}

	if bytes.Equal(onDisk, data) {
		t.Fatal("file on disk contains plaintext")
	}

	plaintext, err := DecryptFile(filename, "pointer stew")
	if err != nil {
		t.Fatalf("DecryptFile: %s", err)
	}

	if !bytes.Equal(plaintext, data) {
		t.Fatalf("got %q, expected %q", plaintext, data)
	}
}

func TestCreateHash(t *testing.T) {
	// MD5 of the empty string.
	const emptyMD5 = "d41d8cd98f00b204e9800998ecf8427e"
	if got := CreateHash(""); got != emptyMD5 {
		t.Fatalf("CreateHash(\"\") = %s, expected %s", got, emptyMD5)
	}

	// The hex digest is 32 characters (128 bits) and usable as an AES-256 key
	// once decoded.
	h := CreateHash("passphrase")
	if len(h) != 32 {
		t.Fatalf("CreateHash returned %d characters, expected 32", len(h))
	}
}
