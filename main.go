package main

import (
	config "bot-whatsapp/config"
	service "bot-whatsapp/service"
	"os"
	"strings"

	"github.com/go-playground/validator"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mau.fi/whatsmeow"
)

func init() {
	godotenv.Load()
	config.InitSqlite()
}

func main() {

	WhatsMeowService := service.WhatsMeowService{
		DB:     config.DB,
		Client: &whatsmeow.Client{},
	}
	go WhatsMeowService.StartClient()

	e := echo.New()
	e.Validator = &config.CustomValidator{Validator: validator.New()}
	// Middleware CORS for all request (allow all)
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Origin, Content-Length, Content-Type, Authorization")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Max-Age", "12h")
			return next(c)
		}
	})
	// Allow index.html in folder public to be accessed
	e.Static("/", "public")

	e.GET("/get-qr", func(c echo.Context) error {
		qr, jid := WhatsMeowService.GetQR()
		return c.JSON(200, map[string]string{
			"qr":  qr,
			"jid": jid,
		})
	})

	e.POST("/message", func(c echo.Context) error {
		payload := struct {
			Phone   string `json:"phone" form:"phone" query:"phone" validate:"required"`
			Message string `json:"message" form:"message" query:"message" validate:"required"`
		}{}
		if err := c.Bind(&payload); err != nil {
			return c.JSON(400, map[string]string{
				"message": err.Error(),
			})
		}
		if err := c.Validate(&payload); err != nil {
			stringme := strings.Split(err.Error(), "Key:")
			stringme = stringme[1:]
			return c.JSON(400, map[string]interface{}{
				"message": "Bad request",
				"field":   stringme,
			})
		}

		err := WhatsMeowService.SendTextMessage(payload.Phone, payload.Message)
		if err != nil {
			return c.JSON(500, map[string]string{
				"message": err.Error(),
			})
		}

		return c.JSON(200, map[string]string{
			"message": "success",
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	e.Logger.Fatal(e.Start(":" + port))
}
