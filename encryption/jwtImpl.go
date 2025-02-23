package encryption

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
	"os"
	"time"
)

func CreateRefreshToken(id bson.ObjectID) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": id,
		"iss": "swaip.com/api/v2/refresh",
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "access",
	})

	privateKeyBytes, err := os.ReadFile("private.pem")
	if err != nil {
		return "", err
	}

	privateKey, err := jwt.ParseEdPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		return "", err
	}

	return token.SignedString(privateKey)
}
func createAccessToken() (string, error) {

	return "", nil
}

func ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return token, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
