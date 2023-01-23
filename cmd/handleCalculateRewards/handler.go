package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/MowMowchow/c1rewards/internal/models"
	"github.com/MowMowchow/c1rewards/internal/responses"
	"github.com/MowMowchow/c1rewards/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody models.CalculateTransactionRequest
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		errMsg := "error unmarshalling response body from calculate transaction request"
		log.Println(errMsg)
		return responses.ServerError(errMsg), fmt.Errorf(errMsg)
	}

	if requestBody.Transactions == nil {
		errMsg := "transactions object does not exist in request"
		log.Println(errMsg)
		return responses.ServerError(errMsg), fmt.Errorf(errMsg)
	}

	if requestBody.Rules == nil {
		errMsg := "rules object does not exist in request"
		log.Println(errMsg)
		return responses.ServerError(errMsg), fmt.Errorf(errMsg)
	}

	// for time grouping:
	// date := now.Format("20060102")
	// fmt.Println(date)
	// date = now.Format("2006-01-02")
	// fmt.Println(date)

	// split at "-" -> get 3 groups: [year][month][day]

	var transactionsSummary models.TransactionsSummary
	transactionsSummary.Merchants = make(map[string]int)
	transactionsSummary.MaxCents = 0
	for _, transaction := range requestBody.Transactions {
		transactionsSummary.Merchants[transaction.MerchantCode] += transaction.AmountCents
		transactionsSummary.MaxCents = utils.IntMax(transactionsSummary.MaxCents, transactionsSummary.Merchants[transaction.MerchantCode])
	}

	// copy to make it easier to read (visually)
	var rules models.Rules
	leftoverRuleCost := 1
	leftoverRulePoints := 0
	for _, rule := range requestBody.Rules {
		rules.Rules = append(rules.Rules, rule)
		if _, exists := rule.Merchants["any"]; exists {
			leftoverRuleCost = rule.Merchants["any"]
			leftoverRulePoints = rule.Points
		}
	}
	rules.Length = len(rules.Rules)

	var calculateMaxRewards func(int, models.TransactionsSummary, int) int
	calculateMaxRewards = func(currPoints int, currTransactionsSummary models.TransactionsSummary, ruleInd int) int {
		// Check if we've gone through all the rules
		currBest := 0
		if ruleInd == rules.Length {
			totalLeftoverMoney := 0
			for _, leftoverMoney := range currTransactionsSummary.Merchants {
				totalLeftoverMoney += leftoverMoney
			}
			return currPoints + int(totalLeftoverMoney/leftoverRuleCost)*leftoverRulePoints
		}
		cost := 0
		pointsEarned := 0
		var canUseReward bool

		// uses: how many times we use the rule
		for uses := 0; uses < currTransactionsSummary.MaxCents; uses++ {
			canUseReward = true
			newTransactionsSummary := models.TransactionsSummary{
				Merchants: make(map[string]int),
				MaxCents:  int(currTransactionsSummary.MaxCents),
			}
			// copy current transation summary to new transaction summary
			for merchant, merchantCost := range currTransactionsSummary.Merchants {
				newTransactionsSummary.Merchants[merchant] = merchantCost
			}
			newTransactionsSummary.MaxCents = int(currTransactionsSummary.MaxCents)
			// iterate over all merchants required in the rule
			for merchant, merchantCost := range rules.Rules[ruleInd].Merchants {
				cost = merchantCost * uses
				pointsEarned = rules.Rules[ruleInd].Points * uses
				// check if the merchant exists is the transaction summary
				if _, merchantExists := currTransactionsSummary.Merchants[merchant]; merchantExists {
					// check if we've spent enough at this merchant to use this rule
					if currTransactionsSummary.Merchants[merchant] >= cost {
						newTransactionsSummary.Merchants[merchant] -= cost
					} else {
						canUseReward = false
					}
				} else {
					canUseReward = false
				}
			}
			if canUseReward {
				currBest = utils.IntMax(currBest, calculateMaxRewards(currPoints+pointsEarned, newTransactionsSummary, ruleInd+1))
			} else {
				currBest = utils.IntMax(currBest, calculateMaxRewards(currPoints, newTransactionsSummary, ruleInd+1))
				break
			}
		}
		return currBest
	}

	maxRewards := map[string]int{
		"best": calculateMaxRewards(0, transactionsSummary, 0),
	}

	// package request body
	responseBody, err := json.Marshal(maxRewards)
	if err != nil {
		errMsg := ("error marshalling calculate rewards to response body")
		log.Println(errMsg)
		return responses.ServerError(err), fmt.Errorf(errMsg)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: string(responseBody),
	}, nil
}
