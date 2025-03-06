package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"testProject/internal/services"
)

func LoginHandler(c echo.Context) error {
	tokens, err := services.AuthenticateUser(c)
	if err != nil {
		return err
	}

	jsonResponse, err := json.Marshal(tokens)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}
	fmt.Println(string(jsonResponse))

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(http.StatusOK)

	if _, err := c.Response().Write(jsonResponse); err != nil {
		return err
	}
	return nil
}

/*func RegisterHandler(c echo.Context) error {
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
}
*/
