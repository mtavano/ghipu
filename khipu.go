package ghipu

import "net/http"

// KhipuClient implements a basic struct with revelant data to handle
// khipu's API requests.
type KhipuClient struct {
	Secret     string
	ReceiverID string
	httpClient *http.Client
}

// NewKhipuClient returns an instance of khipu that is the client to make payment request
func NewKhipuClient(secret, receiverID string) *KhipuClient {
	return &KhipuClient{
		Secret:     secret,
		ReceiverID: receiverID,
		httpClient: http.DefaultClient,
	}
}
