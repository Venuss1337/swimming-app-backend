package services

import (
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testProject/internal/config"
)

func GenerateJWT(id bson.ObjectID, claims map[string]interface{}) (string, error) {
	rawToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"iss": claims["iss"],
		"sub": id,
		"exp": claims["sub"],
		"iat": claims["sub"],
		"typ": claims["typ"],
	})

	signedToken, err := rawToken.SignedString(config.JWTPrivateKey)
	if err != nil {
		return "", err
	}

	encryptedToken, err := EncryptJWT(signedToken, claims["typ"] == "refresh_token")
	if err != nil {
		return "", err
	}

	return encryptedToken, nil
}
