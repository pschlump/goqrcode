package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
)

// CreateHash returns the hex encoded MD5 sum of key. The 32 byte digest is
// used as an AES-256 key by the encrypt/decrypt functions below.
func CreateHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// EncryptWithPassphrase encrypts data with AES-256-GCM using a key derived
// from passphrase. The returned bytes are nonce || ciphertext.
func EncryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	block, err := aes.NewCipher([]byte(CreateHash(passphrase)))
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("reading nonce: %w", err)
	}
	return gcm.Seal(nonce, nonce, data, nil), nil
}

// DecryptWithPassphrase decrypts data produced by EncryptWithPassphrase.
func DecryptWithPassphrase(data []byte, passphrase string) ([]byte, error) {
	key := []byte(CreateHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short, got %d bytes, need at least %d", len(data), nonceSize)
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("gcm.Open: %w", err)
	}
	return plaintext, nil
}

// DecryptWithPasphrase is a misspelled deprecated alias for
// DecryptWithPassphrase.
//
// Deprecated: use DecryptWithPassphrase.
func DecryptWithPasphrase(data []byte, passphrase string) ([]byte, error) {
	return DecryptWithPassphrase(data, passphrase)
}
