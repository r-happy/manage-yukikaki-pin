package main

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r-happy/yukikaki-system/src/handler"
)

func newRouter() *echo.Echo {
	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes
	e.GET("/", func(c echo.Context) error { return c.JSON(http.StatusOK, "Hello world") })
	e.POST("/signup", handler.SignUp)
	e.GET("/signin", handler.SignIn)

	// jwt middleware
	r := e.Group("/api")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte("secret"),
	}
	r.Use(echojwt.WithConfig(config))

	// api routes
	r.GET("/user", handler.GetUser)

	return e
}
