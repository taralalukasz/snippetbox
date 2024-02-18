#Stage 1
FROM golang:latest AS builder

ARG APP_NAME=snippetbox

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o ${APP_NAME} ./cmd/web 

#Stage 2
FROM golang:latest

ARG APP_NAME snippetbox
ENV appname=${APP_NAME}
WORKDIR /app

COPY --from=builder /app/${APP_NAME} .
COPY --from=builder /app/cmd ./cmd

EXPOSE 4000

CMD ./$appname