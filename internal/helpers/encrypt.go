package helpers

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"working-day-api/config"
)

func Encrypt(plaintext string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(config.AppConfig.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("invalid encryption key format: %w", err)
	}

	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: expected 32 bytes, got %d", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(plaintext), nil)
	encrypted := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func Decrypt(encrypted string) (string, error) {

	key, err := base64.StdEncoding.DecodeString(config.AppConfig.EncryptionKey)
	if err != nil {
		return "", fmt.Errorf("invalid encryption key format: %v", err)
	}

	if len(key) != 32 {
		return "", fmt.Errorf("invalid key size: expected 32 bytes, got %d bytes", len(key))
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)

	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := aesGCM.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
