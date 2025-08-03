FROM golang:1.24-bookworm AS builder

WORKDIR /spike-app

RUN apt-get update && \
    apt-get -y install locales wait-for-it && \
    localedef -f UTF-8 -i ja_JP ja_JP.UTF-8

ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8
ENV TZ JST-9


COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./app.go

EXPOSE 8080

CMD ["./app"]