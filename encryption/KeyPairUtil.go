package encryption

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func GenerateNewKeyPair() error {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return err
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	if err := savePrivateKeyToFile(privateKeyBytes); err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	if err := savePublicKeyToFile(publicKeyBytes); err != nil {
		return err
	}
	return nil
}

func savePrivateKeyToFile(privateKeyBytes []byte) error {
	file, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
}
func readPrivateFromFile() (any, error) {
	b, err := os.ReadFile("private.pem")
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(b)
	if block == nil || block.Type != "PRIVATE KEY" {
		return "", errors.New("failed to decode PEM block containing private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)

	if err != nil {
		return "", err
	}
	return privateKey, nil
}
func savePublicKeyToFile(publicKeyBytes []byte) error {
	file, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer file.Close()

	return pem.Encode(file, &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
}
func readPublicFromFile() (any, error) {
	b, err := os.ReadFile("public.pem")
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(b)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	privateKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return "", err
	}
	return privateKey, nil
}
