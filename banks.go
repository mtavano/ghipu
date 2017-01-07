package ghipu

// Bank represents the schema of bank defined by khipu
type Bank struct {
	BankID    string `json:"bank_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Message   string `json:"message,omitempty"`
	MinAmount string `json:"min_amount,omitempty"`
	Type      string `json:"type,omitempty"`
	Parent    string `json:"parent,omitempty"`
}
