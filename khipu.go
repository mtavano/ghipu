package ghipu

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var baseURL = url.URL{
	Scheme: "https",
	Host:   "khipu.com",
	Path:   "api/2.0",
}

// Client represents an khipu's REST API.
type Client struct {
	PaymentsService
	BankService
	ReceiversService
}

// NewClient returns an instance of khipu that is the client to make payment request
func NewClient(secret string, receiverID int) *Client {
	hclient := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}

	return &Client{
		PaymentsService:  PaymentsService{&hclient},
		BankService:      BankService{&hclient},
		ReceiversService: ReceiversService{&hclient},
	}
}

// AuthorizationError represents an authorization error of the Kiphu's REST API.
type AuthorizationError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (err *AuthorizationError) Error() string {
	return fmt.Sprintf("ghipu: unauthorized request, %v", err.Message)
}

// ServiceError represents an service error of the Kiphu's REST API.
type ServiceError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func (err *ServiceError) Error() string {
	return fmt.Sprintf("ghipu: unauthorized request, %v", err.Message)
}

// ErrorItem represents a validation error item.
type ErrorItem struct{ Field, Message string }

// ValidationError represents an validation error of the Kiphu's REST API.
type ValidationError struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  []ErrorItem `json:"errors"`
}

func (err *ValidationError) Error() string {
	var buff bytes.Buffer
	buff.WriteString("ghipu: invalid request")

	for _, e := range err.Errors {
		buff.WriteString(", ")
		buff.WriteString(e.Field)
		buff.WriteByte(':')
		buff.WriteString(e.Message)
	}

	return buff.String()
}

type httpClient struct {
	client *http.Client
	secret string
	recid  int
}

func (hc *httpClient) Do(req *http.Request, values url.Values) (*http.Response, error) {
	hc.signRequest(req, values)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := hc.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: %s resquest failed, %v", req.URL, err)
	}

	switch resp.StatusCode {
	case http.StatusBadRequest:
		var valErr ValidationError
		if err = unmarshalJSON(resp.Body, &valErr); err != nil {
			return nil, fmt.Errorf("ghipu: error parsing response, %v", err)
		}
		return nil, &valErr
	case http.StatusForbidden:
		var authErr AuthorizationError
		if err = unmarshalJSON(resp.Body, &authErr); err != nil {
			return nil, fmt.Errorf("ghipu: error parsing response, %v", err)
		}
		return nil, &authErr
	case http.StatusServiceUnavailable:
		var svcErr ServiceError
		if err = unmarshalJSON(resp.Body, &svcErr); err != nil {
			return nil, fmt.Errorf("ghipu: error parsing response, %v", err)
		}
		return nil, &svcErr
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

func (hc *httpClient) Delete(path string, values url.Values) (*http.Response, error) {
	uri := baseURL.String() + path
	req, err := http.NewRequest(http.MethodDelete, uri, nil)
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

var percentEncode = strings.NewReplacer(
	"+", "%20",
	"*", "%2A",
	"%7A", "~",
)

func (hc *httpClient) signRequest(req *http.Request, values url.Values) {
	var buff bytes.Buffer
	buff.WriteString(url.QueryEscape(req.Method))
	buff.WriteByte('&')
	buff.WriteString(url.QueryEscape(req.URL.Scheme + "://" + req.URL.Host + req.URL.Path))

	if values != nil {
		buff.WriteByte('&')
		buff.WriteString(percentEncode.Replace(values.Encode()))
	}

	sig := hmac.New(sha256.New, []byte(hc.secret))
	sig.Write(buff.Bytes())

	sign := hex.EncodeToString(sig.Sum(nil))

	buff.Reset()
	buff.WriteString(strconv.Itoa(hc.recid))
	buff.WriteByte(':')
	buff.WriteString(sign)

	req.Header.Set("Authorization", buff.String())
}

func unmarshalJSON(r io.ReadCloser, v interface{}) error {
	defer r.Close()

	body, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, v)
}
