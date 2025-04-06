# build
FROM golang:1.23.4-alpine AS builder

ENV CGO_ENABLED=0 \
  GOOS=linux \
  GOARCH=amd64

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o auth .

# run
FROM scratch
COPY --from=builder /app/auth /auth
EXPOSE 80
ENTRYPOINT ["/auth"]