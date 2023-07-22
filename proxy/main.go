package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"
)

var (
	// Hardcoded for POC purpose only.
	basicAuthUserName = "johndoe"
	basicAuthPassword = "supersecret"
	proxyPort         = 8080
	logger            = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	auth.ProxyBasic(proxy, "realm", func(user, password string) bool {
		return user == basicAuthUserName && password == basicAuthPassword
	})

	addr := fmt.Sprintf("127.0.0.1:%d", proxyPort)
	logger.Printf("Starting proxy server on %s\n", addr)
	logger.Fatalln(http.ListenAndServe(addr, proxy))
}
