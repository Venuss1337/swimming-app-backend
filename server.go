package main

import (
	"crypto/tls"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"net/http"
	"testProject/database"
	"testProject/encryption"
	"testProject/model"
	"time"
)

var argon encryption.Argon2id = encryption.Argon2id{
	Memory:      128 * 1024,
	Iterations:  6,
	Parallelism: 4,
	SaltLength:  32,
	KeyLength:   64,
}

type ResponseLogin struct {
	RefreshToken string
	AccessToken  string
}
type ResponeRefresh struct {
	AccessToken string
}

var db database.Database

func main() {
	e := echo.New()
	e.Pre(middleware.HTTPSNonWWWRedirect())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         86400,
	}))
	/*
			{
				"username":"veanut",
				"password":"pass123",
				"email":"example@gmail.com"
			}

		  0l3CsTZ2pfdVOkuQdA2HhvHWfICj05ZGRTECjyHAPTHCoVPQEioR+qSQVz3JcMUfd5VKbDaH4Hzd8DKQuFsd/vidfoJikMn7HdqI87nNA5bYEnzScGhqTN4d01JcpTQ+yTfFzOsGqzenuEPVFLm/J0dj7gpYXad9+wR+YsATe/gMY1HEqz1gl1BhvKvz86r7ROBNXpJL0mxfdeU6X8oMBLdZYxahX2UDbSPAaEx4KrJkEmjHaOqrgUgZJpzcVNkawUZcf/Ybb0LyuOu25PoY/ibZwAavZbz/xj2jtFMiJgVGk491oDxlgmA0jHY//uSI0ZbNKhaxHaensSEECUq1lYwOmHyS0XDD311zS74UHrVFN6IUsNfEaekdNRaE+O3S4lJDCiWrggppIemcHyOHvFwKtYoTUq4CrDHJ6ucihOIGQrdmPwEhsgXq

	*/
	e.POST("/login", login)
	e.POST("/register", register)
	e.POST("/refresh", authMiddleware(refreshToken))

	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
	}
	server := &http.Server{
		Addr:      ":8080",
		Handler:   e,
		TLSConfig: tlsConfig,
	}

	db = database.Database{}
	err := db.Connect("mongodb://localhost:27017")
	if err != nil {
		return
	}
	println("Successfully connected to database")

	println("Server started listening on :8080")
	e.Logger.Fatal(server.ListenAndServeTLS("cert.pem", "key.pem"))
}

func refreshToken(c echo.Context) error {
	claims := c.Get("claims").(jwt.MapClaims)

	if tokenType, ok := claims["typ"].(string); !ok || tokenType != "refresh" {
		return c.JSON(http.StatusUnauthorized, "Invalid token7")
	}
	if nbf, ok := claims["nbf"].(float64); !ok || time.Now().Before(time.Unix(int64(nbf), 0)) {
		return c.JSON(http.StatusUnauthorized, "Invalid token8")
	}
	if exp, ok := claims["exp"].(float64); !ok || time.Now().After(time.Unix(int64(exp), 0)) {
		return c.JSON(http.StatusUnauthorized, "Token expired")
	}
	if iss, ok := claims["iss"].(string); !ok || (iss != "swimply.pl/api/v2/register" && iss != "swimply.pl/api/v2/login") {
		return c.JSON(http.StatusUnauthorized, "Invalid token9")
	}
	if sub, ok := claims["sub"].(string); !ok || sub == "" {
		return c.JSON(http.StatusUnauthorized, "Invalid token10")
	}
	sub, err := bson.ObjectIDFromHex(claims["sub"].(string))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, "Invalid token11")
	}

	accessToken, err := encryption.CreateToken(sub, map[string]interface{}{
		"iss": "swimply.pl/api/v2/refresh-token",
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "access",
	})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	c.Response().Header().Set("Content-Type", "application/json")
	c.Response().WriteHeader(http.StatusOK)
	b, err := json.Marshal(&ResponeRefresh{AccessToken: accessToken})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if _, err := c.Response().Write(b); err != nil {
		return err
	}
	return nil
}
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !c.IsTLS() {
			return echo.NewHTTPError(http.StatusBadRequest, "Connection not secured")
		}
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "Missing token")
		}
		tokenString, err := encryption.ParseJWT(authHeader)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token1")
		}
		if _, ok := tokenString.Claims.(jwt.MapClaims); !ok || !tokenString.Valid {
			return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token2")
		}
		c.Set("claims", tokenString.Claims)
		return next(c)
	}
}
func login(c echo.Context) error {
	if !c.IsTLS() {
		return c.String(http.StatusBadRequest, "Connection not secured")
	}
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "Invalid JSON")
	}
	retrievedHash, userId, err := db.RetrievePasswordHashAndId(user.Username)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if match, err := argon.Verify([]byte(user.Password), retrievedHash); err != nil {
		return c.String(http.StatusBadRequest, "Something went wrong")
	} else if !match {
		return c.String(http.StatusBadRequest, "Invalid password")
	}

	refreshToken, err := encryption.CreateToken(userId, map[string]interface{}{
		"iss": "swimply.pl/api/v2/login",
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "refresh",
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	accessToken, err := encryption.CreateToken(userId, map[string]interface{}{
		"iss": "swimply.pl/api/v2/login",
		"exp": time.Now().Add(time.Minute * 30).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"typ": "access",
	})
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	responseJson := ResponseLogin{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	responeBytes, err := json.Marshal(responseJson)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	c.Response().WriteHeader(http.StatusOK)
	_, err = c.Response().Write(responeBytes)
	if err != nil {
		return err
	}
	return nil
}

func register(c echo.Context) error {
	if !c.IsTLS() {
		return c.String(http.StatusBadRequest, "Connection not secured.")
	}
	var user model.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var passwordHash string
	if err := argon.Hash(&passwordHash, []byte(user.Password)); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	if err := db.RegisterUser(user.Username, passwordHash); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	// Disabled until I figure it out how to do this part to not call database 92834190382 times :)
	/*refreshToken, err := encryption.CreateRefreshToken(userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error()+" tutaj")
	}
	accessToken, err := encryption.CreateAccessToken(userId)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	responseJson := ResponseJson{
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	responeBytes, err := json.Marshal(responseJson)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	c.Response().Header().Set(echo.HeaderContentType, "application/json")
	c.Response().WriteHeader(http.StatusOK)
	_, err = c.Response().Write(responeBytes)
	if err != nil {
		return err
	}*/
	return c.String(http.StatusOK, "Successfully registered user")
}
