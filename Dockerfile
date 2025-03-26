FROM golang:1.24.1-alpine AS builder

WORKDIR /app/source

# restore dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o infinitude-prometheus .


FROM alpine AS final

WORKDIR /app

EXPOSE 8080

COPY --from=builder /app/source/infinitude-prometheus ./

ENTRYPOINT ["./infinitude-prometheus"]
