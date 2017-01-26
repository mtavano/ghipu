package ghipu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type PaymentsService struct {
	client *httpClient
}

func NewPaymentService(secret string, receiverID int) *PaymentsService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &PaymentsService{&client}
}

func (ps *PaymentsService) ReceiverID() int { return ps.client.recid }

// Payment returns informatin of the payment with the given id.
func (ps *PaymentsService) Payment(id string) (*PaymentResponse, error) {
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

func (ps *PaymentsService) PaymentStatus(notificationToken string) (*PaymentResponse, error) {
	values := url.Values{"notification_token": {notificationToken}}

	resp, err := ps.client.Get("/payments?"+values.Encode(), values)
	if err != nil {
		return nil, err
	}

	var payr PaymentResponse
	if err := unmarshalJSON(resp.Body, &payr); err != nil {
		return nil, err
	}

	if payr.ReceiverID != ps.ReceiverID() {
		return nil, errors.New("ghipu: invalid receiver id")
	}

	return &payr, nil
}

// CreatePayment creates a new payment and returns the URLs to complete the payment.
func (ps *PaymentsService) CreatePayment(p *Payment) (*PaymentCreateResponse, error) {
	resp, err := ps.client.PostForm("/payments", p.Params())
	if err != nil {
		return nil, err
	}

	var pcr PaymentCreateResponse
	if err := unmarshalJSON(resp.Body, &pcr); err != nil {
		return nil, err
	}
	return &pcr, nil
}

func (ps *PaymentsService) DeletePayment(id string) (*SuccessResponse, error) {
	resp, err := ps.client.Delete("/payments/"+id, nil)
	if err != nil {
		return nil, err
	}

	var sr SuccessResponse
	if err := unmarshalJSON(resp.Body, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Refund refunds in whole or in part the amount of a payment. This operation can only be
// carried out in the stores that collect in khipu account and before the surrender of the
// corresponding funds.
// If the amount is zero, then the whole amount of the payment is refunded.
func (ps *PaymentsService) Refund(id string, amount float64) (*SuccessResponse, error) {
	var values url.Values
	if amount > 0 {
		values = url.Values{"amount": {fmt.Sprintf("%.4f", amount)}}
	}

	resp, err := ps.client.PostForm("/payments/"+id+"/refunds", values)
	if err != nil {
		return nil, err
	}

	var sr SuccessResponse
	if err := unmarshalJSON(resp.Body, &sr); err != nil {
		return nil, err
	}
	return &sr, nil
}

// Payment represents the payment form requires by khipu to make a payment POST
type Payment struct {
	Amount      float64   `json:"amount"`
	ExpiresDate time.Time `json:"expires_date"`

	NotifyURL  string `json:"notify_url"`
	CancelURL  string `json:"cancel_url"`
	ContactURL string `json:"contact_url"`

	BankID                       string `json:"bank_id"`
	Body                         string `json:"body"`
	Currency                     string `json:"currency"`
	Custom                       string `json:"custom"`
	FixedPayerPersonalIdentifier string `json:"fixed_payer_personal_identifier"`
	IntegratorFee                string `json:"integrator_fee"`
	NotifyAPIVersion             string `json:"notify_api_verion"`
	PayerEmail                   string `json:"payer_email"`
	PayerName                    string `json:"payer_name"`
	PictureURL                   string `json:"picture_url"`
	ResponsibleUserEmail         string `json:"responsible_user_email"`
	ReturnURL                    string `json:"return_url"`

	SendEmail     bool `json:"send_email"`
	SendReminders bool `json:"send_reminders"`

	Subject       string `json:"subject"`
	TransactionID string `json:"transaction_id"`
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

type PaymentCreateResponse struct {
	PaymentID             string `json:"payment_id"`
	PaymentURL            string `json:"payment_url"`
	SimplifiedTransferURL string `json:"simplified_tranfer_url"`
	TransferURL           string `json:"transfer_url"`
	AppURL                string `json:"app_url"`
	ReadyForTerminal      bool   `json:"ready_for_terminal"`
}

// PaymentResponse represets the information returnes by khipu's api after a payment action
type PaymentResponse struct {
	PaymentID             string    `json:"payment_id"`
	PaymentURL            string    `json:"payment_url"`
	SimplifiedTransferURL string    `json:"simplified_tranfer_url"`
	TransferURL           string    `json:"transfer_url"`
	AppURL                string    `json:"app_url"`
	ReadyForTerminal      bool      `json:"ready_for_terminal"`
	NotificationToken     string    `json:"notification_token"`
	ReceiverID            int       `json:"receiver_id"`
	ConciliationDate      time.Time `json:"conciliation_date"`
	Subject               string    `json:"subject"`
	Amount                float64   `json:"amount"`
	Status                string    `json:"status"`
	StatusDetail          string    `json:"status_detail"`
	Body                  string    `json:"body"`
	PictureURL            string    `json:"picture_url"`
	ReceiptURL            string    `json:"receipt_url"`
	NotifyURL             string    `json:"notify_url"`
	ReturnURL             string    `json:"return_url"`
	CancelURL             string    `json:"cancel_url"`
	ContactURL            string    `json:"contact_url"`
	NotifyAPIVersion      string    `json:"notify_api_verion"`
	ExpiresDate           time.Time `json:"expires_date"`
	AttachmentURLS        []string  `json:"attachment_urls"`
	Bank                  string    `json:"bank"`
	BankID                string    `json:"bank_id"`
	PayerName             string    `json:"payer_name"`
	PayerEmail            string    `json:"payer_email"`
	PersonalIdentifier    string    `json:"personal_identifier"`
	BankAccountNumber     string    `json:"bank_account_number"`
	OutOfDateConciliation bool      `json:"out_of_date_conciliation"`
	TransactionID         string    `json:"transaction_id"`
	Custom                string    `json:"custom"`
	ResponsibleUserEmail  string    `json:"responsible_user_email"`
	SendReminders         bool      `json:"send_reminders"`
	SendEmail             bool      `json:"send_email"`
	PaymentMethod         string    `json:"payment_method"`
}

type floatString float64

func (fs *floatString) UnmarshalJSON(b []byte) error {
	str := string(bytes.Trim(b, `"`))

	amount, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return err
	}

	*fs = floatString(amount)

	return nil
}

func (p *PaymentResponse) UnmarshalJSON(b []byte) error {
	raw := struct {
		Amount                floatString `json:"amount"`
		PaymentID             *string     `json:"payment_id"`
		PaymentURL            *string     `json:"payment_url"`
		SimplifiedTransferURL *string     `json:"simplified_tranfer_url"`
		TransferURL           *string     `json:"transfer_url"`
		AppURL                *string     `json:"app_url"`
		ReadyForTerminal      *bool       `json:"ready_for_terminal"`
		NotificationToken     *string     `json:"notification_token"`
		ReceiverID            *int        `json:"receiver_id"`
		ConciliationDate      *time.Time  `json:"conciliation_date"`
		Subject               *string     `json:"subject"`
		Status                *string     `json:"status"`
		StatusDetail          *string     `json:"status_detail"`
		Body                  *string     `json:"body"`
		PictureURL            *string     `json:"picture_url"`
		ReceiptURL            *string     `json:"receipt_url"`
		NotifyURL             *string     `json:"notify_url"`
		CancelURL             *string     `json:"cancel_url"`
		ReturnURL             *string     `json:"return_url"`
		ContactURL            *string     `json:"contact_url"`
		NotifyAPIVersion      *string     `json:"notify_api_verion"`
		ExpiresDate           *time.Time  `json:"expires_date"`
		AttachmentURLS        *[]string   `json:"attachment_urls"`
		Bank                  *string     `json:"bank"`
		BankID                *string     `json:"bank_id"`
		PayerName             *string     `json:"payer_name"`
		PayerEmail            *string     `json:"payer_email"`
		PersonalIdentifier    *string     `json:"personal_identifier"`
		BankAccountNumber     *string     `json:"bank_account_number"`
		OutOfDateConciliation *bool       `json:"out_of_date_conciliation"`
		TransactionID         *string     `json:"transaction_id"`
		Custom                *string     `json:"custom"`
		ResponsibleUserEmail  *string     `json:"responsible_user_email"`
		SendReminders         *bool       `json:"send_reminders"`
		SendEmail             *bool       `json:"send_email"`
		PaymentMethod         *string     `json:"payment_method"`
	}{
		PaymentID:             &p.PaymentID,
		PaymentURL:            &p.PaymentURL,
		SimplifiedTransferURL: &p.SimplifiedTransferURL,
		TransferURL:           &p.TransferURL,
		AppURL:                &p.AppURL,
		ReceiverID:            &p.ReceiverID,
		ReadyForTerminal:      &p.ReadyForTerminal,
		NotificationToken:     &p.NotificationToken,
		ConciliationDate:      &p.ConciliationDate,
		Subject:               &p.Subject,
		Status:                &p.Status,
		StatusDetail:          &p.StatusDetail,
		Body:                  &p.Body,
		PictureURL:            &p.PictureURL,
		ReceiptURL:            &p.ReceiptURL,
		NotifyURL:             &p.NotifyURL,
		ReturnURL:             &p.ReturnURL,
		CancelURL:             &p.CancelURL,
		ContactURL:            &p.ContactURL,
		NotifyAPIVersion:      &p.NotifyAPIVersion,
		ExpiresDate:           &p.ExpiresDate,
		AttachmentURLS:        &p.AttachmentURLS,
		Bank:                  &p.Bank,
		BankID:                &p.BankID,
		PayerName:             &p.PayerName,
		PayerEmail:            &p.PayerEmail,
		PersonalIdentifier:    &p.PersonalIdentifier,
		BankAccountNumber:     &p.BankAccountNumber,
		OutOfDateConciliation: &p.OutOfDateConciliation,
		TransactionID:         &p.TransactionID,
		Custom:                &p.Custom,
		ResponsibleUserEmail:  &p.ResponsibleUserEmail,
		SendReminders:         &p.SendReminders,
		SendEmail:             &p.SendEmail,
		PaymentMethod:         &p.PaymentMethod,
	}

	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}

	p.Amount = float64(raw.Amount)

	return nil
}

// SuccessResponse represents a success response defined by khipu
type SuccessResponse struct {
	Message string `json:"message"`
}
