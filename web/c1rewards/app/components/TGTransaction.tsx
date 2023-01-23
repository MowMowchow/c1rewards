export type TGTransactionRowProps = {
  date: string;
  merchant_code: string;
  amount_cents: number;
  name?: string;
  max_rewards?: number;
  deleteTGTransaction?: (tGTransactionName: string) => void;
};

export default function TGTransactionRow({
  name,
  date,
  merchant_code,
  amount_cents,
  deleteTGTransaction,
  max_rewards,
}: TGTransactionRowProps) {
  return (
    <li className="my-2 border-2 rounded-md shadow-sm">
      <div className="grid grid-cols-5 m-2">
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Date: ${date}`}</h1>
        </div>
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Merchant Code: ${merchant_code}`}</h1>
        </div>
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Amount (cents): ${amount_cents}`}</h1>
        </div>
        {max_rewards ? (
          <h1>{`Max Rewards: ${max_rewards}`}</h1>
        ) : (
          <h1>No Reward Available</h1>
        )}

        <div className="mx-4 flex flex-col justify-center items-center">
          {deleteTGTransaction === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="py-1 px-2 shadow-md rounded-md bg-rose-400"
              onClick={() => deleteTGTransaction(name)}
            >
              <h1 className="text-md text-stone-50"> Delete</h1>
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
