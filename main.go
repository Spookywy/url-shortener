package main

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod != http.MethodGet {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Body:       "Method Not Allowed",
			Headers: map[string]string{
				"Content-Type": "text/plain",
			},
		}, nil
	}

	shorten_url := request.PathParameters["shorten_url"]

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusPermanentRedirect,
		Headers: map[string]string{
			"Location": shorten_url,
		},
	}, nil
}

func main() {
	lambda.Start(handler)
}
