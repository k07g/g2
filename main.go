package main

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Name string `json:"name"`
}

type Response struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func handler(ctx context.Context, req Request) (Response, error) {
	log.Printf("Received request: %+v", req)

	name := req.Name
	if name == "" {
		name = "World"
	}

	return Response{
		Message:    fmt.Sprintf("Hello, %s!", name),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
