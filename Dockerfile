FROM golang:1.17.6-alpine as builder

WORKDIR /app/source

# restore dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy source
COPY . .

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -v -o infinitude-prometheus .


FROM alpine as final

WORKDIR /app

EXPOSE 8080

COPY --from=builder /app/source/infinitude-prometheus ./

ENTRYPOINT ./infinitude-prometheus
