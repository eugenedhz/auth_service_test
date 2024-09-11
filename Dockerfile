FROM golang:1.23.1

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go mod tidy

RUN go build -o ./build cmd/api/main.go

EXPOSE ${SERVER_PORT}

CMD ["./build"]