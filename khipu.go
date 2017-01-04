package gokhipu

import (
	"fmt"
	"net/http"
)

var (
	basePath = "https://khipu.com/api/2.0"
)

// Khipu implements a basic struct with revelant data to handle
// khipu's API requests.
type Khipu struct {
	Secret     string
	ReceiverID string
}

// NewKhipuClient returns an instance of khipu that is the client to make payment request
func NewKhipuClient(secret, receiverID string) *Khipu {
	return &Khipu{
		Secret:     secret,
		ReceiverID: receiverID,
	}
}

// Banks ...
func (kc *Khipu) Banks() (*http.Response, error) {
	requestPath := basePath + "/banks"
	req, err := http.NewRequest("GET", requestPath, nil)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to create request to %s\n%s", requestPath, err))
	}

	req.Header.Set("Authorization", setAuth(nil, "GET", requestPath, kc.Secret, kc.ReceiverID))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	cl := http.Client{}
	res, err := cl.Do(req)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("failed to made request to %s\n%s", requestPath, err))
	}

	return res, nil
}
