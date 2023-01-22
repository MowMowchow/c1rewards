package models

type TransactionsSummary struct {
	Merchants map[string]int // Total amount spent at each merchant
	MaxCents  int
}

type Transaction struct {
	Date         string `json:"date,omitempty"`
	MerchantCode string `json:"merchant_code,omitempty"`
	AmountCents  int    `json:"amount_cents"`
}

type CalculateTransactionRequest struct {
	Transactions map[string]Transaction `json:"transactions"`
}
