FROM golang:1.22.7-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/main.go
ENV REDIS_ADDR redis:6379
ENV AUTH_SERVICE_PORT 8080

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/auth-service .

EXPOSE 8080
# RUN chmod +x deployment-service
CMD ["./auth-service"]