FROM golang:latest 

WORKDIR /go/src/app

COPY . .

RUN go get -d -v ./...

RUN go build -o snippetbox ./cmd/web 

EXPOSE 4000

CMD ["./snippetbox"]