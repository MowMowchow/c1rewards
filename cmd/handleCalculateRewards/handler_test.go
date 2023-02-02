package main

import (
	"encoding/json"
	"log"
	"testing"

	"github.com/MowMowchow/c1rewards/internal/constants"
	"github.com/MowMowchow/c1rewards/internal/models"
	"github.com/aws/aws-lambda-go/events"
)

func TestHappyPath(t *testing.T) {
	transactions := map[string]models.Transaction{
		"T01": {
			Name:         "T01",
			Date:         "2021-05-01",
			MerchantCode: "sportcheck",
			AmountCents:  21000,
		},
		"T02": {
			Name:         "T02",
			Date:         "2021-05-02",
			MerchantCode: "sportcheck",
			AmountCents:  8700,
		},
		"T03": {
			Name:         "T03",
			Date:         "2021-05-03",
			MerchantCode: "tim_hortons",
			AmountCents:  323,
		},
		"T04": {
			Name:         "T04",
			Date:         "2021-05-04",
			MerchantCode: "tim_hortons",
			AmountCents:  1267,
		},
		"T05": {
			Name:         "T05",
			Date:         "2021-05-05",
			MerchantCode: "tim_hortons",
			AmountCents:  2116,
		},
		"T06": {
			Name:         "T06",
			Date:         "2021-05-06",
			MerchantCode: "tim_hortons",
			AmountCents:  2211,
		},
		"T07": {
			Name:         "T07",
			Date:         "2021-05-07",
			MerchantCode: "subway",
			AmountCents:  1853,
		},
		"T08": {
			Name:         "T08",
			Date:         "2021-05-08",
			MerchantCode: "subway",
			AmountCents:  2153,
		},
		"T09": {
			Name:         "T09",
			Date:         "2021-05-09",
			MerchantCode: "sportcheck",
			AmountCents:  7326,
		},
		"T10": {
			Name:         "T10",
			Date:         "2021-05-10",
			MerchantCode: "tim_hortons",
			AmountCents:  1321,
		},
	}
	rules := []models.RewardRule{
		{
			Name:   "rule 1",
			Points: 500,
			Merchants: map[string]int{
				"sportcheck":  7500,
				"tim_hortons": 2500,
				"subway":      2500,
			},
		},
		{
			Name:   "rule 2",
			Points: 300,
			Merchants: map[string]int{
				"sportcheck":  7500,
				"tim_hortons": 2500,
			},
		},
		{
			Name:   "rule 3",
			Points: 200,
			Merchants: map[string]int{
				"sportcheck": 7500,
			},
		},
		{
			Name:   "rule 4",
			Points: 150,
			Merchants: map[string]int{
				"sportcheck":  2500,
				"tim_hortons": 1000,
				"subway":      1000,
			},
		},
		{
			Name:   "rule 5",
			Points: 75,
			Merchants: map[string]int{
				"sportcheck":  2500,
				"tim_hortons": 1000,
				"subway":      1000,
			},
		},
		{
			Name:   "rule 6",
			Points: 75,
			Merchants: map[string]int{
				"sportcheck": 2000,
			},
		},
		{
			Name:   "rule 7",
			Points: 1,
			Merchants: map[string]int{
				"anyMerchant": 1,
			},
		},
	}
	groupingType := constants.TransactionGroupingMonth

	transactionRequest := models.CalculateRewardsRequest{
		Transactions: transactions,
		Rules:        rules,
		Grouping:     groupingType,
	}

	requestBody, err := json.Marshal(transactionRequest)
	if err != nil {
		errMsg := ("error marshalling unit test request body")
		log.Println(errMsg, err)
	}

	request := events.APIGatewayProxyRequest{
		Body: string(requestBody),
	}

	handler := Handler{}
	response, err := handler.HandleRequest(request)

	if response.StatusCode != 200 {
		t.Error("error | happy path test | expected 200 response")
	}

	var responseBody models.TransactionGroups
	err = json.Unmarshal([]byte(response.Body), &responseBody)
	if err != nil {
		errMsg := "error unmarshalling response body from calculate transaction response"
		log.Println(errMsg, err)
	}

	if transactionGroup, exists := responseBody.Groups["2021-05"]; exists {
		if transactionGroup.MaxRewards != 48270 {
			t.Error("error | happy path test | expected transaction group MaxRewards == 48270, got", transactionGroup.MaxRewards)
		}
	} else {
		t.Error("error | happy path test | date grouping 2021-05 does not exist")
	}

}
