FROM golang:latest as builder
WORKDIR /build/

COPY . .

RUN go mod tidy && go mod vendor && go mod download
RUN go install github.com/daixiang0/gci@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.0
RUN gci write --skip-generated -s standard -s default .
RUN golangci-lint run
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o /app/onlinestore .

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl
WORKDIR /app
COPY --from=builder /app/onlinestore .

CMD ["./onlinestore"]
