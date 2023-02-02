package models

type Transaction struct {
	Name         string `json:"name,omitempty"`
	Date         string `json:"date,omitempty"`
	MerchantCode string `json:"merchant_code,omitempty"`
	AmountCents  int    `json:"amount_cents,omitempty"`
	MaxRewards   int    `json:"max_rewards,omitempty"`
}

type TransactionsSummary struct {
	Merchants     map[string]int // map[merchant_code](amount of cents spent at merchant)
	MaxCentsSpent int            // maximum amount of cents spent at any one of the merchants
	// TotalCentsSpent int            // Not used, may potentially be needed, makes sense to have semantically
}

type TransactionGroup struct {
	Transactions []Transaction
	MaxRewards   int // maximum amount of rewards using ALL of the transactions
}

type TransactionGroups struct {
	Grouping string // date grouping (day, month, year)
	Groups   map[string]*TransactionGroup
}

type RewardRule struct {
	Name      string         `json:"name,omitempty"`
	Merchants map[string]int `json:"merchants,omitempty"` // map[merchant_code](amount of required to spend at merchant)
	Points    int            `json:"points,omitempty"`    // points earned from reward
}

type CalculateRewardsRequest struct {
	Transactions map[string]Transaction `json:"transactions,omitempty"`
	Rules        []RewardRule           `json:"rules,omitempty"`
	Grouping     string                 `json:"grouping,omitempty"`
}

type CalculateRewardsResponse struct {
	Grouping string
	Groups   map[string]TransactionGroup
}
