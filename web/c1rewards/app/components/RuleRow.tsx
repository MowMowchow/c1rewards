import { Rule } from "../../internal/models";

export type RuleRowProps = {
  name?: string;
  points: number;
  merchants: string[];
  deleteRule?: (ruleName: string) => void;
};

export default function RuleRow({
  name,
  points,
  merchants,
  deleteRule,
}: RuleRowProps) {
  return (
    <li className="my-2 border-2 rounded-md w-fit shadow-sm">
      <div className="m-2 flex flex-row justify-start">
        <div className="mx-4">
          <h1>{`Name: ${name}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Points: ${points}`}</h1>
        </div>
        <div className="mx-4">
          <h1>{`Merchants: ${merchants}`}</h1>
        </div>
        <div className="mx-4">
          {deleteRule === undefined || name === undefined ? (
            <></>
          ) : (
            <button
              className="py-1 px-2 rounded-md bg-rose-400"
              onClick={() => deleteRule(name)}
            >
              delete
            </button>
          )}
        </div>
      </div>
    </li>
  );
}
