package main

import (
	"fmt"
	"log"
	"midtrans-adapter-go/handler"
	"midtrans-adapter-go/libs"
	"midtrans-adapter-go/midtransHandler"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := libs.ConnectMysql()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	midtransRepository := midtransHandler.NewRepository(db)
	midtransService := midtransHandler.NewService(midtransRepository)
	handlerMidtrans := handler.NewMidtransHandler(midtransService)

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/create", handlerMidtrans.CreatePayment)
	r.POST("/webhook", handlerMidtrans.HandleWebhook)

	r.Run(fmt.Sprintf("127.0.0.1:%s", os.Getenv("PORT")))
}
