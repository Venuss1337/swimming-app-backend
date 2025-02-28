package encryption

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func CreateToken(id bson.ObjectID, claims map[string]interface{}) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": id,
		"iss": claims["iss"],
		"exp": claims["exp"],
		"iat": claims["iat"],
		"nbf": claims["nbf"],
		"typ": claims["typ"],
	})

	privateKey, err := readPrivateFromFile()
	if err != nil {
		return "", nil
	}
	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}
	return EncryptAES(signedToken)
}

func ParseJWT(encryptedToken string) (*jwt.Token, error) {

	tokenString, err := DecryptAES(encryptedToken)
	if err != nil {
		return nil, err
	}

	publicKey, err := readPublicFromFile()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}
