package services

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/joho/godotenv"
	"io"
	"os"
	"testProject/pkg/utils"
)

func DecryptJWT(token string, refresh bool) (string, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	key, err := hex.DecodeString(utils.If(refresh, os.Getenv("JWT_REFRESH_SECRET"), os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.RawStdEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	decryptedBytes, err := gcm.Open(nil, encryptedBytes[:gcm.NonceSize()], encryptedBytes[gcm.NonceSize():], nil)
	if err != nil {
		return "", err
	}
	return string(decryptedBytes), nil
}
func EncryptJWT(token string, refresh bool) (string, error) {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	key, err := hex.DecodeString(utils.If(refresh, os.Getenv("JWT_REFRESH_SECRET"), os.Getenv("JWT_ACCESS_SECRET")))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	encryptedBytes := gcm.Seal(nonce, nonce, []byte(token), nil)

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
