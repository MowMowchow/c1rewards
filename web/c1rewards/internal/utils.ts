import { stringify } from "querystring";
import { MaximumRewardsRequestBody, Rule, Transaction } from "./models";

export const formatMaximumRewardsRequest = (
  transactions: Transaction[],
  rules: Rule[]
): any => {
  const transactionMapping = new Map<string, Transaction>();
  transactions.forEach((transaction) => {
    if (transaction?.name) {
      transactionMapping.set(transaction.name, transaction);
    }
  });

  const ruleMapping: any[] = [];
  rules.forEach((rule) => {
    const tempRule = {
      name: rule.name,
      merchants: Object.fromEntries(rule.merchants),
      points: Number(rule.points),
    };
    ruleMapping.push(tempRule);
  });

  const maximumRewardsRequestBody = {
    transactions: Object.fromEntries(transactionMapping),
    rules: ruleMapping,
    grouping: "month",
  };

  return maximumRewardsRequestBody;
};
