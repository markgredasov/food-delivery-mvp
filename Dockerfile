FROM golang:1.26.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/cmd/app/exe /app/cmd/app/main.go
CMD ["./app"]

FROM alpine:3.24
WORKDIR /app
COPY --from=builder /app/cmd/app/exe /app
CMD ["/app/exe"] 
