"use client";
import Image from "next/image";
import { Inter } from "@next/font/google";
import styles from "./page.module.css";
import { useState } from "react";
import {
  MaximumRewardsResponseBody,
  Rule,
  RuleMerchant,
  Transaction,
  TransactionGroup,
} from "../internal/models";
import TransactionRow from "./components/TransactionRow";
import RuleRow from "./components/RuleRow";
import RuleMerchantRow from "./components/RuleMerchantRow";
import { fetchMaximumRewards } from "@/internal/services/rewards";
import { formatMaximumRewardsRequest } from "@/internal/utils";
import defaultValues from "./defaultValues.json";
import { stringify } from "querystring";
import TransactionGroupRow from "./components/TransactionGroup";

export default function Home() {
  // const [transactionList, setTransactionList] = useState<Transaction[]>([]);
  const [transactionList, setTransactionList] = useState<Transaction[]>(
    defaultValues.transactions
  );
  const [newTransaction, setNewTransaction] = useState<Transaction>(
    {} as Transaction
  );

  // const [ruleList, setRuleList] = useState<Rule[]>([]);
  const defaultRules: Rule[] = [];
  defaultValues.rules.forEach((rule) => {
    const tempRule: Rule = {
      name: rule.name,
      merchants: new Map(
        Object.entries(JSON.parse(JSON.stringify(rule.merchants)))
      ),
      points: rule.points,
    };
    rule.merchants;
    defaultRules.push(tempRule);
  });
  // const [ruleList, setRuleList] = useState<Rule[]>();
  const [ruleList, setRuleList] = useState<Rule[]>(defaultRules);
  const [newRule, setNewRule] = useState<Rule>({
    merchants: new Map<string, number>(),
  } as Rule);

  const [newRuleMerchant, setNewRuleMerchant] = useState<RuleMerchant>(
    {} as RuleMerchant
  );

  const [transactionGroups, setTransactionGroups] =
    useState<Map<string, TransactionGroup>>();

  // const handleSubmitTransaction = (event) => {
  //   setTransactionList((prevState) => [...prevState, event.target]);
  // };
  const handleNewTransactionChange = (field: string, value: any) => {
    setNewTransaction((prevState) => ({
      ...prevState,
      name: field === "name" ? value : prevState.name,
      date: field === "date" ? value : prevState.date,
      merchant_code:
        field === "merchant_code" ? value : prevState.merchant_code,
      amount_cents:
        field === "amount_cents" ? Number(value) : prevState.amount_cents,
    }));
  };

  const handleSubmitNewTransaction = () => {
    if (
      transactionList.some(
        (transaction) => transaction.name === newTransaction.name
      )
    ) {
      return;
    }
    setTransactionList((prevState) => [...prevState, newTransaction]);
  };

  const handleNewRuleChange = (field: string, value: any) => {
    setNewRule((prevState) => ({
      ...prevState,
      name: field === "name" ? value : prevState.name,
      points: field === "points" ? Number(value) : prevState.points,
      merchants: field === "merchants" ? value : prevState.merchants,
    }));
  };

  const handleNewRuleMerchantChange = (field: string, value: any) => {
    setNewRuleMerchant((prevState) => ({
      ...prevState,
      name: field === "name" ? value : prevState.name,
      cost: field === "cost" ? Number(value) : prevState.cost,
    }));
  };

  const handleSubmitNewRule = () => {
    if (ruleList.some((rule) => rule.name === newRule.name)) {
      return;
    }
    setRuleList((prevState) => [...prevState, newRule]);
    setNewRule({
      merchants: new Map<string, number>(),
    } as Rule);
  };

  const handleSubmitNewRuleMerchant = () => {
    setNewRule((prevState) => ({
      ...prevState,
      merchants: new Map(newRule.merchants).set(
        newRuleMerchant.name,
        Number(newRuleMerchant.cost)
      ),
    }));
  };

  const handleDeleteTransaction = (transactionName: string) => {
    setTransactionList((prevState) =>
      transactionList.filter(
        (transaction) => transaction.name !== transactionName
      )
    );
  };

  const handleDeleteRule = (ruleName: string) => {
    setRuleList((prevState) =>
      ruleList.filter((rule) => rule.name !== ruleName)
    );
  };

  const handleDeleteRuleMerchantRow = (merchantName: string) => {
    const tempNewRule = new Map(newRule.merchants);
    tempNewRule.delete(merchantName);
    setNewRule((prevState) => ({
      ...prevState,
      merchants: tempNewRule,
    }));
  };

  const calculateMaximumRewards = async () => {
    const maximumRewardsRequestBody = formatMaximumRewardsRequest(
      transactionList,
      ruleList
    );
    const responseBody: MaximumRewardsResponseBody = await fetchMaximumRewards(
      maximumRewardsRequestBody
    );
    setTransactionGroups(responseBody.Groups);
  };

  return (
    <div className="grid gird-row-1 justify-center w-11/12 md:w-11/12 lg:w-9/12 ">
      <div className="grid gird-cols-1">
        <div className="px-5 font-poppins">
          <div className="grid grid-cols-1">
            <div className="m-0 sm:m-4 md:m-10" />
            <div className="flex flex-row justify-center">
              <div className="flex flex-col justify-around">
                <h1 className="text-3xl">C1 Rewards Calculator</h1>
                <h3 className="text-xl">
                  by Jason Hou (jasonhou0299@gmail.com)
                </h3>
              </div>
            </div>
            {/* Transactions */}
            <div className="m-0 sm:m-4 md:m-10" />
            <div className="">
              <h1 className="text-2xl">Add New Transaction</h1>
              <div className="mt-4 flex flex-row justify-center">
                <form
                  onSubmit={(event) => {
                    event.preventDefault();
                    handleSubmitNewTransaction();
                  }}
                  className="w-full"
                >
                  <div className="flex flex-row justify-around">
                    <input
                      title="name"
                      type="text"
                      value={newTransaction?.name}
                      onChange={(event) => {
                        handleNewTransactionChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="name"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <input
                      title="date"
                      type="text"
                      value={newTransaction?.date}
                      onChange={(event) => {
                        handleNewTransactionChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="date"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <input
                      title="merchant_code"
                      type="text"
                      value={newTransaction?.merchant_code}
                      onChange={(event) => {
                        handleNewTransactionChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="merchant_code"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <input
                      title="amount_cents"
                      type="text"
                      value={newTransaction?.amount_cents}
                      onChange={(event) => {
                        handleNewTransactionChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="amount_cents"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <button
                      className="py-1 px-2 rounded-md bg-emerald-400"
                      type="submit"
                    >
                      Add Transaction
                    </button>
                  </div>
                </form>
              </div>
              <div className="m-6">
                <ul>
                  {transactionList.map((transaction) => (
                    <TransactionRow
                      key={transaction.name}
                      name={transaction.name}
                      date={transaction.date}
                      merchant_code={transaction.merchant_code}
                      amount_cents={transaction.amount_cents}
                      max_rewards={-2}
                      deleteTransaction={handleDeleteTransaction}
                    />
                  ))}
                </ul>
              </div>
            </div>

            {/* Rules */}
            <div className="flex flex-col justify-center">
              <h1 className="text-2xl">Add New Rules</h1>
              <p>
                *NOTE: add a merchant name and cost for the rule, and then add
                the rule
              </p>
              <div className="mt-4 flex flex-col justify-center">
                <form
                  onSubmit={(event) => {
                    event.preventDefault();
                    handleSubmitNewRule();
                  }}
                  className="w-full"
                >
                  <div className="flex flex-row justify-around">
                    <input
                      title="name"
                      type="text"
                      value={newRule?.name}
                      onChange={(event) => {
                        handleNewRuleChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="name"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <input
                      title="points"
                      type="text"
                      value={newRule?.points}
                      onChange={(event) => {
                        handleNewRuleChange(
                          event.target.title,
                          event.target.value
                        );
                      }}
                      placeholder="points"
                      className="border-2 border-stone-400 rounded-md"
                    />
                    <div className="flex flex-row justify-around">
                      <button
                        className="py-1 px-2 rounded-md bg-emerald-400"
                        type="submit"
                      >
                        Add Rule
                      </button>
                    </div>
                  </div>
                </form>

                {/* New Rule Merchants */}
                <div className="mt-4 flex flex-col justify-center">
                  <form
                    onSubmit={(event) => {
                      event.preventDefault();
                      handleSubmitNewRuleMerchant();
                    }}
                    className="w-full"
                  >
                    <div className="flex flex-row justify-around">
                      <input
                        title="name"
                        type="text"
                        value={newRuleMerchant.name}
                        onChange={(event) => {
                          handleNewRuleMerchantChange(
                            event.target.title,
                            event.target.value
                          );
                        }}
                        placeholder="merchant name"
                        className="border-2 border-stone-400 rounded-md"
                      />
                      <input
                        title="cost"
                        type="text"
                        value={newRuleMerchant.cost}
                        onChange={(event) => {
                          handleNewRuleMerchantChange(
                            event.target.title,
                            event.target.value
                          );
                        }}
                        placeholder="merchant cost"
                        className="border-2 border-stone-400 rounded-md"
                      />
                      <div className="flex flex-row justify-around">
                        <button
                          className="py-1 px-2 rounded-md bg-emerald-400"
                          type="submit"
                        >
                          Add Merchant
                        </button>
                      </div>
                    </div>
                  </form>
                  <div className="m-6">
                    <h1>Merchants for current rule:</h1>
                    <ul>
                      {Array.from(newRule.merchants.entries()).map(
                        ([merchantName, merchantCost]) => (
                          <RuleMerchantRow
                            name={merchantName}
                            cost={merchantCost}
                            deleteRuleMerchantRow={handleDeleteRuleMerchantRow}
                          />
                        )
                      )}
                    </ul>
                  </div>
                </div>
                <div className="mt-4">
                  <h1 className="text-2xl">Rules:</h1>
                  <ul>
                    {ruleList.map((rule) => (
                      <RuleRow
                        key={rule.name}
                        name={rule.name}
                        points={rule.points}
                        merchants={Array.from(rule.merchants.entries()).map(
                          ([name, cost]) => ` name: ${name} and cost:${cost}`
                        )}
                        deleteRule={handleDeleteRule}
                      />
                    ))}
                  </ul>
                </div>
              </div>
            </div>
            {/* Maximum Rewards Button */}
            <div className="flex flex-row justify-center mt-8">
              <button
                className="py-1 px-2 rounded-md bg-emerald-400"
                onClick={() => calculateMaximumRewards()}
              >
                CALCULATE REWARDS
              </button>
            </div>
            {/* Maximum Rewards Response */}
            <div>
              <h1>Maximum Rewards</h1>
              <div>
                <ul>
                  {transactionGroups ? (
                    Array.from(transactionGroups.entries()).map(
                      ([date, transactionGroup]) => (
                        <TransactionGroupRow
                          date={date}
                          transactiongroup={transactionGroup}
                        />
                      )
                    )
                  ) : (
                    <></>
                  )}
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
