package ghipu

import "net/http"

// BankItem represents the schema of bank defined by khipu
type BankItem struct {
	BankID    string  `json:"bank_id,omitempty"`
	Name      string  `json:"name,omitempty"`
	Message   string  `json:"message,omitempty"`
	MinAmount float64 `json:"min_amount,omitempty"`
	Type      string  `json:"type,omitempty"`
	Parent    string  `json:"parent,omitempty"`
}

// BanksResponse represents schema of banks response.
type BanksResponse struct {
	Banks []BankItem
}

type BankService struct {
	client *httpClient
}

func NewBankService(secret string, receiverID int) *BankService {
	client := httpClient{
		client: &http.Client{},
		secret: secret,
		recid:  receiverID,
	}
	return &BankService{&client}
}

func (bs *BankService) Banks() ([]BankItem, error) {
	resp, err := bs.client.Get("/banks", nil)
	if err != nil {
		return nil, err
	}

	var res BanksResponse
	if err := unmarshalJSON(resp.Body, &res); err != nil {
		return nil, err
	}
	return res.Banks, nil
}
