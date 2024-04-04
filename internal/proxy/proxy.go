package proxy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const dnsFormat = "http://%s-svc-%s.default.svc.cluster.local:8080"

func HandleProxy(w http.ResponseWriter, r *http.Request) {
    uid := r.Header.Get("X-POL-UID")
	stage := r.Header.Get("X-POL-STAGE")

    // uid 유효성 검사
    if !isValidHeader(uid) || !isValidHeader(stage) {
        log.Printf("Invalid headers: %s %s", uid, stage)
        http.Error(w, "Invalid Headers", http.StatusBadRequest)
        return
    }

	// 프록시 경로 설정
    targetURL := fmt.Sprintf(dnsFormat, stage, uid)
    target, err := url.Parse(targetURL)
    if err != nil {
        log.Println(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
	log.Println("targetURL" + targetURL)

    proxy := httputil.NewSingleHostReverseProxy(target)

    // WebSocket 지원을 위한 설정
    proxy.ModifyResponse = modifyResponse

    // 에러 처리 로직 추가
    proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
        log.Printf("Proxy error: %v", err)
        http.Error(w, "Proxy error", http.StatusBadGateway)
    }

    logRequestHeaders(r)

    // WebSocket 핸드셰이크 처리
    if isWebSocketRequest(r) {
        ctx, cancel := context.WithCancel(r.Context())
        defer cancel()
        r = r.WithContext(ctx)
    }

    proxy.ServeHTTP(w, r)
}

func isValidHeader(value string) bool {
    return value != ""
}

func modifyResponse(resp *http.Response) error {
    if resp.Header.Get("Upgrade") == "websocket" {
        resp.Header.Set("Connection", "Upgrade")
    }
    return nil
}

func isWebSocketRequest(r *http.Request) bool {
    return r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "websocket"
}

func logRequestHeaders(r *http.Request) {
    var logBuffer bytes.Buffer

    logBuffer.WriteString(r.Method + " " + r.URL.Path + " " + r.Proto + "\n\n")
    for key, values := range r.Header {
        for _, value := range values {
            logBuffer.WriteString(key + ": " + value + "\n")
        }
    }
    logBuffer.WriteString("\n")

    log.Println(logBuffer.String())
}