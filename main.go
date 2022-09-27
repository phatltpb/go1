package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/phatltb/gotraining/config"
	"github.com/phatltb/gotraining/controller"
)

func main() {
	e := echo.New()
	config.SetupDatabaseConnection()
	isLogin := middleware.JWT([]byte("secret"))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	ser1 := e.Group("/server1")
	ser1.GET("/", controller.GetAuth)
	ser1.POST("/create", controller.CreateAuth)
	ser1.POST("/checkauth", controller.CheckAuthRequest)
	ser2 := e.Group("/server2")
	ser2.POST("/checkJWTp", controller.DecryptionJWTPhone, isLogin)
	ser2.POST("/checkJWTs", controller.DecryptionJWTStatus, isLogin)
	e.Logger.Fatal(e.Start(":3080"))
}
