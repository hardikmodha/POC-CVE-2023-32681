package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	certFile = "/path/to/cert.pem"
	keyFile  = "/path/to/cert.key"

	redirectionServerPort = 443
	echoServerPort        = 4431
	logger                = log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	errs := make(chan error)

	redirectionServer := http.NewServeMux()
	redirectionServer.HandleFunc("/redirect", redirectHandler())

	echoServer := http.NewServeMux()
	echoServer.HandleFunc("/echoHeaders", echoHeadersHandler())

	startServer(redirectionServer, redirectionServerPort, errs)
	startServer(echoServer, echoServerPort, errs)

	logger.Fatalln(<-errs)
}

func startServer(mux *http.ServeMux, port int32, errs chan<- error) {
	go func() {
		defer func() {
			if x := recover(); x != nil {
				errs <- fmt.Errorf("panic: %v", x)
			}
		}()

		logger.Printf("Starting server on port %s\n", port)

		errs <- http.ListenAndServeTLS(fmt.Sprintf(":%d", port), certFile, keyFile, mux)
	}()
}

func redirectHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Println("[Redirection Server] Redirecting request")
		http.Redirect(w, r, "https://127.0.0.1:4431/echoHeaders", http.StatusMovedPermanently)
	}
}

func echoHeadersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		mData, err := json.MarshalIndent(r.Header, "", " ")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		logger.Printf("[Echo Headers Server] %s", string(mData))

		if _, err := w.Write(mData); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}
