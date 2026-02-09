package config

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/joho/godotenv"
)

var (
	Cliente *dynamodb.Client
	Table   string
	Cfg     aws.Config // Mudei o nome para não conflitar com o pacote 'config'
)

func InitAWS() error {
	_ = godotenv.Load() // Na Lambda, o .env precisa estar na raiz do ZIP

	endpoint := os.Getenv("AWS_ENDPOINT_URL")
    if endpoint == "" {
        endpoint = os.Getenv("AWS_ENDPOINT")
    }
    if endpoint == "" {
        endpoint = "http://172.17.0.1:4566" 
    }
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}
	Table = os.Getenv("TABLE_NAME")
	
	ctx := context.TODO()
	
	// CORREÇÃO: Usando '=' para não criar variável local (Shadowing)
	var err error
	Cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return err
	}

	Cliente = dynamodb.NewFromConfig(Cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return nil
}

func GetClient() *dynamodb.Client {
	return Cliente
}