FROM golang:1.21


WORKDIR /app

COPY wait-for-it.sh /app/wait-for-it.sh
RUN chmod +x /app/wait-for-it.sh

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

CMD ["/app/wait-for-it.sh", "mysql:3306", "--", "/app/main"]