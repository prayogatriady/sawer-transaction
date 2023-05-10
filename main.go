package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/midtrans/midtrans-go/snap"
	"github.com/prayogatriady/sawer-transaction/controller"
	"github.com/prayogatriady/sawer-transaction/db"
	"github.com/prayogatriady/sawer-transaction/repository"
	"github.com/prayogatriady/sawer-transaction/service"
)

func main() {
	// set default environment
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "9001"
	}
	DB_USER := os.Getenv("DB_USER")
	if DB_USER == "" {
		DB_USER = "root"
	}
	DB_PASSWORD := os.Getenv("DB_PASSWORD")
	if DB_PASSWORD == "" {
		DB_PASSWORD = "root"
	}
	DB_HOST := os.Getenv("DB_HOST")
	if DB_HOST == "" {
		DB_HOST = "127.0.0.1"
	}
	DB_PORT := os.Getenv("DB_PORT")
	if DB_PORT == "" {
		DB_PORT = "3306"
	}
	DB_NAME := os.Getenv("DB_NAME")
	if DB_NAME == "" {
		DB_NAME = "sawer"
	}

	db, err := db.NewConnectDB(
		DB_USER,
		DB_PASSWORD,
		DB_HOST,
		DB_PORT,
		DB_NAME,
	).InitMySQL()
	if err != nil {
		log.Fatal(err)
	}

	trRepo := repository.NewTransactionRepository(db)
	trServ := service.NewTransactionService(trRepo, snap.Client{})
	trCont := controller.NewTransactionController(trServ)

	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/sawer", trCont.Sawer)
		api.POST("/midtrans/notification", trCont.PaymentNotification)
	}

	log.Fatal(r.Run(":" + PORT))
}
