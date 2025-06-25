package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetUser(c echo.Context) error {
	user, err := UserFromToken(c)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	// user.UserPassword = ""

	return c.JSON(http.StatusOK, user)
}
