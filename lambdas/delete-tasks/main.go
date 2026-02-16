package main

import (
	"context"
	"encoding/json"

	// Adicionado para formatação
	"github.com/Cannedsans/aws-API/internal/config"
	"github.com/Cannedsans/aws-API/internal/models"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	// Inicializa a configuração (Certifique-se que o client do Dynamo está acessível)
	err := config.InitAWS()
	if err != nil {
		panic("Erro ao iniciar AWS: " + err.Error())
	}

	// Passa a função, não o resultado da execução!
	lambda.Start(delete_task)
}

func delete_task(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var id string

	if id == "" && request.Body != "" {
		var tarefa models.Tarefa
		err := json.Unmarshal([]byte(request.Body), &tarefa)
		if err == nil {
			id = tarefa.ID
		}
	}

	if id == "" {
		return events.APIGatewayProxyResponse{StatusCode: 400, Body: "ID não fornecido na URL ou no JSON"}, nil
	}

	// 3. Chamar o Update
	err := models.DeleData(ctx, config.Cliente, id)
	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Erro no Dynamo: " + err.Error()}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       `{"message": "Tarefa deletada com sucesso!"}`,
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}
