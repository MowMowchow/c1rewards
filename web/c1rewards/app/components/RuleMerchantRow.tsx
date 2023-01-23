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
    <li className="my-2 border-2 rounded-md shadow-sm">
      <div className="m-2 grid grid-cols-3 justify-start">
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-4 flex flex-col justify-center items-center">
          <h1>{`Cost: ${cost}`}</h1>
        </div>
        <div className="mx-4 flex flex-col justify-center items-center">
          {deleteRuleMerchantRow === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="y-1 px-2 shadow-md rounded-md bg-rose-400"
              onClick={() => deleteRuleMerchantRow(name)}
            >
              <h1 className="text-md text-stone-50">delete</h1>
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
