export type Transaction = {
  date: string;
  merchant_code: string;
  amount_cents: number;
  name?: string;
  max_rewards?: number;
};

export type Rule = {
  name?: string;
  merchants: Map<string, number>;
  points: number;
};

export type RuleMerchant = {
  name: string;
  cost: number;
};

export type MaximumRewardsRequestBody = {
  transactions: Object;
  rules: Rule[];
};

export type TransactionGroup = {
  Transactions: Transaction[];
  MaxRewards: number;
};

export type MaximumRewardsResponseBody = {
  Grouping: string;
  Groups: Map<string, TransactionGroup>;
};

// export type MaximumRewardsRequestBody2 = {
//   Grouping: string;
//   Groups: Object<string, TransactionGroup>;

// }
