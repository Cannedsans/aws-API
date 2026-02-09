package models

import (
	"context"
	"fmt" // Use fmt em vez de log.Fatal

	"github.com/Cannedsans/aws-API/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Tarefa struct {
	ID        string `dynamodbav:"id"`
	Tarefa    string `dynamodbav:"tarefa"`
	Descricao string `dynamodbav:"descricao"`
	Feito     bool   `dynamodbav:"feito"`
}

// CORREÇÃO: Retorne um erro em vez de matar o programa
func GetData(ctx context.Context, client *dynamodb.Client) ([]Tarefa, error) {
	if config.Table == "" {
		return nil, fmt.Errorf("TABLE_NAME não configurada")
	}

	res, err := client.Scan(ctx, &dynamodb.ScanInput{
		TableName: aws.String(config.Table),
	})
	if err != nil {
		return nil, fmt.Errorf("erro ao fazer scan: %w", err)
	}

	var tarefas []Tarefa
	err = attributevalue.UnmarshalListOfMaps(res.Items, &tarefas)
	if err != nil {
		return nil, fmt.Errorf("erro ao unmarshal: %w", err)
	}

	return tarefas, nil
}


func SaveData(ctx context.Context, client *dynamodb.Client, tarefa Tarefa) error {
	// Transforma a struct em um mapa compatível com o DynamoDB
	item, err := attributevalue.MarshalMap(tarefa)
	if err != nil {
		return err
	}

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(config.Table),
		Item:      item,
	})
	return err
}