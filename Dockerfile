FROM docker.io/library/golang:1.24.1-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY ./config/.env ./config/

RUN go mod tidy
RUN apk --no-cache add curl

COPY . .
WORKDIR /app/cmd/api/

RUN CGO_ENABLED=0 GOOS=linux go build -o app 
 
EXPOSE 8080

# Comando final para esperar o banco de dados e depois iniciar a aplicação
CMD ["./app"]
