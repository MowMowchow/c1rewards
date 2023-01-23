import { TransactionGroup } from "../../internal/models";
import TransactionRow from "./TransactionRow";

export type TransactionGroupRowProps = {
  date: string;
  transactiongroup: TransactionGroup;
};

export default function TransactionGroupRow({
  date,
  transactiongroup,
}: TransactionGroupRowProps) {
  console.log("TRANSACTION GROUP", transactiongroup);
  return (
    <li className="my-2 border-2 rounded-md w-fit shadow-sm">
      <div className="m-2 flex flex-row justify-start">
        <div className="mx-4">
          <h1>{`Max Rewards for month: ${date} = ${transactiongroup.MaxRewards}`}</h1>
        </div>
        <div>
          {transactiongroup.Transactions.map((transaction) => (
            <TransactionRow
              name={transaction.name}
              date={transaction.date}
              merchant_code={transaction.merchant_code}
              amount_cents={transaction.amount_cents}
              max_rewards={transaction.max_rewards}
            />
          ))}
        </div>
      </div>
    </li>
  );
}
