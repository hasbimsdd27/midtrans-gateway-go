package libs

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type CallWebhookResult struct {
	Status   int
	Response string
}

type CallWebhookInput struct {
	WebhookUrl      string `json:"-"`
	TransactionCode string `json:"transaction_code"`
	Status          string `json:"status"`
}

func CallWebhook(input CallWebhookInput) (CallWebhookResult, error) {
	result := CallWebhookResult{}

	postBody, _ := json.Marshal(input)

	postBodyBuffer := bytes.NewBuffer(postBody)

	resp, err := http.Post(input.WebhookUrl, "application/json", postBodyBuffer)

	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return result, err
	}

	result.Status = resp.StatusCode
	result.Response = string(bodyBytes)

	return result, nil

}
