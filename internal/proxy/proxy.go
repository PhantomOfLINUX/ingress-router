package proxy

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/PhantomOfLINUX/ingressRouter/internal/handler"
)

const dnsFormat = "http://%s-%s.default.svc.cluster.local:8080"

func HandleProxy(w http.ResponseWriter, r *http.Request) {
    queryParams := r.URL.Query()
    uid := queryParams.Get("uid")
    stage := queryParams.Get("stage")

    // uid 유효성 검사
    if !isValidParam(uid) || !isValidParam(stage) {
        log.Printf("PSHELL_NOT_FOUND, uid=%s, stage=%s\n", uid, stage)
		handler.RespondWithError(w, http.StatusBadRequest, "PSHELL_NOT_FOUND", "4602_PSHELL_NOT_FOUND", "")
		return
    }

    // 프록시 경로 설정
	targetURL := fmt.Sprintf(dnsFormat, stage, uid)
	target, err := url.Parse(targetURL)
	if err != nil {
		log.Println(err)
		handler.RespondWithError(w, http.StatusInternalServerError, "Internal Server Error", "5000_INTERNAL_SERVER_ERROR", "")
		return
	}
	log.Printf("TargetURL: %s\n", targetURL)

	// 리버스 프록시
	proxy := httputil.NewSingleHostReverseProxy(target)

	// WebSocket 지원을 위한 설정
	proxy.ModifyResponse = modifyResponse

	// 에러 처리 로직 추가
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		handler.RespondWithError(w, http.StatusBadGateway, "Proxy error", "5001_PROXY_ERROR", "")
	}

	// 요청 헤더 출력
	go logRequestHeaders(r)

	// WebSocket 핸드셰이크 처리
	if isWebSocketRequest(r) {
		ctx, cancel := context.WithCancel(r.Context())
		defer cancel()
		r = r.WithContext(ctx)
	}

	proxy.ServeHTTP(w, r)
}

func isValidParam(value string) bool {
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

	logBuffer.WriteString(r.Method + " " + r.URL.Path + " " + r.Proto + "\n")
	for key, values := range r.Header {
		for _, value := range values {
			logBuffer.WriteString(key + ": " + value + "\n")
		}
	}

	log.Println(logBuffer.String())
}
