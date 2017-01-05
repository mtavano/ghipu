package gokhipu

// Bank represents the schema of bank defined by khipu
type Bank struct {
	BankID    string `json:"bank_id"`
	Name      string `json:"name"`
	Message   string `json:"message"`
	MinAmount string `json:"min_amount"`
	Type      string `json:"type"`
	Parent    string `json:"parent"`
}
