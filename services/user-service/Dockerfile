FROM golang:1.22.7-alpine AS build

WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o user-service ./cmd/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=build /app/user-service .

EXPOSE 8081

# RUN chmod +x deployment-service
CMD ["./user-service"]