import { Rule } from "../../internal/models";

export type RuleRowProps = {
  name?: string;
  points: number;
  merchants: Map<string, number>;
  deleteRule?: (ruleName: string) => void;
};

export default function RuleRow({
  name,
  points,
  merchants,
  deleteRule,
}: RuleRowProps) {
  return (
    <li className="my-2 border-2 rounded-md shadow-sm">
      <div className="m-2 flex flex-row justify-around">
        <div className="mx-2 w-2/12 flex flex-col justify-center items-center">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-2 w-2/12 flex flex-col justify-center items-center">
          <h1>{`Points: ${points}`}</h1>
        </div>
        <div className="mx-2 w-4/12 flex flex-col justify-center items-center">
          <ul>
            {Array.from(merchants.entries()).map(([name, cost]) => (
              <li className="grid grid-cols-2 space-x-8">
                <p>{`Name: ${name}`}</p>
                <p>{`Cost: ${cost}`}</p>
              </li>
            ))}
          </ul>
        </div>
        <div className="mx-2 w-2/12 flex flex-col justify-center items-center">
          {deleteRule === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="py-1 px-2 shadow-md rounded-md bg-rose-400"
              onClick={() => deleteRule(name)}
            >
              <h1 className="text-md text-stone-50">delete</h1>
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
