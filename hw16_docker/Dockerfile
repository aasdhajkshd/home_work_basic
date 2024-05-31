FROM golang:1.21.5 as builder
WORKDIR /build/

COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG GIT_COMMIT
ARG BUILD_TIME
ARG GIT_TAG

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$GIT_TAG -X main.gitCommit=$GIT_COMMIT -X main.buildTime=$BUILD_TIME" -a -installsuffix cgo -o app main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY . .
COPY --from=builder /build/app/ ./
WORKDIR /root/

CMD ["./app"]