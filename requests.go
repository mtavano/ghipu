package gokhipu

import (
	"fmt"
	"net/http"

	"github.com/mtavano/go-khipu/utils"
)

// GetBanks returns the list of banks availables to make a payment through khipu's API
func (kc *KhipuClient) GetBanks() ([]*Bank, error) {
	requestPath := basePath + "/banks"
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(nil, "GET", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := kc.DoRequest(req, requestPath)
	if err != nil {
		return nil, err
	}

	payload := new(struct {
		Banks []*Bank `json:"banks"`
	})

	err = utils.UnmarshalJSON(res.Body, payload)
	if err != nil {
		return nil, err
	}

	return payload.Banks, nil
}

// GetPayment ...
func (kc *KhipuClient) GetPayment(id string) (*PaymentResponse, error) {
	requestPath := basePath + "/payments/" + id
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(nil, "GET", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := kc.DoRequest(req, requestPath)
	if err != nil {
		return nil, err
	}

	pr := new(PaymentResponse)
	err = utils.UnmarshalJSON(res.Body, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// MakePay makes a pay through a POST calls to khipu's API and return the PaymentResponse
// Given by Khipu's API
func (kc *KhipuClient) MakePay(p *Payment) (*PaymentResponse, error) {
	requestPath := basePath + "/payments"

	req, err := http.NewRequest("POST", requestPath, utils.MakeForm(p.Params()))
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(p.Params(), "POST", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	res, err := kc.DoRequest(req, requestPath)
	if err != nil {
		return nil, err
	}

	pr := new(PaymentResponse)
	err = utils.UnmarshalJSON(res.Body, pr)
	if err != nil {
		return nil, err
	}

	return pr, nil
}

// DoRequest makes a request to khipu's API using theyr schema
func (kc *KhipuClient) DoRequest(req *http.Request, path string) (*http.Response, error) {
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
