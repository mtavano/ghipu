package gokhipu

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
