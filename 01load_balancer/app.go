package mybalancer

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

var mu sync.Mutex
var idx int = 0

// lbHandler is the handler for loadbalancing
func lbHandler(w http.ResponseWriter, r *http.Request) {
	maxLen := len(cfg.Backends)

	// Round Robin
	mu.Lock()

	// currentBackend := cfg.Backends[idx%maxLen]

	//choosing next backend server by incremented idx
	targetURL, err := url.Parse(cfg.Backends[idx%maxLen].URL)
	if err != nil {
		log.Fatal(err.Error())
	}
	idx++
	mu.Unlock()

	reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)

	reverseProxy.ServeHTTP(w, r)
}

var cfg Config

func Serve() {
	cfg.Init()

	//director config func for reverse Proxy
	//changes incoming request's host and Scheme to actual servers which are running
	director := func(request *http.Request) {
		request.URL.Scheme = "http"
		request.URL.Host = ":8081"
	}

	//lbHandler is an HTTP Handler that takes an incoming request and
	// sends it to another server, proxying(redo it again from the server to client) the response back to the client.
	lbHandler := &httputil.ReverseProxy{
		Director: director,
	}

	server := http.Server{
		Addr:    ":8080",
		Handler: lbHandler,
	}

	//starting the server
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err.Error())
	}
}
