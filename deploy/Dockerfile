FROM docker.io/library/golang:1.24.1-alpine

WORKDIR /app

COPY go.mod  /app
COPY go.sum  /app
COPY ./config/.env /app

RUN go mod tidy

COPY . .

WORKDIR /app/cmd/api/

RUN CGO_ENABLED=0 GOOS=linux go build -o app 
 
EXPOSE 8080

CMD ["./app"]
