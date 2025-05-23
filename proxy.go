package main

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ProxyServer struct {
	forwardURL string
	cache      Cache
}

func (p *ProxyServer) proxyRequest(w http.ResponseWriter, r *http.Request) {
	var response *http.Response
	var err error

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	targetURL := fmt.Sprintf("%s%s", p.forwardURL, r.URL.Path)

	response, cacheHit := p.cache.Read(targetURL, r.Header)
	if cacheHit {
		writeOKResponse(w, response)
		return
	}

	response, err = http.Get(targetURL)
	if err != nil {
		sugar.Error(err)
		return
	}

	writeOKResponse(w, response)
	p.cache.Write(targetURL, response, r.Header)
}

func (p *ProxyServer) StartServer(listenAddr string) error {
	http.HandleFunc("GET /", p.proxyRequest)
	return http.ListenAndServe(listenAddr, nil)
}

func NewProxyServer(forwardingHost string, cacheExpireDuration time.Duration) *ProxyServer {
	return &ProxyServer{
		forwardURL: forwardingHost,
		cache:      NewCache(cacheExpireDuration),
	}
}

func writeOKResponse(w http.ResponseWriter, response *http.Response) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		sugar.Error(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(body))
}
