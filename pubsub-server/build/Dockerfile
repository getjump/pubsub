FROM golang:1.18.3 as builder

COPY . /app
WORKDIR /app

CMD go mod download && go run /app/cmd/grpc_app.go

EXPOSE ${APP_GRPC_PORT} ${APP_PORT}