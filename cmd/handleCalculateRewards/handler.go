package main

import (
	"encoding/json"
	"fmt"

	"github.com/MowMowchow/c1rewards/internal/models"
	"github.com/aws/aws-lambda-go/events"
)

type Handler struct {
}

func (h *Handler) HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var requestBody models.CalculateTransactionRequest
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		errMsg := "error unmarshalling response body from calculate transaction request"
		return responses.ServerError(errMsg), fmt.Errorf(errMsg)
	}

	if requestBody.Transactions != nil {

	}
}
