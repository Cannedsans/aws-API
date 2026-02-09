#!/bin/bash
# comandos CLI para criar a tabela e as lambdas

aws dynamodb create-table --profile localstack\
    --table-name minha-tabela \
    --attribute-definitions AttributeName=id,AttributeType=S \
    --key-schema AttributeName=id,KeyType=HASH \
    --provisioned-throughput ReadCapacityUnits=5,WriteCapacityUnits=5 >> /dev/null
