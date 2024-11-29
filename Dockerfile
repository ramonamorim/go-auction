FROM golang:1.23.2 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o auction cmd/auction/main.go

FROM scratch

WORKDIR /

COPY --from=builder /app/auction /auction
COPY --from=builder /app/cmd/auction/.env /cmd/auction/.env


ENTRYPOINT ["/auction"]
