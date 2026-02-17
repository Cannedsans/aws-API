# Api com aws

Uma api rest feita inteiramente com serviços da aws
serviços usados:

* lambdas functions
* API gateway
* dynamoDB

## Variavéis de ambiente

Para poder utilizar o projeto preencha o `.env`  com as configurações e o nome do banco de dados:
```env
  AWS_REGION=us-east-1
  AWS_ACCESS_KEY_ID=test
  AWS_SECRET_ACCESS_KEY=test
  AWS_ENDPOINT=http://host.docker.internal:4566
  TABLE_NAME=minha-tabela
```

## Rodando localmente

Para rodar pela primeira vez é nescessário criar o container e a tabela do banco de dados, rode:

``` bash
docker compose up -d
```

e para criar o banco de dados rode:

``` bash
./scripts/create-sources.sh   
```

caso já tenha rodado esses comandos pelo menos uma vez, poderá **compilar e configurar** o projeto com o comando:

``` bash
./scripts/build.sh
```

esse ultimo comando gera toda a configuração automáticamente, caso edite as variavéis no começo do código é possivél enviar o projeto diretamente para a aws real.