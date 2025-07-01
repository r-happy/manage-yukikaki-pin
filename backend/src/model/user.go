package model

import (
	"errors"
	"regexp"
	"time"

	"github.com/google/uuid"
)

// Userの型定義
type User struct {
	UserID       uuid.UUID `json:"user_id" gorm:"primaryKey"`
	UserName     string    `json:"user_name"`
	UserEmail    string    `json:"user_email"`
	UserPassword string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	DeletedAt    time.Time `json:"deleted_at"`
}

func isValidEmail(email string) bool {
	const emailRegexPattern = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	matched, _ := regexp.MatchString(emailRegexPattern, email)
	return matched
}

// Userの作成
func CreateUser(user *User) error {
	if !isValidEmail(user.UserEmail) {
		return errors.New("invalid email")
	}

	if u, _ := FindUserByUserEmail(user.UserEmail); u != nil {
		return errors.New("Emails is already used")
	}

	r := db.Create(user)

	if r.Error != nil {
		return r.Error
	}

	return nil
}

// Userを探す関数群 //

// UserIDを用いてUserを探す
func FindUserByUserID(user_id uuid.UUID) (*User, error) {
	var user User
	r := db.Where("user_id = ?", user_id).First(&user)
	if r.Error != nil {
		return nil, r.Error
	}
	return &user, nil
}

// UserIDを用いてUserを探す
func FindUserByUserEmail(email string) (*User, error) {
	var user User
	r := db.Where("user_email = ?", email).First(&user)
	if r.Error != nil {
		return nil, r.Error
	}
	return &user, nil
}
