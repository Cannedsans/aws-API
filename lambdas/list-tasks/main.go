package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Cannedsans/aws-API/internal/config"
	"github.com/Cannedsans/aws-API/internal/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	err := config.InitAWS()
	if err != nil {
		panic("Erro ao iniciar AWS: " + err.Error())
	}
	lambda.Start(list_tasks)
}

func list_tasks(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Agora tratamos o erro sem quebrar a Lambda
	data, err := models.GetData(ctx, config.GetClient())
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       fmt.Sprintf(`{"error": "%s"}`, err.Error()),
		}, nil // Retornamos nil no error para o Gateway mostrar o nosso JSON de erro
	}

	jsonBody, _ := json.Marshal(data)
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(jsonBody),
	}, nil
}
