package midtransHandler

import "time"

type PaymentRequest struct {
	ID                          int
	UUID                        string
	SourceID                    string
	WebhookUrl                  string
	PaymentUrl                  string
	IsReceiveWebhook            int
	Status                      string
	ClientWebhookStatusCode     int
	ClientWebhookStatusResponse string
	MidtransWebhookPayload      string
	RequestPayload              string
	CreatedAt                   time.Time
	UpdatedAt                   time.Time
}
