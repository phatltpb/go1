package controller

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/phatltb/gotraining/config"
	"github.com/phatltb/gotraining/middlewares"
	"github.com/phatltb/gotraining/model"
)

type token struct {
	token string `json:"string"`
}

func GetAuth(c echo.Context) error {
	db := config.SetupDatabaseConnection()
	authRequest := []model.AuthRequest{}
	if err := db.Find(&authRequest).Error; err != nil {
		return nil
	}
	return c.JSON(http.StatusOK, authRequest)
}

func CreateAuth(c echo.Context) error {
	db := config.SetupDatabaseConnection()
	req := new(model.AuthRequest)
	c.Bind(req)
	authRequest := model.AuthRequest{
		PhoneNumber: req.PhoneNumber,
		Status:      "pending",
		UID:         rand.Uint64(),
		Expired_at:  time.Now().Add(time.Hour * 72).Unix(),
	}
	db.Create(&authRequest)
	return c.JSON(http.StatusCreated, authRequest)
}

func CheckAuthRequest(c echo.Context) error {
	db := config.SetupDatabaseConnection()
	req := new(model.AuthRequest)
	c.Bind(req)
	// log.Println(req.UID)
	authRequest := model.AuthRequest{
		PhoneNumber: req.PhoneNumber,
		UID:         req.UID,
	}
	if err := db.First(&authRequest).Error; err != nil {
		return nil
	}
	if authRequest.Status == "authenticated" || authRequest.Status == "finnish" {
		authRequest.Status = "finnish"
	} else {
		timeStop := make(map[string]int64)
		timeStop["expired_at"] = authRequest.Expired_at
		return c.JSON(http.StatusOK, timeStop)
	}

	db.Save(&authRequest)
	t := middlewares.CreateJWT(authRequest.PhoneNumber, authRequest.Status, authRequest.Expired_at)
	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func DecryptionJWTPhone(c echo.Context) error {
	db := config.SetupDatabaseConnection()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	phone := claims["phone_number"].(string)
	authRequest := model.AuthRequest{
		PhoneNumber: phone,
	}
	if err := db.First(&authRequest).Error; err != nil {
		return nil
	}
	if authRequest.Status != "authenticated" {
		authRequest.Status = "authenticated"
	}
	db.Save(&authRequest)

	return c.JSON(http.StatusOK, authRequest)
}

func DecryptionJWTStatus(c echo.Context) error {
	db := config.SetupDatabaseConnection()
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	phone := claims["status"].(string)
	authRequest := model.AuthRequest{
		PhoneNumber: phone,
	}
	if err := db.First(&authRequest).Error; err != nil {
		return nil
	}
	if authRequest.Status != "request" {
		authRequest.Status = "request"
	}
	db.Save(&authRequest)

	return c.JSON(http.StatusOK, authRequest)
}
