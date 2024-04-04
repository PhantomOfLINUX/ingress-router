FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o pol-proxy ./

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /app/pol-proxy /pol-proxy

CMD ["./pol-proxy"]
