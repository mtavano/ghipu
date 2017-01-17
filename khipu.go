package ghipu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "khipu.com",
	Path:   "api/2.0",
}

// Client represents an khipu's REST API.
type Client struct {
	PaymentService
	BankService
}

// NewClient returns an instance of khipu that is the client to make payment request
func NewClient(secret, receiverID string) *Client {
	hclient := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}

	return &Client{
		PaymentService: PaymentService{&hclient},
		BankService:    BankService{&hclient},
	}
}

// ErrorResponse represents an error of the Kiphu's REST API.
type ErrorResponse string

func (err ErrorResponse) Error() string {
	return "ghipu: " + string(err)
}

// Khipu's REST API errors
const (
	ErrInvalidPayload       ErrorResponse = "invalid payload"
	ErrInvalidAuthorization ErrorResponse = "invalid request"
	ErrInvalidOperation     ErrorResponse = "invalid operation"
)

type httpClient struct {
	client *http.Client
	secret string
	recid  string
}

func (hc *httpClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	hc.signRequest(req, values)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	reqdump, _ := httputil.DumpRequest(req, true)
	fmt.Printf("%s\n", reqdump)

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: %s resquest failed, %v", req.URL, err)
	}

	respdump, _ := httputil.DumpResponse(resp, false)
	fmt.Printf("%s\n", respdump)

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, ErrInvalidPayload
	case http.StatusForbidden:
		return nil, ErrInvalidAuthorization
	case http.StatusServiceUnavailable:
		return nil, ErrInvalidOperation
	default:
		return resp, nil
	}
}

func (hc *httpClient) Get(path string, values url.Values) (*http.Response, error) {
	uri := baseURL.String() + path
	req, err := http.NewRequest(http.MethodGet, uri, nil)
	if err != nil {
		return nil, err
	}
	return hc.Do(req, values)
}

func (hc *httpClient) PostForm(path string, values url.Values) (*http.Response, error) {
	uri := baseURL.String() + path
	req, err := http.NewRequest(http.MethodPost, uri, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	return hc.Do(req, values)
}

func (hc *httpClient) signRequest(req *http.Request, values url.Values) {
	var buff bytes.Buffer
	buff.WriteString(url.QueryEscape(req.Method))
	buff.WriteByte('&')
	buff.WriteString(url.QueryEscape(req.URL.RequestURI()))

	if values != nil {
		buff.WriteByte('&')
		buff.WriteString(strings.Replace(values.Encode(), "+", "%20", -1)[1:])
	}

	sig := hmac.New(sha256.New, []byte(hc.secret))
	sig.Write(buff.Bytes())

	sign := hex.EncodeToString(sig.Sum(nil))
	req.Header.Set("Authorization", hc.recid+":"+sign)
}
