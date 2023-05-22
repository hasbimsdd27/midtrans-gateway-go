package libs

import (
	"fmt"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransRequestPayload struct {
	FName    string
	LName    string
	Email    string
	Phone    string
	GrossAmt int
	OrderID  string
}

func RequestPaymentToMidtrans(input MidtransRequestPayload) (string, error) {
	snp := snap.Client{}

	env := midtrans.Sandbox

	if os.Getenv("MIDTRANS_ENV") == "production" {
		env = midtrans.Production
	}

	snp.New(os.Getenv("MIDTRANS_SERVER_KEY"), env)

	customerDetail := midtrans.CustomerDetails{}

	customerDetail.FName = input.FName
	customerDetail.Email = input.Email

	if input.LName != "" {
		customerDetail.LName = input.LName
	}
	if input.Phone != "" {
		customerDetail.Phone = input.Phone
	}

	transactionDetail := midtrans.TransactionDetails{}
	transactionDetail.GrossAmt = int64(input.GrossAmt)
	transactionDetail.OrderID = input.OrderID

	req := &snap.Request{
		TransactionDetails: transactionDetail,
		CustomerDetail:     &customerDetail,
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	snapRes, err := snp.CreateTransaction(req)

	if err != nil {
		fmt.Println("12345", os.Getenv("MIDTRANS_SERVER_KEY"))
		return "", err
	}

	return snapRes.RedirectURL, nil

}
