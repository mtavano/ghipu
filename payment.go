package ghipu

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type PaymentService struct {
	client *httpClient
}

func NewPaymentService(secret, receiverID string) *PaymentService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &PaymentService{&client}
}

func (ps *PaymentService) Payment(id string) (*PaymentResponse, error) {
	resp, err := ps.client.Get("/payments/"+id, nil)
	if err != nil {
		return nil, err
	}
	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (ps *PaymentService) CreatePayment(p *Payment) (*PaymentResponse, error) {
	resp, err := ps.client.PostForm("/payments", p.Params())
	if err != nil {
		return nil, err
	}
	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (ps *PaymentService) Refund(id string) (*PaymentResponse, error) {
	resp, err := ps.client.PostForm("/payments/"+id+"/refunds", nil)
	if err != nil {
		return nil, err
	}

	var pr PaymentResponse
	if err := unmarshalJSON(resp.Body, &pr); err != nil {
		return nil, err
	}
	return &pr, nil
}

func (ps *PaymentService) PaymentStatus(notificationToken string) (*SuccessResponse, error) {
	values := url.Values{"notification_token": {notificationToken}}

	resp, err := ps.client.Get("/payments?"+values.Encode(), values)
	if err != nil {
		return nil, err
	}

	var sr SuccessResponse
	if err := unmarshalJSON(resp.Body, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

func unmarshalJSON(r io.ReadCloser, v interface{}) error {
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}

// Payment represents the payment form requires by khipu to make a payment POST
type Payment struct {
	Amount      float64   `json:"amount"`
	ExpiresDate time.Time `json:"expires_date,omitempty"`

	NotifyURL  string `json:"notify_url,omitempty"`
	CancelURL  string `json:"cancel_url,omitempty"`
	ContactURL string `json:"contact_url,omitempty"`

	BankID                       string `json:"bank_id,omitempty"`
	Body                         string `json:"body,omitempty"`
	Currency                     string `json:"currency"`
	Custom                       string `json:"custom,omitempty"`
	FixedPayerPersonalIdentifier string `json:"fixed_payer_personal_identifier,omitempty"`
	IntegratorFee                string `json:"integrator_fee,omitempty"`
	NotifyAPIVersion             string `json:"notify_api_verion,omitempty"`
	PayerEmail                   string `json:"payer_email,omitempty"`
	PayerName                    string `json:"payer_name,omitempty"`
	PictureURL                   string `json:"picture_url,omitempty"`
	ResponsibleUserEmail         string `json:"responsible_user_email,omitempty"`
	ReturnURL                    string `json:"return_url,omitempty"`
	SendEmail                    bool   `json:"send_email,omitempty"`
	SendReminders                bool   `json:"send_reminders,omitempty"`

	Subject       string `json:"subject"`
	TransactionID string `json:"transaction_id,omitempty"`
}

// Params returns a map used to sign the requests
func (p *Payment) Params() url.Values {
	form := url.Values{
		"subject":  {p.Subject},
		"currency": {p.Currency},
		"amount":   {fmt.Sprintf("%.4f", p.Amount)},

		"fixed_payer_personal_identifier": {p.FixedPayerPersonalIdentifier},
		"responsible_user_email":          {p.ResponsibleUserEmail},
		"transaction_id":                  {p.TransactionID},
		"bank_id":                         {p.BankID},

		"custom": {p.Custom},
		"body":   {p.Body},

		"return_url":  {p.ReturnURL},
		"cancel_url":  {p.CancelURL},
		"picture_url": {p.PictureURL},
		"notify_url":  {p.NotifyURL},
		"contact_url": {p.ContactURL},

		"notify_api_version": {p.NotifyAPIVersion},
		"integrator_fee":     {p.IntegratorFee},
	}

	if !p.ExpiresDate.IsZero() {
		form.Set("expires_date", p.ExpiresDate.Format(time.RFC3339))
	}

	if p.SendEmail {
		form.Set("send_email", fmt.Sprint(p.SendEmail))
		form.Set("payer_name", p.PayerName)
		form.Set("payer_email", p.PayerEmail)
		form.Set("send_reminders", fmt.Sprint(p.SendReminders))
	}

	return form
}

// PaymentResponse represets the information returnes by khipu's api after a payment action
type PaymentResponse struct {
	AppURL                string    `json:"app_url,omitempty"`
	AttachmentURLS        []string  `json:"attachment_urls,omitempty"`
	Bank                  string    `json:"bank,omitempty"`
	BankAccountNumber     string    `json:"bank_account_number,omitempty"`
	ConciliationDate      time.Time `json:"conciliation_date,omitempty"`
	NotificationToken     string    `json:"notification_token,omitempty"`
	OutOfDateConciliation bool      `json:"aout_of_date_conciliation,omitempty"`
	Payment
	PaymentID             string `json:"payment_id,omitempty"`
	PaymentURL            string `json:"payment_url,omitempty"`
	PersonalIdentifier    string `json:"personal_identifier,omitempty"`
	ReadyForTerminal      bool   `json:"ready_for_terminal,omitempty"`
	ReceiptURL            string `json:"receipt_url,omitempty"`
	ReceiverID            string `json:"reveiver_id,omitempty"`
	SimplifiedTransferURL string `json:"simplified_tranfer_url,omitempty"`
	Status                string `json:"status,omitempty"`
	StatusDetail          string `json:"status_detail,omitempty"`
	TransferURL           string `json:"transfer_url,omitempty"`
}

// SuccessResponse represents a success response defined by khipu
type SuccessResponse struct {
	Message string `json:"message"`
}
