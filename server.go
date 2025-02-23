package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"testProject/encryption"
)

var argon encryption.Argon2id = encryption.Argon2id{
	Memory:      128 * 1024,
	Iterations:  6,
	Parallelism: 4,
	SaltLength:  32,
	KeyLength:   64,
}

type RegUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSMaxAge:         86400,
	}))
	e.POST("/login", login)
	e.POST("/register", register)
	e.POST("/refresh", authMiddleware(refreshToken))
	e.Logger.Fatal(e.Start(":8080"))
}
func refreshToken(c echo.Context) error {
	if !c.IsTLS() {
		return c.String(http.StatusUnauthorized, "Connection not secured.")
	}
	var user RegUser
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	println(user.Username)
	println(user.Password)
	println(user.Email)
	return nil
}
func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		return next(c)
	}
}
func login(c echo.Context) error {
	if !c.IsTLS() {
		return c.String(http.StatusUnauthorized, "Connection not secured.")
	}
	var user RegUser
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	println(user.Username)
	println(user.Password)
	println(user.Email)
	return nil
}
func register(c echo.Context) error {
	if !c.IsTLS() {
		return c.String(http.StatusUnauthorized, "Connection not secured.")
	}
	var user RegUser
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	println(user.Username)
	println(user.Password)
	println(user.Email)
	return nil
}
