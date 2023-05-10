package controller

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prayogatriady/sawer-transaction/middleware"
	"github.com/prayogatriady/sawer-transaction/model"
	"github.com/prayogatriady/sawer-transaction/service"
)

type TransactionContInterface interface {
	Sawer(c *gin.Context)
	PaymentNotification(c *gin.Context)
}

type TransactionController struct {
	TransactionService service.TransactionServInterface
}

func NewTransactionController(trService service.TransactionServInterface) TransactionContInterface {
	return &TransactionController{
		TransactionService: trService,
	}
}

func (tc *TransactionController) PaymentNotification(c *gin.Context) {
	ctx := context.Background()

	var notifPayload model.NotificationPayload
	if err := c.BindJSON(&notifPayload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "500 - INTERNAL SERVER ERROR",
			"message": err.Error(),
		})
		return
	}

	fmt.Printf("%+v\n", notifPayload)

	if err := tc.TransactionService.PaymentNotification(ctx, &notifPayload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "500 - INTERNAL SERVER ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "200 - STATUS OK",
		"message": "Notif received",
	})
}

func (tc *TransactionController) Sawer(c *gin.Context) {
	ctx := context.Background()

	// get payload from token
	userId, err := middleware.ExtractToken(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": "401 - Unauthorized",
			"msg":    "Unauthorized - Missing JWT Token",
		})
		return
	}

	var trRequest model.TransactionRequest
	if err := c.BindJSON(&trRequest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "500 - INTERNAL SERVER ERROR",
			"message": err.Error(),
		})
		return
	}

	trRequest.UserId = userId

	trResponse, err := tc.TransactionService.Sawer(ctx, &trRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "500 - INTERNAL SERVER ERROR",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "200 - STATUS OK",
		"message": "Sawer succeed",
		"body":    trResponse,
	})
}
