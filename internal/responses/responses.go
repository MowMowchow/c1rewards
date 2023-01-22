package responses

import (
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type MessageBody struct {
	Message interface{} `json:"message"`
}

const (
	BadRequestMessage  = "sadly your message was bad :)"
	ServerErrorMessage = "oops :O"
)

func makeResponse(statusCode int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Headers: map[string]string{
			"Access-Control-Allow-Origin":      "*",
			"Access-Control-Allow-Headers":     "*",
			"Access-Control-Allow-Credentials": "true",
		},
		Body: body,
	}
}

func makeMessageBody(bodyMessage interface{}) string {
	body := MessageBody{
		Message: bodyMessage,
	}
	bodyBytes, _ := json.Marshal(body)
	return string(bodyBytes)
}

func Ok(body []byte) events.APIGatewayProxyResponse {
	return makeResponse(http.StatusOK, string(body))
}

func Created(body []byte) events.APIGatewayProxyResponse {
	return makeResponse(http.StatusCreated, string(body))
}

func Accepted(body []byte) events.APIGatewayProxyResponse {
	return makeResponse(http.StatusAccepted, string(body))
}

func BadRequest(message ...interface{}) events.APIGatewayProxyResponse {
	responseBody := makeMessageBody(BadRequestMessage)
	if message != nil {
		responseBody = makeMessageBody(message[0])
	}
	return makeResponse(http.StatusBadRequest, responseBody)
}

func ServerError(message ...interface{}) events.APIGatewayProxyResponse {
	responseBody := makeMessageBody(ServerErrorMessage)
	if message != nil {
		responseBody = makeMessageBody(message[0])
	}
	return makeResponse(http.StatusInternalServerError, responseBody)
}
