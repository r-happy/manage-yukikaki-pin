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
	if err := godotenv.Load(); err != nil {
		log.Printf("[warn] .env file not found, continuing with existing environment: %v", err)
	}

	secretKey, ok := os.LookupEnv("SECRET_KEY")
	if !ok || secretKey == "" {
		log.Fatal("SECRET_KEY environment variable must be set")
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
	e.GET("/public/groups/:groupID/pins", handler.GetPublicPinsByGroupID)

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
	r.GET("/groups/:groupID/members", handler.GetGroupMembers)
	r.POST("/groups/:groupID/add-members", handler.AddGroupMemberByAdmin)
	r.POST("/groups/:groupID/join-requests", handler.AddGroupMemberByAdmin)
	r.PUT("/groups/:groupID/members/:memberID", handler.UpdateGroupMember)
	r.DELETE("/groups/:groupID/members/:memberID", handler.DeleteGroupMember)
	// pin type
	r.POST("/groups/:groupID/pin-types", handler.AddPinType)
	r.GET("/groups/:groupID/pin-types", handler.GetPinTypesByGroupID)

	// pin
	r.POST("/groups/:groupID/pins", handler.AddPinByMember)
	r.GET("/groups/:groupID/pins", handler.GetPinsByGroupID)
	r.PUT("/groups/:groupID/pins/:pinID", handler.UpdatePinByMember)
	r.DELETE("/groups/:groupID/pins/:pinID", handler.DeletePinByMember)

	return e
}
