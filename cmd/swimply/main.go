package main

import (
	"crypto/tls"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"testProject/internal/config"
	"testProject/internal/handlers"
	"testProject/internal/middleware"
	"testProject/internal/repository"
	"testProject/pkg/utils"
)

func main() {

	if err := repository.Connect("mongodb://localhost:27017"); err != nil {
		panic(err)
	}
	if err := config.LoadSecrets(); err != nil {
		panic(err)
	}
	utils.InitConfig(&utils.Argon2id{
		Memory: 47104, Iterations: 1, Parallelism: 1, SaltLength: 32, KeyLength: 64,
	})

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodPost},
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		XFrameOptions:      "SAMEORIGIN",
		ContentTypeNosniff: "nosniff",
		HSTSMaxAge:         3600,
	}))
	e.Use(middleware.BodyLimit("1M"))
	e.POST("/login", imiddleware.AuthenticationMiddleware(handlers.LoginHandler))
	/*e.POST("/register", imiddleware.AuthenticationMiddleware(handlers.RegisterHandler))*/
	server := http.Server{
		Addr:    ":8080",
		Handler: e,
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	}
	e.Logger.Fatal(server.ListenAndServeTLS(config.TlsCert, config.TLSKey))

}
