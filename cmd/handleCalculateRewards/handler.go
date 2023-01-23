package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strings"

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
		return responses.BadRequest(errMsg), fmt.Errorf(errMsg)
	}

	if requestBody.Rules == nil {
		errMsg := "rules object does not exist in request"
		log.Println(errMsg)
		return responses.BadRequest(errMsg), fmt.Errorf(errMsg)
	}

	transactionList := []models.Transaction{}
	for transactionName, transaction := range requestBody.Transactions {
		transaction.Name = transactionName
		transactionList = append(transactionList, transaction)
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

	// sort transactions
	sort.Slice(transactionList, func(i, j int) bool {
		merchantIDateArr := strings.Split(transactionList[i].Date, "-")
		merchantJDateArr := strings.Split(transactionList[j].Date, "-")
		if merchantIDateArr[0] == merchantJDateArr[0] { // year is the same
			if merchantIDateArr[1] == merchantJDateArr[1] { // month is the same
				return merchantIDateArr[2] < merchantJDateArr[2] // compare days
			} else {

				return merchantIDateArr[1] < merchantJDateArr[1] // compare months
			}
		} else {

			return merchantIDateArr[0] < merchantJDateArr[0] // compare years
		}
	})

	transactionGroups := models.TransactionGroups{
		Groups:   make(map[string]*models.TransactionGroup),
		Grouping: requestBody.Grouping,
	}
	transactionGroups.Grouping = requestBody.Grouping

	for i := 0; i < len(transactionList); i++ {
		if transactionGroups.Grouping == "day" {
			if _, exists := transactionGroups.Groups[transactionList[i].Date]; !exists {
				transactionGroups.Groups[transactionList[i].Date] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[transactionList[i].Date].Transactions = append(transactionGroups.Groups[transactionList[i].Date].Transactions, transactionList[i])
		} else if transactionGroups.Grouping == "year" {
			dateKey := strings.Join(strings.Split(transactionList[i].Date, "-")[:1], "-")
			if _, exists := transactionGroups.Groups[dateKey]; !exists {
				transactionGroups.Groups[dateKey] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[dateKey].Transactions = append(transactionGroups.Groups[dateKey].Transactions, transactionList[i])
		} else { // month is default case
			dateKey := strings.Join(strings.Split(transactionList[i].Date, "-")[:2], "-")
			if _, exists := transactionGroups.Groups[dateKey]; !exists {
				transactionGroups.Groups[dateKey] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[dateKey].Transactions = append(transactionGroups.Groups[dateKey].Transactions, transactionList[i])
		}
	}

	for dateKey, transactionGroup := range transactionGroups.Groups {
		var transactionsSummary models.TransactionsSummary
		transactionsSummary.Merchants = make(map[string]int)
		transactionsSummary.MaxCents = 0

		for i, transaction := range transactionGroup.Transactions {
			transactionGroups.Groups[dateKey].Transactions[i].MaxRewards = calculateMaxRewards(
				0,
				models.TransactionsSummary{
					Merchants: map[string]int{
						transaction.MerchantCode: transaction.AmountCents,
					},
					MaxCents: transaction.AmountCents,
				},
				0,
			)
			if _, merchantExists := transactionsSummary.Merchants[transaction.MerchantCode]; !merchantExists {
				transactionsSummary.Merchants[transaction.MerchantCode] = 0
			}
			transactionsSummary.MaxCents = utils.IntMax(transactionsSummary.MaxCents, transaction.AmountCents)
			transactionsSummary.Merchants[transaction.MerchantCode] += transaction.AmountCents
			transactionGroups.Groups[dateKey].MaxRewards = utils.IntMax(transactionsSummary.MaxCents, transactionsSummary.Merchants[transaction.MerchantCode])
		}
		transactionGroups.Groups[dateKey].MaxRewards = calculateMaxRewards(0, transactionsSummary, 0)
	}

	outTransactionGroups := models.OutTransactionGroups{
		Groups:   make(map[string]models.TransactionGroup),
		Grouping: transactionGroups.Grouping,
	}
	outTransactionGroups.Grouping = transactionGroups.Grouping
	for date, group := range transactionGroups.Groups {
		outTransactionGroups.Groups[date] = *group
	}

	// package request body
	responseBody, err := json.Marshal(outTransactionGroups)
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
