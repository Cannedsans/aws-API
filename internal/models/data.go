package models

import (
	"context"
	"fmt" // Use fmt em vez de log.Fatal
	"strings"

	"github.com/Cannedsans/aws-API/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Tarefa struct {
    ID        string `json:"id" dynamodbav:"id"` 
    Tarefa    string `json:"tarefa" dynamodbav:"tarefa"`
    Descricao string `json:"descricao" dynamodbav:"descricao"`
    Feito     bool   `json:"feito" dynamodbav:"feito"`
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
func UpdateData(ctx context.Context, client *dynamodb.Client, tarefaInput Tarefa, id string) (Tarefa, error) {
    if id == "" {
        return Tarefa{}, fmt.Errorf("id do item não pode estar vazio")
    }

    var tarefaAtualizada Tarefa
    exprParts := []string{}
    attrValues := make(map[string]types.AttributeValue)

    // Montagem dinâmica dos campos
    if tarefaInput.Tarefa != "" {
        exprParts = append(exprParts, "tarefa = :t")
        attrValues[":t"] = &types.AttributeValueMemberS{Value: tarefaInput.Tarefa}
    }
    if tarefaInput.Descricao != "" {
        exprParts = append(exprParts, "descricao = :d")
        attrValues[":d"] = &types.AttributeValueMemberS{Value: tarefaInput.Descricao}
    }
    
    // Campo booleano sempre enviamos para garantir o estado
    exprParts = append(exprParts, "feito = :f")
    attrValues[":f"] = &types.AttributeValueMemberBOOL{Value: tarefaInput.Feito}

    updateExpression := "SET " + strings.Join(exprParts, ", ")

    result, err := client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
        TableName: aws.String(config.Table), // Certifique-se que config.Table está carregado
        Key: map[string]types.AttributeValue{
            "id": &types.AttributeValueMemberS{Value: id},
        },
        UpdateExpression:          aws.String(updateExpression),
        ExpressionAttributeValues: attrValues,
        ReturnValues:              types.ReturnValueAllNew,
    })

    if err != nil {
        return Tarefa{}, err
    }

    err = attributevalue.UnmarshalMap(result.Attributes, &tarefaAtualizada)
    return tarefaAtualizada, err
}