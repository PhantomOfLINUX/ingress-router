package proxy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func HandleProxy(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path
    uid := strings.TrimPrefix(path, "/")

    // uid 유효성 검사
    // if !isValidUid(uid) {
    //     log.Printf("Invalid uid: %s", uid)
    //     http.Error(w, "Invalid uid", http.StatusBadRequest)
    //     return
    // }

    // uidInt, _ := strconv.Atoi(uid)
    targetURL := fmt.Sprintf("http://svc-%s.default.cluster.local", uid)
	log.Println(targetURL)
    target, err := url.Parse(targetURL)
    if err != nil {
        log.Println(err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

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

// func isValidUid(uid string) bool {
//     uidInt, err := strconv.Atoi(uid)
//     if err != nil {
//         return false
//     }
//     return uidInt >= 10000 && uidInt <= 30000
// }

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