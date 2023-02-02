package main

import (
	"encoding/json"
	"log"
	"sort"
	"strings"

	"github.com/MowMowchow/c1rewards/internal/constants"
	"github.com/MowMowchow/c1rewards/internal/models"
	"github.com/MowMowchow/c1rewards/internal/responses"
	"github.com/MowMowchow/c1rewards/internal/utils"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody models.CalculateRewardsRequest
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		errMsg := "error unmarshalling response body from calculate transaction request"
		log.Println(errMsg, err)
		return responses.ServerError(), nil
	}

	if requestBody.Transactions == nil {
		errMsg := "transactions object does not exist in request"
		log.Println(errMsg)
		return responses.BadRequest(), nil
	}
	if requestBody.Rules == nil {
		errMsg := "rules object does not exist in request"
		log.Println(errMsg)
		return responses.BadRequest(), nil
	}
	if requestBody.Grouping == "" {
		errMsg := "no grouping specified in request"
		log.Println(errMsg)
		return responses.BadRequest(), nil
	}

	// copy to make it easier to read
	transactions := []models.Transaction{}
	for transactionName, transaction := range requestBody.Transactions {
		transaction.Name = transactionName
		transactions = append(transactions, transaction)
	}

	// copy to make it easier to read
	rules := []models.RewardRule{}
	leftoverRuleCost := 1
	leftoverRulePoints := 0
	for _, rule := range requestBody.Rules {
		rules = append(rules, rule)
		if _, exists := rule.Merchants[constants.AnyMerchant]; exists {
			leftoverRuleCost = rule.Merchants[constants.AnyMerchant]
			leftoverRulePoints = rule.Points
		}
	}

	var calculateMaxRewards func(int, models.TransactionsSummary, int) int
	calculateMaxRewards = func(currPoints int, currTransactionsSummary models.TransactionsSummary, ruleInd int) int {
		// Check if we've gone through all the rules
		if ruleInd == len(rules) {
			totalLeftoverMoney := 0
			for _, leftoverMoney := range currTransactionsSummary.Merchants {
				totalLeftoverMoney += leftoverMoney
			}
			return currPoints + int(totalLeftoverMoney/leftoverRuleCost)*leftoverRulePoints
		}

		currMaxRewards := 0
		cost := 0
		pointsEarned := 0
		var canUseRule bool
		// uses: how many times we use the rule
		for uses := 0; uses < currTransactionsSummary.MaxCentsSpent; uses++ {
			canUseRule = true
			newTransactionsSummary := models.TransactionsSummary{
				Merchants:     make(map[string]int),
				MaxCentsSpent: currTransactionsSummary.MaxCentsSpent,
			}
			// copy current transation summary to new transaction summary
			for merchant, merchantCost := range currTransactionsSummary.Merchants {
				newTransactionsSummary.Merchants[merchant] = merchantCost
			}
			// iterate over all merchants required in the rule
			for merchant, merchantCost := range rules[ruleInd].Merchants {
				cost = merchantCost * uses // required expenditure
				// check if the merchant exists is the transaction summary AND if we've spent enough at this merchant to use this rule
				if _, merchantExists := currTransactionsSummary.Merchants[merchant]; merchantExists && currTransactionsSummary.Merchants[merchant] >= cost {
					newTransactionsSummary.Merchants[merchant] -= cost
				} else {
					canUseRule = false
				}
			}

			pointsEarned = rules[ruleInd].Points * uses
			if canUseRule {
				currMaxRewards = utils.IntMax(currMaxRewards, calculateMaxRewards(currPoints+pointsEarned, newTransactionsSummary, ruleInd+1))
			} else {
				currMaxRewards = utils.IntMax(currMaxRewards, calculateMaxRewards(currPoints, newTransactionsSummary, ruleInd+1))
				break
			}
		}
		return currMaxRewards
	}

	// sort transactions
	sort.Slice(transactions, func(i, j int) bool {
		merchantIDateArr := strings.Split(transactions[i].Date, "-")
		merchantJDateArr := strings.Split(transactions[j].Date, "-")
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

	for i := 0; i < len(transactions); i++ {
		switch transactionGroups.Grouping {
		case constants.TransactionGroupingYear:
			dateKey := strings.Join(strings.Split(transactions[i].Date, "-")[:1], "-")
			if _, exists := transactionGroups.Groups[dateKey]; !exists {
				transactionGroups.Groups[dateKey] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[dateKey].Transactions = append(transactionGroups.Groups[dateKey].Transactions, transactions[i])
		case constants.TransactionGroupingMonth:
			dateKey := strings.Join(strings.Split(transactions[i].Date, "-")[:2], "-")
			if _, exists := transactionGroups.Groups[dateKey]; !exists {
				transactionGroups.Groups[dateKey] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[dateKey].Transactions = append(transactionGroups.Groups[dateKey].Transactions, transactions[i])
		case constants.TransactionGroupingDay:
			if _, exists := transactionGroups.Groups[transactions[i].Date]; !exists {
				transactionGroups.Groups[transactions[i].Date] = &models.TransactionGroup{}
			}
			transactionGroups.Groups[transactions[i].Date].Transactions = append(transactionGroups.Groups[transactions[i].Date].Transactions, transactions[i])
		default:
			errMsg := "invalid grouping supplied"
			log.Println(errMsg, err)
			return responses.BadRequest(), nil
		}
	}

	for dateKey, transactionGroup := range transactionGroups.Groups {
		var transactionsSummary models.TransactionsSummary
		transactionsSummary.Merchants = make(map[string]int)
		transactionsSummary.MaxCentsSpent = 0

		for i, transaction := range transactionGroup.Transactions {
			transactionGroups.Groups[dateKey].Transactions[i].MaxRewards = calculateMaxRewards(
				// curr points
				0,
				// transaction summary containing only the current transaction
				models.TransactionsSummary{
					Merchants: map[string]int{
						transaction.MerchantCode: transaction.AmountCents,
					},
					MaxCentsSpent: transaction.AmountCents,
				},
				// curr rule -> 0-indexed
				0,
			)
			if _, merchantExists := transactionsSummary.Merchants[transaction.MerchantCode]; !merchantExists {
				transactionsSummary.Merchants[transaction.MerchantCode] = 0
			}
			transactionsSummary.MaxCentsSpent = utils.IntMax(transactionsSummary.MaxCentsSpent, transaction.AmountCents)
			transactionsSummary.Merchants[transaction.MerchantCode] += transaction.AmountCents
			transactionGroups.Groups[dateKey].MaxRewards = utils.IntMax(transactionsSummary.MaxCentsSpent, transactionsSummary.Merchants[transaction.MerchantCode])
		}
		transactionGroups.Groups[dateKey].MaxRewards = calculateMaxRewards(0, transactionsSummary, 0)
	}

	calculateRewardsResponse := models.CalculateRewardsResponse{
		Groups:   make(map[string]models.TransactionGroup),
		Grouping: transactionGroups.Grouping,
	}
	calculateRewardsResponse.Grouping = transactionGroups.Grouping
	for date, group := range transactionGroups.Groups {
		calculateRewardsResponse.Groups[date] = *group
	}

	// package response body
	responseBody, err := json.Marshal(calculateRewardsResponse)
	if err != nil {
		errMsg := ("error marshalling calculate rewards to response body")
		log.Println(errMsg, err)
		return responses.ServerError(), nil
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
