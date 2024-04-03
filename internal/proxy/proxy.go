package proxy

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func HandleProxy(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	uid := strings.TrimPrefix(path, "/")

	target, err := url.Parse("http://example-svc.default.svc.cluster/" + uid)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// WebSocket 지원을 위한 설정
	proxy.ModifyResponse = func(resp *http.Response) error {
		if resp.Header.Get("Upgrade") == "websocket" {
			resp.Header.Set("Connection", "Upgrade")
		}
		return nil
	}

	// WebSocket 핸드셰이크 처리
	if r.Header.Get("Connection") == "Upgrade" && r.Header.Get("Upgrade") == "websocket" {
		proxy.ServeHTTP(w, r)
		return
	}

	// 일반적인 HTTP 요청 처리
	proxy.ServeHTTP(w, r)
}
