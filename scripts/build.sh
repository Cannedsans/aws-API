#!/bin/bash

# Configurações
ENDPOINT="--endpoint-url=http://localhost:4566"
PROFILE="--profile localstack"
REGION="us-east-1"
ACCOUNT_ID="000000000000"

echo "🚀 Iniciando Build e Limpeza..."
rm -rf build && mkdir -p build/list build/create

# 1. Compilação Estática (CGO_ENABLED=0 é vital para Linux/LocalStack)
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "⚙️ Compilando binários Go..."
go build -ldflags="-s -w" -o build/list/bootstrap lambdas/list-tasks/main.go
go build -ldflags="-s -w" -o build/create/bootstrap lambdas/inset-tasks/main.go

# 2. Empacotamento
cp .env build/list/.env
cp .env build/create/.env

cd build/list && zip -j list_function.zip bootstrap .env && cd ../..
cd build/create && zip -j create_function.zip bootstrap .env && cd ../..

echo "🧹 Removendo recursos antigos..."
aws $ENDPOINT $PROFILE lambda delete-function --function-name list_tasks_lambda 2>/dev/null
aws $ENDPOINT $PROFILE lambda delete-function --function-name create_task_lambda 2>/dev/null

# 3. Criando Lambdas
echo "📦 Enviando Lambdas para o LocalStack..."

aws $ENDPOINT $PROFILE lambda create-function \
    --function-name list_tasks_lambda \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://build/list/list_function.zip \
    --role arn:aws:iam::$ACCOUNT_ID:role/lambda-role

aws $ENDPOINT $PROFILE lambda create-function \
    --function-name create_task_lambda \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://build/create/create_function.zip \
    --role arn:aws:iam::$ACCOUNT_ID:role/lambda-role

# 4. Configurando API Gateway
echo "🌐 Configurando API Gateway..."

API_ID=$(aws $ENDPOINT $PROFILE apigateway create-rest-api --name 'TaskAPI' --query 'id' --output text)
PARENT_ID=$(aws $ENDPOINT $PROFILE apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text)
RESOURCE_ID=$(aws $ENDPOINT $PROFILE apigateway create-resource --rest-api-id $API_ID --parent-id $PARENT_ID --path-part tasks --query 'id' --output text)

# --- MÉTODOS ---

# GET (Listar)
aws $ENDPOINT $PROFILE apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method GET --authorization-type "NONE"
aws $ENDPOINT $PROFILE apigateway put-integration --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method GET --type AWS_PROXY --integration-http-method POST --uri arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:$ACCOUNT_ID:function:list_tasks_lambda/invocations

# POST (Criar)
aws $ENDPOINT $PROFILE apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method POST --authorization-type "NONE"
aws $ENDPOINT $PROFILE apigateway put-integration --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method POST --type AWS_PROXY --integration-http-method POST --uri arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:$ACCOUNT_ID:function:create_task_lambda/invocations

# --- PERMISSÕES (Crucial para evitar Error 500/502) ---
aws $ENDPOINT $PROFILE lambda add-permission --function-name list_tasks_lambda --statement-id apigateway-get --action lambda:InvokeFunction --principal apigateway.amazonaws.com
aws $ENDPOINT $PROFILE lambda add-permission --function-name create_task_lambda --statement-id apigateway-post --action lambda:InvokeFunction --principal apigateway.amazonaws.com

# 5. Deploy Final
aws $ENDPOINT $PROFILE apigateway create-deployment --rest-api-id $API_ID --stage-name prod

echo "--------------------------------------------------"
echo "✅ Deploy Finalizado!"
echo "🔗 URL (GET/POST):"
echo "http://localhost:4566/restapis/$API_ID/prod/_user_request_/tasks"
echo "--------------------------------------------------"