FROM golang:1.22 as builder

WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

# 정적 링크된 바이너리 생성
RUN CGO_ENABLED=0 go build -v -tags netgo -ldflags '-extldflags "-static"' -o pol-proxy ./

# Distroless static 이미지 사용
FROM gcr.io/distroless/static

WORKDIR /

COPY --from=builder /app/pol-proxy /pol-proxy

CMD ["./pol-proxy"]
