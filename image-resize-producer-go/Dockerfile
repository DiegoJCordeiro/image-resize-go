FROM golang:tip-bookworm AS builder
WORKDIR /app

RUN apt-get update && apt-get install -y gcc

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/google/wire/cmd/wire@latest
RUN cd ./cmd/api && wire
RUN go build -o /app/image-resizing-go-producer-api ./cmd/api/

FROM debian:bookworm-slim
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
WORKDIR /root/

ARG profileArg=docker
ENV profile=${profileArg}

COPY --from=builder --chmod=755 /app/image-resizing-go-producer-api .
COPY --from=builder --chmod=755 /app/application-${profile}.yaml .
COPY --from=builder --chmod=755 /app/.env.${profile} .

EXPOSE 8080
CMD ["./image-resizing-go-producer-api"]