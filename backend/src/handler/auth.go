package handler

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/r-happy/yukikaki-system/src/model"
)

// jwtの設定 //
type JwtCustomClaims struct {
	UserId uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

// SignUp処理 //
// SignUpリクエストに必要な型
type signUpRequest struct {
	UserEmail    string `form:"user_email"`
	UserPassword string `form:"user_password"`
	UserName     string `form:"user_name"`
}

// SignUpのメイン処理
func SignUp(c echo.Context) error {
	req := new(signUpRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}
	// すべてのフィールドが埋まってるかどうか
	if err := ValidateStruct(req); err != nil {
		return err
	}

	createdUserID := uuid.New()
	sha256Password := sha256.Sum256([]byte(req.UserPassword))
	user := &model.User{
		UserID:       createdUserID,
		UserName:     req.UserName,
		UserEmail:    req.UserEmail,
		UserPassword: fmt.Sprintf("%x", sha256Password),
	}

	if err := model.CreateUser(user); err != nil {
		return c.JSON(http.StatusInternalServerError, "Error")
	}

	user.UserPassword = ""

	return c.JSON(http.StatusOK, user)
}

// SignIn処理 //
// SignInリクエストに必要な型
type signInRequest struct {
	UserEmail    string `form:"email"`
	UserPassword string `form:"password"`
}

// SignInのメイン処理
func SignIn(c echo.Context) error {
	req := new(signInRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid Request")
	}
	if err := ValidateStruct(req); err != nil {
		return err
	}

	// 認証処理
	sha256Password := sha256.Sum256([]byte(req.UserPassword))
	user, err := model.FindUserByUserEmail(req.UserEmail)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error")
	}
	if user == nil {
		return c.JSON(http.StatusUnauthorized, "Invalid Email or Password")
	}
	if user.UserPassword != fmt.Sprintf("%x", sha256Password) {
		return c.JSON(http.StatusUnauthorized, "Invalid Email or Password")
	}

	// jwtの発行
	claims := &JwtCustomClaims{
		user.UserID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 3)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{"token": t})
}
