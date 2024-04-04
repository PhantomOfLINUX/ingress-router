FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -v -o ingress-router main.go

FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY --from=builder /app/ingress-router /ingress-router

CMD ["/ingress-router"]
