import { Transaction } from "../../internal/models";

export type TransactionRowProps = {
  date: string;
  merchant_code: string;
  amount_cents: number;
  name?: string;
  max_rewards?: number;
  deleteTransaction?: (transactionName: string) => void;
};

export default function TransactionRow({
  name,
  date,
  merchant_code,
  amount_cents,
  max_rewards = -1,
  deleteTransaction,
}: TransactionRowProps) {
  console.log(
    "TRANSACTION ROW",
    name,
    date,
    merchant_code,
    amount_cents,
    max_rewards
  );
  return (
    <li className="my-2 border-2 rounded-md w-fit shadow-sm">
      <div className="m-2 flex flex-row justify-start">
        <div className="mx-4">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Date: ${date}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Merchant Code: ${merchant_code}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Amount (cents): ${amount_cents}`}</h1>
        </div>
        <div className="mx-4">
          {max_rewards < 0 ? (
            max_rewards == -2 ? (
              <></>
            ) : (
              <h1>No Reward Available</h1>
            )
          ) : (
            <h1>{`Max Reward: ${max_rewards}`}</h1>
          )}
        </div>
        <div className="mx-4">
          {deleteTransaction === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="py-1 px-2 rounded-md bg-rose-400"
              onClick={() => deleteTransaction(name)}
            >
              delete
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
