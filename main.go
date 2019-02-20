package main

//! Golang(API) =>Echo + Mondgo Atlas

import (
	"fmt"
	"net/http"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/teamming/getmethod"
	"github.com/teamming/postmethod"
	"github.com/teamming/updatemethod"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Welcome to api teammate")
	})
	// Query all data
	e.GET("/read", getmethod.Getdata)
	e.POST("/register", postmethod.Postdata)
	e.PUT("/updateavatar", updatemethod.UploadImage)
	e.GET("/image/:username", getmethod.GetImage)
	e.PUT("/checkin", updatemethod.UpdateTimeAndLocation)
	e.PUT("/telno", updatemethod.UpdateTelNumber)
	e.PUT("/email", updatemethod.UpdateEmail)
	e.PUT("/team", updatemethod.UpdateTeam)
	e.GET("/sort", getmethod.SortDateAndTime)

	// Start server
	e.Logger.Fatal(e.Start(getPort()))
}
func getPort() string {
	var port = os.Getenv("PORT")
	if port == "" {
		port = "1323"
		fmt.Println("No Port In Heroku" + port)
	}
	return ":" + port

}
