FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o addrgo .

FROM gcr.io/distroless/base-debian11 AS runtime
WORKDIR /app
COPY --from=builder /app/addrgo .
EXPOSE 8080
USER nonroot:nonroot
CMD ["./addrgo"]
