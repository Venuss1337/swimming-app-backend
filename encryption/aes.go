package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"github.com/joho/godotenv"
	"io"
	"os"
)

func DecryptAES(plain string) (string, error) {
	if err := godotenv.Load("secret.env"); err != nil {
		panic(err)
	}

	key, err := hex.DecodeString(os.Getenv("JWT_SEC_NONCE"))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	encryptedBytes, err := base64.RawStdEncoding.DecodeString(plain)
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
func EncryptAES(plain string) (string, error) {
	if err := godotenv.Load("secret.env"); err != nil {
		panic(err)
	}

	key, err := hex.DecodeString(os.Getenv("JWT_SEC_NONCE"))
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
	encryptedBytes := gcm.Seal(nonce, nonce, []byte(plain), nil)

	return base64.StdEncoding.EncodeToString(encryptedBytes), nil
}
