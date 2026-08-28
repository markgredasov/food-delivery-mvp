FROM golang:1.26-alpine AS builder
WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /out/restaurant-simulator ./cmd/simulator

FROM alpine:3.22
RUN apk add --no-cache ca-certificates
COPY --from=builder /out/restaurant-simulator /usr/local/bin/restaurant-simulator

EXPOSE 8081
ENTRYPOINT ["restaurant-simulator"]
