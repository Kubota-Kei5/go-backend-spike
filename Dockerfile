# 開発用ベースイメージ
FROM golang:1.24-bookworm as base

WORKDIR /spike-app

RUN apt-get update && \
    apt-get -y install locales wait-for-it && \
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8

ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9

# 本番環境の依存関係（ビルドに必要なライブラリ等）をインストール
FROM base as builder

COPY spike-app/go.mod spike-app/go.sum ./
RUN go mod download


# 開発環境
FROM builder as dev

# 開発用のツールをインストール
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

COPY spike-app/ .


# 本番環境用ビルドステージ
FROM builder AS build-prod

COPY spike-app/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# 本番環境のイメージ
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /app
COPY --from=build-prod /spike-app/main .
COPY --from=build-prod /spike-app/templates ./templates
ENV PORT=8080
EXPOSE 8080
CMD ["./main"]