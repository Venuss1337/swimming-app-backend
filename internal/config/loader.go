package config

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
)

const (
	secretDirs        = "secrets"
	tlsCertPath       = "cert.pem"
	tlsKeyPath        = "key.pem"
	jwtPrivateKeyPath = "private.pem"
	jwtPublicKeyPath  = "public.pem"
)

var (
	TlsCert       string
	TLSKey        string
	JWTPrivateKey any
	JWTPublicKey  any
)

func LoadSecrets() error {
	certPath := filepath.Join(secretDirs, tlsCertPath)
	keyPath := filepath.Join(secretDirs, tlsKeyPath)
	jwtPrivateKeyPath := filepath.Join(secretDirs, jwtPrivateKeyPath)
	jwtPublicKeyPath := filepath.Join(secretDirs, jwtPublicKeyPath)

	TlsCert = certPath
	TLSKey = keyPath

	var err error
	JWTPrivateKey, err = ReadPrivateKey(jwtPrivateKeyPath)
	if err != nil {
		return err
	}
	JWTPublicKey, err = ReadPublicKey(jwtPublicKeyPath)
	if err != nil {
		return err
	}
	return nil
}
func ReadPrivateKey(keyPath string) (any, error) {
	b, err := os.ReadFile(keyPath)
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

func ReadPublicKey(keyPath string) (any, error) {
	b, err := os.ReadFile(keyPath)
	if err != nil {
		return "", err
	}

	block, _ := pem.Decode(b)
	if block == nil || block.Type != "PUBLIC KEY" {
		return "", errors.New("failed to decode PEM block containing public key")
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)

	if err != nil {
		return "", err
	}
	return publicKey, nil
}
