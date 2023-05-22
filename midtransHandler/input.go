package midtransHandler

type RequestPayloadTransactionDetail struct {
	OrderID     string `json:"order_id" binding:"required"`
	GrossAmount int    `json:"gross_amount" binding:"required"`
}

type RequestPayloadCustomerDetail struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone"`
}

type RequestPayload struct {
	TransactionDetail RequestPayloadTransactionDetail `json:"transaction_details" binding:"required"`
	CustomerDetail    RequestPayloadCustomerDetail    `json:"customer_detail" binding:"required"`
	WebhookUrl        string                          `json:"webhook_url" binding:"required"`
	GeneratedID       string
	TransactionKey    string `json:"transaction_key" binding:"required"`
}

type WebhookPayload struct {
	OrderID                string `json:"order_id" binding:"required"`
	TransactionStatus      string `json:"transaction_status" binding:"required"`
	FraudStatus            string `json:"fraud_status" binding:"required"`
	StatusCode             string `json:"status_code" binding:"required"`
	GrossAmount            string `json:"gross_amount" binding:"required"`
	SignatureKey           string `json:"signature_key" binding:"required"`
	TransactionTime        string `json:"transaction_time" `
	TransactionID          string `json:"transaction_id"`
	StatusMessage          string `json:"status_message"`
	PaymentType            string `json:"payment_type"`
	MerchantID             string `json:"merchant_id"`
	MaskedCard             string `json:"masked_card"`
	Eci                    string `json:"eci"`
	Currency               string `json:"currency"`
	ChannelResponseMessage string `json:"channel_response_message"`
	ChannelResponseCode    string `json:"channel_response_code"`
	CardType               string `json:"card_type"`
	Bank                   string `json:"bank"`
	ApprovalCode           string `json:"approval_code"`
}
