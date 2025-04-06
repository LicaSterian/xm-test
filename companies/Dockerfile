# build
FROM golang:1.23.4 AS builder

ENV GOOS=linux \
  GOARCH=amd64

# The Kafka client needs CGO_ENABLED
RUN apt-get update && apt-get install -y gcc librdkafka-dev

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -a -installsuffix=cgo -ldflags '-extldflags "-static"' -o companies .

# run
FROM gcr.io/distroless/base-debian10
COPY --from=builder /app/companies /companies
EXPOSE 8080
ENTRYPOINT ["/companies"]