import { TransactionGroup } from "../../internal/models";
import TGTransactionRow from "./TGTransaction";

export type TransactionGroupRowProps = {
  date: string;
  transactiongroup: TransactionGroup;
};

export default function TransactionGroupRow({
  date,
  transactiongroup,
}: TransactionGroupRowProps) {
  return (
    <li className="my-2 border-2 rounded-md shadow-sm">
      <div className="m-2 flex flex-row justify-start">
        <div className="w-3/12">
          <h1 className="text=lg">{`Date (Month): ${date}`}</h1>
          <h1 className="text=lg">{`Max Rewards (for Month): ${transactiongroup.MaxRewards}`}</h1>
        </div>
        <div className="w-9/12">
          <h1 className="text-lg">Transactions:</h1>
          <ul>
            {transactiongroup.Transactions.map((transaction) => (
              <TGTransactionRow
                name={transaction.name}
                date={transaction.date}
                merchant_code={transaction.merchant_code}
                amount_cents={transaction.amount_cents}
                max_rewards={transaction.max_rewards}
              />
            ))}
          </ul>
        </div>
      </div>
    </li>
  );
}
