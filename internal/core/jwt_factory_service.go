package core

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"time"
)

type JWTTokens struct{}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func (j *JWTTokens) NewToken(id bson.ObjectID, iss string, access bool) (string, error) {
	rawToken := jwt.NewWithClaims(jwt.SigningMethodEdDSA, jwt.MapClaims{
		"sub": id,
		"iss": iss,
		"exp": If(access, time.Now().Add(time.Hour*2).Unix(), time.Now().Add(time.Hour*24*7).Unix()),
		"iat": time.Now().Unix(),
		"typ": If(access, "access_token", "refresh_token"),
	})

	return rawToken.SignedString(If(access, Ed25519Keys.AccessPrivateKey, Ed25519Keys.RefreshPrivateKey))
}
func (j *JWTTokens) ParseToken(token string, access bool) (*jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
			return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
		}
		return If(access, Ed25519Keys.AccessPublicKey, Ed25519Keys.RefreshPublicKey), nil
	})
	if err != nil {
		return nil, err
	}
	var (
		claims jwt.MapClaims
		ok     bool
	)
	if claims, ok = parsedToken.Claims.(jwt.MapClaims); !ok || !parsedToken.Valid {
		return nil, &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
	}
	return &claims, nil
}
func (j *JWTTokens) VerifyClaims(claims *jwt.MapClaims, access bool) error {
	if iss, err := claims.GetIssuer(); err != nil || If(access, iss != "https://auth.swimply.pl/refresh-token" && iss != "https://auth.swimply.pl/signin", iss != "https://auth.swimply.pl/signin") {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
	}
	if iat, err := claims.GetIssuedAt(); err != nil || (iat.Unix() >= time.Now().Unix()) {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
	}
	if exp, err := claims.GetExpirationTime(); err != nil || exp.Unix() < time.Now().Unix() {
		return &echo.HTTPError{Code: http.StatusPreconditionRequired, Message: "invalid token"}
	}
	if _, err := claims.GetSubject(); err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid token"}
	}
	return nil
}

var JWTFactory = &JWTTokens{}
