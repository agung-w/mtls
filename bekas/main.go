package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	// Inisialisasi Echo
	e := echo.New()

	// Route untuk Ping
	e.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	// Jalankan server pada port 8080
	e.Logger.Fatal(e.Start(":8080"))
}
