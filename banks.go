package ghipu

import "net/http"

// Bank represents the schema of bank defined by khipu
type Bank struct {
	BankID    string `json:"bank_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Message   string `json:"message,omitempty"`
	MinAmount string `json:"min_amount,omitempty"`
	Type      string `json:"type,omitempty"`
	Parent    string `json:"parent,omitempty"`
}

type BankService struct {
	client *httpClient
}

func NewBankService(secret, receiverID string) *BankService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &BankService{&client}
}

func (bs *BankService) Banks() ([]*Bank, error) {
	path := basePath + "/banks"
	resp, err := bs.client.Get(path, nil)
	if err != nil {
		return nil, err
	}

	var payload struct {
		Banks []*Bank `json:"banks"`
	}
	if err := unmarshalJSON(resp.Body, &payload); err != nil {
		return nil, err
	}
	return payload.Banks, nil
}
