package services

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"testProject/internal/models"
	"testProject/internal/repository"
	"testProject/pkg/utils"
	"time"
)

func AuthenticateUser(c echo.Context) (map[string]string, error) {
	user := c.Get("user").(models.LoginRequest)

	retrievedUser, err := repository.LocalService.GetUserAndPassword(user.Username)
	if err != nil {
		return map[string]string{}, c.String(http.StatusUnauthorized, "invalid username or password")
	}

	if err := utils.Argon.Verify([]byte(user.Password), retrievedUser.PasswordHash); err != nil {
		return map[string]string{}, c.String(http.StatusUnauthorized, "invalid username or password")
	}

	accessToken, err := GenerateJWT(retrievedUser.ID, map[string]interface{}{
		"iss": "https://swimply.pl/api/v2/login",
		"exp": time.Now().Add(time.Hour * 2).Unix(),
		"iat": time.Now().Unix(),
		"typ": "access_token",
	})
	if err != nil {
		return nil, c.String(http.StatusInternalServerError, "something went wrong during token generation")
	}

	refreshToken, err := GenerateJWT(retrievedUser.ID, map[string]interface{}{
		"iss": "https://swimply.pl/api/v2/login",
		"exp": time.Now().Add(time.Hour * 72).Unix(),
		"iat": time.Now().Unix(),
		"typ": "refresh_token",
	})
	if err != nil {
		return map[string]string{}, c.String(http.StatusInternalServerError, "something went wrong during token generation")
	}

	tokensMap := map[string]string{}
	tokensMap["access_token"] = accessToken
	tokensMap["refresh_token"] = refreshToken
	tokensMap["user_id"] = retrievedUser.ID.Hex()

	return tokensMap, nil
}
