package handlers
import (

"encoding/base64"
"encoding/hex"
"github.com/joho/godotenv"
"github.com/labstack/echo/v4"
"net/http"
"os"
"testProject/internal/core"
"testProject/internal/models"
)

func (h *Handler) SignUp(c echo.Context) error {
	user := c.Get("user").(models.User)

	if userExists, err := h.DB.Exists(user.Username); err != nil || userExists {
		return &echo.HTTPError{Code: http.StatusBadRequest, Message: "user already exists"}
	}

	hash, err := core.ArgonHashService.Hash([]byte(user.Password))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "hashing failed"}
	}

	if err := h.DB.NewUser(user.Id, user.Username, hash, user.Weight, user.IsMale); err != nil {
		return &echo.HTTPError{Code: http.StatusInternalServerError, Message: "user not created"}
	}
	user.Password = ""
	user.Weight = 0
	user.IsMale = false

	return c.JSON(http.StatusCreated, user)
}
func (h *Handler) SignIn(c echo.Context) error {
	user := c.Get("user").(models.User)

	if userExists, err := h.DB.Exists(user.Username); err != nil || !userExists {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "invalid user or password"}
	}

	dbUser, err := h.DB.GetUserByName(user.Username)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "invalid user or password"}
	}

	if core.ArgonHashService.Verify([]byte(user.Password), dbUser.Password) != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "invalid user or password"}
	}
	user.Password = ""
	user.Weight = 0
	user.IsMale = false


	rawAccessToken, err := core.JWTFactory.NewToken(user.Id, "https://auth.swimply.pl/signin", true)
	rawRefreshToken, err := core.JWTFactory.NewToken(user.Id, "https://auth.swimply.pl/signin", false)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error signing token"}
	}
	godotenv.Load()

	encAccessKey, err := hex.DecodeString(os.Getenv("JWT_ACCESS_SECRET"))
	encRefreshKey, err := hex.DecodeString(os.Getenv("JWT_REFRESH_SECRET"))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error creating token"}
	}
	encAccessToken, err := core.JWTEncrypter.Encrypt([]byte(rawAccessToken), encAccessKey)
	encRefreshToken, err := core.JWTEncrypter.Encrypt([]byte(rawRefreshToken), encRefreshKey)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error creating token"}
	}

	user.AccessToken = base64.StdEncoding.EncodeToString(encAccessToken)
	user.RefreshToken = base64.StdEncoding.EncodeToString(encRefreshToken)
	user.Id = dbUser.Id
	return c.JSON(http.StatusOK, user)
}
func (h *Handler) GetAccountInfo(c echo.Context) error {
	user := c.Get("user").(models.User)
	if userExists, err := h.DB.Exists(user.Username); err != nil || !userExists {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "user not found"}
	}
	dbUser, err := h.DB.GetUserByName(user.Username)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusNotFound, Message: "user not found"}
	}
	dbUser.Password = ""
	return c.JSON(http.StatusOK, dbUser)
}
func (h *Handler) RefreshToken(c echo.Context) error {
	user := c.Get("user").(*models.User)

	rawAccessToken, err := core.JWTFactory.NewToken(user.Id, "https://auth.swimply.pl/refresh-token", true)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error signing token"}
	}

	encAccessKey, err := hex.DecodeString(os.Getenv("JWT_ACCESS_SECRET"))
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error creating token"}
	}

	encAccessToken, err := core.JWTEncrypter.Encrypt([]byte(rawAccessToken), encAccessKey)
	if err != nil {
		return &echo.HTTPError{Code: http.StatusUnauthorized, Message: "error creating token"}
	}

	user.AccessToken = base64.StdEncoding.EncodeToString(encAccessToken)
	return c.JSON(http.StatusOK, user)
}