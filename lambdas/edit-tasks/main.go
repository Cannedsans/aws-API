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
	lambda.Start(edit_task)
}

func edit_task(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
    var tarefa models.Tarefa

    // 1. Parse do Body
    err := json.Unmarshal([]byte(request.Body), &tarefa)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "JSON Inválido"}, nil
    }

    // 2. Tentar pegar o ID de PathParameters ou do próprio JSON
    id := request.PathParameters["id"]
    if id == "" {
        id = tarefa.ID // Fallback para o ID dentro do corpo do JSON
    }

    if id == "" {
        return events.APIGatewayProxyResponse{StatusCode: 400, Body: "ID não fornecido na URL ou no JSON"}, nil
    }

    // 3. Chamar o Update
    tarefaAtualizada, err := models.UpdateData(ctx, config.Cliente, tarefa, id)
    if err != nil {
        return events.APIGatewayProxyResponse{StatusCode: 500, Body: "Erro no Dynamo: " + err.Error()}, nil
    }

    body, _ := json.Marshal(tarefaAtualizada)
    return events.APIGatewayProxyResponse{
        StatusCode: 200,
        Body:       string(body),
        Headers:    map[string]string{"Content-Type": "application/json"},
    }, nil
}
