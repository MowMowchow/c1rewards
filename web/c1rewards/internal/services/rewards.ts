import {
  MaximumRewardsResponseBody,
  Transaction,
  TransactionGroup,
} from "../models";

export const fetchMaximumRewards = async (
  maximumRewardsRequestBody: any
): Promise<MaximumRewardsResponseBody> => {
  try {
    const requestOptions = {
      method: "POST",
      body: JSON.stringify(maximumRewardsRequestBody),
    };
    const response = await fetch(
      "https://a2oplo4548.execute-api.us-east-1.amazonaws.com/getRewards",
      requestOptions
    );
    const responseBody: MaximumRewardsResponseBody =
      (await response.json()) as MaximumRewardsResponseBody;

    let typedResponseBody: MaximumRewardsResponseBody = {
      Grouping: responseBody.Grouping,
      Groups: new Map<string, TransactionGroup>(
        Object.entries(responseBody.Groups)
      ),
    };

    typedResponseBody.Groups.forEach((transactionGroup: any, date) => {
      let tempTransactionGroup: TransactionGroup = {
        MaxRewards: transactionGroup.MaxRewards,
        Transactions: transactionGroup.Transactions,
      };
      typedResponseBody.Groups.set(date, tempTransactionGroup);
    });

    return typedResponseBody;
  } catch (error) {
    console.error("error fetching maximum rewards", error);
    return {} as Promise<MaximumRewardsResponseBody>;
  }
};
