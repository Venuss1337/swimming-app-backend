package encryption

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"os"
	"time"
)

func CreateRefreshToken(id bson.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": id,
		"iss": "swimply.pl/api/v2/login",
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "refresh",
	})

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

	return token.SignedString(privateKey)
}

func CreateAccessToken(id bson.ObjectID) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": id,
		"iss": "swimply.pl/api/v2/refresh",
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "access",
	})

	privateKey, err := readPrivateFromFile()
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	publicKey, err := readPublicFromFile()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
