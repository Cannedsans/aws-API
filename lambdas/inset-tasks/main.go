package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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
	lambda.Start(insert_tasks)
}

func insert_tasks(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var tarefa models.Tarefa

	err := json.Unmarshal([]byte(request.Body), &tarefa)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "JSON Inválido"}, nil
	}
	if tarefa.ID == "" {
		tarefa.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	err = models.SaveData(ctx, config.GetClient(), tarefa)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       `{"message": "Tarefa criada com sucesso!"}`,
	}, nil
}
