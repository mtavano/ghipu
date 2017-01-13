package ghipu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const basePath = "https://khipu.com/api/2.0"

// Client represents an khipu's REST API.
type Client struct {
	PaymentService
	BankService
}

// NewClient returns an instance of khipu that is the client to make payment request
func NewClient(secret, receiverID string) *Client {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}

	return &Client{
		PaymentService: PaymentService{&client},
		BankService:    BankService{&client},
	}
}

type httpClient struct {
	client *http.Client
	secret string
	recid  string
}

func (hc *httpClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	req.Header.Set("Authorization", hc.authParams(req, values))
	req.Header.Set("Accept", "application/json")

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: %s resquest failed, %v", req.URL, err)
	}

	return resp, nil
}

func (hc *httpClient) Get(path string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return hc.Do(req, values)
}

func (hc *httpClient) PostForm(path string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return hc.Do(req, values)
}

func (hc *httpClient) authParams(req *http.Request, values url.Values) string {
	var buff bytes.Buffer
	buff.WriteString(url.QueryEscape(req.Method))
	buff.WriteByte('&')
	buff.WriteString(url.QueryEscape(req.URL.RequestURI()))

	if values != nil {
		buff.WriteByte('&')
		buff.WriteString(strings.Replace(values.Encode(), "+", "%20", -1)[1:])
	}

	sign := hmac.New(sha256.New, []byte(hc.secret))
	sign.Write(buff.Bytes())
	return hc.recid + ":" + hex.EncodeToString(sign.Sum(nil))
}
