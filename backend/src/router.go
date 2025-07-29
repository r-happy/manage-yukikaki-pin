package main

import (
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/r-happy/yukikaki-system/src/handler"
)

func newRouter() *echo.Echo {
	// godotenv
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	e := echo.New()

	// Global middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Public routes
	e.GET("/", func(c echo.Context) error { return c.JSON(http.StatusOK, "Hello world") })
	e.POST("/signup", handler.SignUp)
	e.POST("/signin", handler.SignIn)

	// jwt middleware
	r := e.Group("/api")
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(handler.JwtCustomClaims)
		},
		SigningKey: []byte(os.Getenv("SECRET_KEY")),
	}
	r.Use(echojwt.WithConfig(config))

	// api routes with jwt middleware
	r.GET("/me", handler.GetUser)
	r.GET("/me/groups", handler.GetGroupsByUser)

	// group
	r.POST("/groups", handler.AddGroup)
	r.GET("/groups/:groupID", handler.GetGroup)

	// group member
	r.POST("/groups/:groupID/add-members", handler.AddGroupMemberByAdmin)
	r.POST("/groups/:groupID/join-requests", handler.AddGroupMemberByAdmin)
	// pin type
	r.POST("/groups/:groupID/pin-types", handler.AddPinType)
	r.GET("/groups/:groupID/pin-types", handler.GetPinTypesByGroupID)

	// pin
	r.POST("/groups/:groupID/pins", handler.AddPinByMember)
	r.GET("/groups/:groupID/pins", handler.GetPinsByGroupID)

	return e
}
