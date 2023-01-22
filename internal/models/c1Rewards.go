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

type RewardMerchant struct {
	Cost   int `json:"cost,omitempty"`
	Points int `json:"points,omitempty"`
}

type RewardRule struct {
	Merchants map[string]RewardMerchant `json:"merchants,omitempty"`
}

type Rules struct {
	Rules  []RewardRule `json:"rules,omitempty"`
	Length int
}

type CalculateTransactionRequest struct {
	Transactions map[string]Transaction `json:"transactions,omitempty"`
	Rules        []RewardRule           `json:"rules,omitempty"`
}
