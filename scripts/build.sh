#!/bin/bash

# Configurações
ENDPOINT="--endpoint-url=http://localhost:4566"
PROFILE="--profile localstack"
REGION="us-east-1"
ACCOUNT_ID="000000000000"

echo "🚀 Iniciando Build e Limpeza..."
rm -rf build && mkdir -p build/list build/create build/edit

# 1. Compilação Estática
export GOOS=linux
export GOARCH=amd64
export CGO_ENABLED=0

echo "⚙️ Compilando binários Go..."
# Verifique se o nome da pasta é 'inset' ou 'insert' (corrigi para o que estava no seu script)
go build -ldflags="-s -w" -o build/list/bootstrap lambdas/list-tasks/main.go
go build -ldflags="-s -w" -o build/create/bootstrap lambdas/inset-tasks/main.go
go build -ldflags="-s -w" -o build/edit/bootstrap lambdas/edit-tasks/main.go
go build -ldflags="-s -w" -o build/delete/bootstrap lambdas/delete-tasks/main.go


echo "📦 Zipando arquivos..."
cp .env build/list/.env && cd build/list && zip -j list_function.zip bootstrap .env && cd ../..
cp .env build/create/.env && cd build/create && zip -j create_function.zip bootstrap .env && cd ../..
cp .env build/edit/.env && cd build/edit && zip -j edit_function.zip bootstrap .env && cd ../..
cp .env build/delete/.env && cd build/delete && zip -j delete_function.zip bootstrap .env && cd ../..


echo "🧹 Removendo recursos antigos..."
aws $ENDPOINT $PROFILE lambda delete-function --function-name list_tasks_lambda 2>/dev/null
aws $ENDPOINT $PROFILE lambda delete-function --function-name create_task_lambda 2>/dev/null
aws $ENDPOINT $PROFILE lambda delete-function --function-name edit_task_lambda 2>/dev/null
aws $ENDPOINT $PROFILE lambda delete-function --function-name delete_task_lambda 2>/dev/null

# Nota: Como estamos criando uma nova REST API a cada execução, os IDs mudam, o que evita conflitos no API Gateway.

# 3. Criando Lambdas
echo "📤 Enviando Lambdas..."

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

aws $ENDPOINT $PROFILE lambda create-function \
    --function-name edit_task_lambda \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://build/edit/edit_function.zip \
    --role arn:aws:iam::$ACCOUNT_ID:role/lambda-role7

aws $ENDPOINT $PROFILE lambda create-function \
    --function-name delete_task_lambda \
    --runtime provided.al2023 \
    --handler bootstrap \
    --zip-file fileb://build/delete/delete_function.zip \
    --role arn:aws:iam::$ACCOUNT_ID:role/lambda-role

# 4. Configurando API Gateway
echo "🌐 Configurando API Gateway..."

API_ID=$(aws $ENDPOINT $PROFILE apigateway create-rest-api --name 'TaskAPI' --query 'id' --output text)
PARENT_ID=$(aws $ENDPOINT $PROFILE apigateway get-resources --rest-api-id $API_ID --query 'items[0].id' --output text)
RESOURCE_ID=$(aws $ENDPOINT $PROFILE apigateway create-resource --rest-api-id $API_ID --parent-id $PARENT_ID --path-part tasks --query 'id' --output text)

# Função auxiliar para configurar Métodos + Integração + Permissão
configurar_metodo() {
    METHOD=$1
    LAMBDA_NAME=$2
    
    # Criar Método
    aws $ENDPOINT $PROFILE apigateway put-method --rest-api-id $API_ID --resource-id $RESOURCE_ID --http-method $METHOD --authorization-type "NONE"
    
    # Criar Integração (A URI corrigida com o caminho completo da Lambda)
    aws $ENDPOINT $PROFILE apigateway put-integration \
        --rest-api-id $API_ID \
        --resource-id $RESOURCE_ID \
        --http-method $METHOD \
        --type AWS_PROXY \
        --integration-http-method POST \
        --uri "arn:aws:apigateway:$REGION:lambda:path/2015-03-31/functions/arn:aws:lambda:$REGION:$ACCOUNT_ID:function:$LAMBDA_NAME/invocations"

    # Adicionar Permissão (Statement ID único para evitar ResourceConflict)
    aws $ENDPOINT $PROFILE lambda add-permission \
        --function-name $LAMBDA_NAME \
        --statement-id "apigateway-invoke-$METHOD-$API_ID" \
        --action lambda:InvokeFunction \
        --principal apigateway.amazonaws.com \
        --source-arn "arn:aws:execute-api:$REGION:$ACCOUNT_ID:$API_ID/*/$METHOD/tasks"
}

echo "🛠️ Configurando endpoints"
configurar_metodo "GET" "list_tasks_lambda"
configurar_metodo "POST" "create_task_lambda"
configurar_metodo "PUT" "edit_task_lambda"
configurar_metodo "DELETE" "delete_task_lambda"


# 5. Deploy Final
echo "🚀 Criando Deployment..."
aws $ENDPOINT $PROFILE apigateway create-deployment --rest-api-id $API_ID --stage-name prod

echo "--------------------------------------------------"
echo "✅ Deploy Finalizado!"
echo "🔗 URL BASE:"
echo "http://localhost:4566/restapis/$API_ID/prod/_user_request_/tasks"
echo "--------------------------------------------------"