package handler

import (
	"fmt"
	"jackpot/db"
	"jackpot/helper"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	YYYYMMDD = "2006-01-02"
)

type User struct {
	UserId       uint              `json:"user_id" gorm:"primaryKey"`
	Name         string            `json:"name" validate:"required, name"`
	Email        string            `json:"email" validate:"required, email"`
	Password     string            `json:"-" validate:"required, password"`
	DateOfBirth  helper.CustomDate `json:"date_of_birth" validate:"required, date_of_birth"`
	IsAdmin      bool              `json:"is_admin"`
	IsSoftBanned bool              `json:"is_soft_banned"`
}

type RegisterResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required, email"`
	Password string `json:"-" validate:"required, password"`
}

type LoginResponse struct {
	Message string `json:"message"`
}

var jwtSecret = []byte("secret")

func RegisterStudent(c echo.Context) error {
	var req User
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Request"})
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Invalid generate password"})
	}

	req.Password = string(hashPassword)

	if err := db.GormDB.Create(&req).Error; err != nil {
		fmt.Println(err.Error())
		if err.Error() == `ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)` {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "email already exist"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Database Error"})
	}

	response := RegisterResponse{
		Message: "Registered",
		Data: map[string]interface{}{
			"id":            req.UserId,
			"name":          req.Name,
			"email":         req.Email,
			"date_of_birth": req.DateOfBirth,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func LoginStudent(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Request"})
	}

	var user User

	if err := db.GormDB.Where("email = ?", req.Email).Table("students").Take(&user).Scan(&user).Error; err != nil {
		fmt.Println(err.Error())
		if err.Error() == "record not found" {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Email or Password"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error Fetching Student Data"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid Email or Password"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.UserId,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println(err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Error generating token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "successfully logged in",
		"token":   tokenString,
	})

}
