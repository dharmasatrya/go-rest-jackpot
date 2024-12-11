package routes

import (
	"jackpot/handler"

	echojwt "github.com/labstack/echo-jwt"
	"github.com/labstack/echo/v4"
)

func Init(e *echo.Echo) {

	u := e.Group("/users")
	u.POST("/register", handler.RegisterStudent)
	u.POST("/login", handler.LoginStudent)

	r := e.Group("/records")
	r.Use(echojwt.JWT([]byte("secret")))
	// r.GET("/wallet", ) check user wallet
	// r.GET("/stats", ) check the user game records

	t := e.Group("/topup")
	t.Use(echojwt.JWT([]byte("secret")))
	// t.GET("", ) get the available topup amount
	// t.POST("topup", ) topup
	// t.GET("/history") get a user topup history

	g := e.Group("/games")
	g.Use(echojwt.JWT([]byte("secret")))
	// g.POST("/play", )

	a := e.Group("/godmode")
	a.Use(echojwt.JWT([]byte("secret")))
	// a.GET("/winrates", ) display users winrate
	// a.GET("/curse")
	// a.POST("/curse", ) apply curse to a user
	// a.DELETE("/")
}
