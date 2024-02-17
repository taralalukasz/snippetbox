#Stage 1
FROM golang:latest AS builder

WORKDIR /app

COPY . .

RUN go get -d -v ./...

RUN go build -o snippetbox ./cmd/web 

#Stage 2
FROM golang:latest

ARG APP_NAME=snippetbox

WORKDIR /app

COPY --from=builder /app/snippetbox .
COPY --from=builder /app/cmd ./cmd

EXPOSE 4000

CMD ./snippetbox