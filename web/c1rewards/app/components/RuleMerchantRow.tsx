import { RuleMerchant } from "../../internal/models";

export type RuleMerchantRowProps = {
  name?: string;
  cost: number;
  deleteRuleMerchantRow?: (ruleMerchantName: string) => void;
};

export default function RuleMerchantRow({
  name,
  cost,
  deleteRuleMerchantRow,
}: RuleMerchantRowProps) {
  return (
    <li className="my-2 border-2 rounded-md w-fit shadow-sm">
      <div className="m-2 flex flex-row justify-start">
        <div className="mx-4">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Cost: ${cost}`}</h1>
        </div>
        <div className="mx-4">
          {deleteRuleMerchantRow === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="py-1 px-2 rounded-md bg-rose-400"
              onClick={() => deleteRuleMerchantRow(name)}
            >
              delete
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
