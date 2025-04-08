package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/v2/bson"
	"log"
	"net/http"
	"testProject/internal/api/handlers"
	imiddleware "testProject/internal/api/middleware"
	"testProject/internal/core"
	database "testProject/internal/data"
)

func main() {
	err := core.LoadKeys()
	if err != nil {
		panic(err)
	}

	client, err := database.Connect()
	if err != nil {
		panic(err)
	}

	db := database.DB{Db: client.Database("swimply")}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderAuthorization, echo.HeaderContentType, echo.HeaderAccept},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost},
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:      "1; mode=block",
		XFrameOptions:      "SAMEORIGIN",
		ContentTypeNosniff: "nosniff",
		HSTSMaxAge:         3600,
	}))
	e.Use(middleware.BodyLimit("1M"))

	h := handlers.NewHandler(&db)
	e.POST("/signup", imiddleware.UserAuth(h.SignUp))
	e.POST("/sign-in", imiddleware.UserAuth(h.SignIn))
	e.POST("/refresh-token", imiddleware.JWTRefreshAuth(h.RefreshToken))
	e.POST("/api/v2/workout", imiddleware.JWTAccessAuth(h.NewWorkout))
	e.GET("/api/v2/workout", imiddleware.JWTAccessAuth(h.GetWorkout))
	e.DELETE("/api/v2/workout", imiddleware.JWTAccessAuth(h.DeleteWorkout))
	e.PUT("/api/v2/workout", imiddleware.JWTAccessAuth(h.ChangeWorkout))
	e.GET("/api/v2/account/info", imiddleware.JWTAccessAuth(h.GetAccountInfo))

	testId, _ := bson.ObjectIDFromHex("67ed678c23c645e03efe21f0")

	wrkts, _ := h.DB.GetAllWorkouts(testId)
	log.Println(wrkts)
	e.Logger.Fatal(e.Start(":8080"))
}
