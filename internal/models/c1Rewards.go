package models

type Transaction struct {
	Date         string `json:"date,omitempty"`
	MerchantCode string `json:"merchant_code,omitempty"`
	AmountCents  int    `json:"amount_cents,omitempty"`
}

type TransactionsSummary struct {
	Merchants map[string]int // Total amount spent at each merchant
	MaxCents  int
}

type RewardRule struct {
	Merchants map[string]int `json:"merchants,omitempty"` // [merchantName]cost
	Points    int            `json:"points,omitempty"`    // points earned from deal
}

type Rules struct {
	Rules  []RewardRule `json:"rules,omitempty"`
	Length int
}

type CalculateTransactionRequest struct {
	Transactions map[string]Transaction `json:"transactions,omitempty"`
	Rules        []RewardRule           `json:"rules,omitempty"`
}
