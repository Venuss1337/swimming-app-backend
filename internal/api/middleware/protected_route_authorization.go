package imiddleware

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"os"
	"strings"
	"testProject/internal/core"
	"testProject/internal/models"
)

func JWTRefreshAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			after = ""
			found = false
		)
		if after, found = strings.CutPrefix(c.Request().Header.Get("Authorization"), "Bearer "); after == "" || !found {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "invalid token1"}
		}

		decodedToken, err := base64.StdEncoding.DecodeString(after)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "invalid token2"}
		}

		key, err := hex.DecodeString(os.Getenv("JWT_REFRESH_SECRET"))
		if err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "encryption failed"}
		}

		decryptedToken, err := core.JWTEncrypter.Decrypt(decodedToken, key)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "encryption failed"}
		}
		claims, err := core.JWTFactory.ParseToken(string(decryptedToken), false)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: err.Error()}
		}

		if err := core.JWTFactory.VerifyClaims(claims, false); err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}

		userId, _ := claims.GetSubject()
		rawUserId, err := bson.ObjectIDFromHex(userId)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: err.Error()}
		}

		c.Set("user", &models.User{Id: rawUserId})
		return next(c)
	}
}
func JWTAccessAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var (
			after = ""
			found = false
		)
		if after, found = strings.CutPrefix(c.Request().Header.Get("Authorization"), "Bearer "); after == "" || !found {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "invalid token"}
		}

		decodedToken, err := base64.StdEncoding.DecodeString(after)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "invalid token"}
		}

		key, err := hex.DecodeString(os.Getenv("JWT_ACCESS_SECRET"))
		if err != nil {
			return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "encryption failed"}
		}

		decryptedToken, err := core.JWTEncrypter.Decrypt(decodedToken, key)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: "encryption failed"}
		}
		claims, err := core.JWTFactory.ParseToken(string(decryptedToken), true)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: err.Error()}
		}

		if err := core.JWTFactory.VerifyClaims(claims, true); err != nil {
			return &echo.HTTPError{Code: http.StatusBadRequest, Message: err.Error()}
		}

		userId, _ := claims.GetSubject()
		rawUserId, err := bson.ObjectIDFromHex(userId)
		if err != nil {
			return &echo.HTTPError{Code: http.StatusForbidden, Message: err.Error()}
		}

		c.Set("user", &models.User{Id: rawUserId})
		return next(c)
	}
}