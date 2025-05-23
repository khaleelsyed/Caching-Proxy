package main

import (
	"flag"
	"fmt"
	"time"

	"go.uber.org/zap"
)

type ServerParams struct {
	Port           int
	ForwardingHost string
}

func getFlags() (string, string) {
	port := flag.Int("port", 3000, "The port to run the proxy server on")
	forwardURL := flag.String("destination", "https://echo.hoppscotch.io", "The destination URL to forward all requests to")

	flag.Parse()

	listenAddr := fmt.Sprintf("localhost:%d", *port)

	return listenAddr, *forwardURL
}

func main() {
	var err error

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar := logger.Sugar()

	listenAddr, forwardURL := getFlags()
	sugar.Infof("Running a caching proxy server on %s", listenAddr)
	sugar.Infof("Forwarding all requests to %s", forwardURL)

	proxy := NewProxyServer(forwardURL, 30*time.Second)

	err = proxy.StartServer(listenAddr)
	if err != nil {
		sugar.Panic(err)
	}
}
