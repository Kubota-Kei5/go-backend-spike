#!/bin/sh

# .env ファイルから環境変数（PROD_DB_USER 等）を読み込む
source ./.env

DB_NAME=spike-app-1-prod
DB_USER=postgres

docker run --rm -it --net=host postgres:16 psql "host=127.0.0.1 sslmode=disable dbname=${DB_NAME} user=${DB_USER}"