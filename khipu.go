package gokhipu

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var (
	basePath = "https://khipu.com/api/2.0"
)

// Khipu implements a basic struct with revelant data to handle
// khipu's API requests.
type Khipu struct {
	Secret     string
	ReceiverID string
	httpClient *http.Client
}

// NewKhipuClient returns an instance of khipu that is the client to make payment request
func NewKhipuClient(secret, receiverID string) *Khipu {
	return &Khipu{
		Secret:     secret,
		ReceiverID: receiverID,
		httpClient: http.DefaultClient,
	}
}

// Banks returns the list of banks availables to make a payment through khipu's API
func (kc *Khipu) Banks() ([]*Bank, error) {
	requestPath := basePath + "/banks"
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(nil, "GET", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := kc.doRequest(req, requestPath)
	if err != nil {
		return nil, err
	}

	payload := new(struct {
		Banks []*Bank `json:"banks"`
	})

	err = unmarshalJSON(res.Body, payload)
	if err != nil {
		return nil, err
	}

	return payload.Banks, nil
}

// Pay makes a pay through a POST calls to khipu's API
func (kc *Khipu) Pay(p *Payment) (*PaymentResponse, error) {
	requestPath := basePath + "/payments"

	req, err := http.NewRequest("POST", requestPath, makeForm(p.Params()))
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(p.Params(), "POST", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := kc.doRequest(req, requestPath)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	pr := new(PaymentResponse)
	err = unmarshalJSON(res.Body, pr)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%#v", pr)

	return &PaymentResponse{}, nil
}

// Private functions section

func (kc *Khipu) doRequest(req *http.Request, path string) (*http.Response, error) {
	res, err := kc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("error, failde to made request to %s\n%#v", path, err))
	}
	var status int
	status = res.StatusCode / 10
	if status != 20 {
		return nil, fmt.Errorf("error, unexpected status code: %d", res.StatusCode)
	}

	return res, nil
}

func makeForm(m map[string]string) *strings.Reader {
	vals := url.Values{}
	for k, v := range m {
		if v == "" {
			continue
		}
		vals.Set(k, v)
	}

	return strings.NewReader(vals.Encode())
}
