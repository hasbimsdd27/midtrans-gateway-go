package handler

import (
	"midtrans-adapter-go/midtransHandler"
	"net/http"

	"github.com/gin-gonic/gin"
)

type handlerStruct struct {
	service midtransHandler.Service
}

func NewMidtransHandler(service midtransHandler.Service) *handlerStruct {
	return &handlerStruct{service}
}

func (h *handlerStruct) CreatePayment(c *gin.Context) {
	var input midtransHandler.RequestPayload

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	newPaymentRequest, err := h.service.CreateTransaction(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"payment_url": newPaymentRequest.PaymentUrl,
		},
	})

}

func (h *handlerStruct) HandleWebhook(c *gin.Context) {
	var input midtransHandler.WebhookPayload

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}
	_, err = h.service.UpdatePayment(input)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
	})

}
