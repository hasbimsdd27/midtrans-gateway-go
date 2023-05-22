package midtransHandler

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"midtrans-adapter-go/libs"
	"os"
	"strings"

	"github.com/google/uuid"
)

type Service interface {
	CreateTransaction(input RequestPayload) (PaymentRequest, error)
	UpdatePayment(input WebhookPayload) (PaymentRequest, error)
}

type service struct {
	paymentRepository Repository
}

func NewService(paymentRepository Repository) *service {
	return &service{paymentRepository}
}

func (s *service) CreateTransaction(input RequestPayload) (PaymentRequest, error) {

	paymentRequest := PaymentRequest{}

	json, err := json.Marshal(input)

	if err != nil {
		return paymentRequest, err
	}

	UUID := uuid.New().String()

	paymentRequest.UUID = UUID
	paymentRequest.SourceID = input.TransactionDetail.OrderID
	paymentRequest.WebhookUrl = input.WebhookUrl
	paymentRequest.RequestPayload = string(json)
	paymentRequest.Status = "PENDING"

	newPaymentRequest, err := s.paymentRepository.Save(paymentRequest)

	if err != nil {
		return newPaymentRequest, err
	}

	paymentRequestPayload := libs.MidtransRequestPayload{}
	paymentRequestPayload.Email = input.CustomerDetail.Email
	paymentRequestPayload.FName = input.CustomerDetail.FirstName
	paymentRequestPayload.LName = input.CustomerDetail.LastName
	paymentRequestPayload.Phone = input.CustomerDetail.Phone
	paymentRequestPayload.GrossAmt = input.TransactionDetail.GrossAmount
	paymentRequestPayload.OrderID = UUID

	paymentUrl, err := libs.RequestPaymentToMidtrans(paymentRequestPayload)

	if err != nil {
		return newPaymentRequest, err
	}

	newPaymentRequest.PaymentUrl = paymentUrl

	newPaymentRequest, err = s.paymentRepository.Update(newPaymentRequest)

	if err != nil {
		return newPaymentRequest, err
	}

	return newPaymentRequest, nil

}

func (s *service) UpdatePayment(input WebhookPayload) (PaymentRequest, error) {

	paymentRequest, err := s.paymentRepository.FindByUUID(input.OrderID)

	if err != nil {
		return paymentRequest, err
	}

	signatureKey := fmt.Sprintf("%s%s%s%s", input.OrderID, input.StatusCode, input.GrossAmount, os.Getenv("MIDTRANS_SERVER_KEY"))

	hashedSignatureKey := sha512.New()
	hashedSignatureKey.Write([]byte(signatureKey))

	ownSignatureKey := fmt.Sprintf("%x", hashedSignatureKey.Sum(nil))

	if ownSignatureKey != input.SignatureKey {
		return paymentRequest, errors.New("invalid signature key")
	}

	if paymentRequest.Status == "SUCCESS" {
		return paymentRequest, errors.New("transaction already paid")
	}

	if input.TransactionStatus == "capture" {
		if input.FraudStatus == "challange" {
			paymentRequest.Status = "CHALLENGE"
		} else if input.FraudStatus == "accept" {
			paymentRequest.Status = "SUCCESS"
		}
	} else if input.TransactionStatus == "settlement" {
		paymentRequest.Status = "SUCCESS"
	} else if input.TransactionStatus == "cancel" ||
		input.TransactionStatus == "deny" ||
		input.TransactionStatus == "expire" {
		paymentRequest.Status = strings.ToUpper(input.TransactionStatus)
	} else if input.TransactionStatus == "pending" {
		paymentRequest.Status = "PENDING"
	}

	json, err := json.Marshal(input)

	if err != nil {
		return paymentRequest, err
	}
	paymentRequest.IsReceiveWebhook = 1
	paymentRequest.MidtransWebhookPayload = string(json)

	callWebhookInputParams := libs.CallWebhookInput{}
	callWebhookInputParams.WebhookUrl = paymentRequest.WebhookUrl
	callWebhookInputParams.TransactionCode = paymentRequest.SourceID
	callWebhookInputParams.Status = paymentRequest.Status

	webhookResponse, err := libs.CallWebhook(callWebhookInputParams)
	if err != nil {
		return paymentRequest, err
	}
	paymentRequest.ClientWebhookStatusCode = webhookResponse.Status
	paymentRequest.ClientWebhookStatusResponse = webhookResponse.Response

	paymentRequest, err = s.paymentRepository.Update(paymentRequest)

	if err != nil {
		return paymentRequest, err
	}

	return paymentRequest, nil
}
