FROM golang:1.25.5-nanoserver AS golang-builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN ["CGO_ENABLED=0 GOOS=linux GOARCH=amd64", "go build -o image-resizing-go ./cmd/api/main.go"]

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=golang-builder --chmod=755 /app/image-resizing-go .
COPY .env* ./
EXPOSE 8080
CMD ["./image-resizing-go"]