package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/BKellogg/iugaserver/handlers"
)

const (
	apiPathV1 = "/api/v1/"
)

func main() {
	addr := getRequiredENV("ADDR", ":443")
	tlsKey := getRequiredENV("TLSKEY", "")
	tlsCert := getRequiredENV("TLSCERT", "")

	// redirect all HTTP requests to HTTPS
	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(handlers.HTTPSRedirect)); err != nil {
			log.Fatal(err)
		}
	}()

	mux := http.NewServeMux()

	attachServiceFromENVToMux("IUGASITEADDR", "/", mux)
	attachServiceFromENVToMux("IUGAEVENTSADDR", apiPathV1+"events", mux)
	attachServiceFromENVToMux("NODEMICROSERVICEADDR", "/node", mux)

	log.Printf("server is listening at %s...\n", addr)
	log.Fatal(http.ListenAndServeTLS(addr, tlsCert, tlsKey, mux))
}

// returns a proxy of the given address
func getServiceProxy(addr string) *httputil.ReverseProxy {
	return &httputil.ReverseProxy{
		Director: func(r *http.Request) {
			r.URL.Scheme = "http"
			r.URL.Host = addr
			r.URL.Path = trimAPIPrefix(r.URL.Path)
		},
	}
}

// Gets the service address from the given env environment variable and attaches
// it to the given mux at the given addr. Does nothing if there is no environment
// variable for env.
func attachServiceFromENVToMux(env, addr string, mux *http.ServeMux) {
	serviceAddr := os.Getenv(env)
	if len(serviceAddr) == 0 {
		return
	}
	mux.Handle(addr, getServiceProxy(serviceAddr))
}

// Trim the API prefix off the given addr and return it
func trimAPIPrefix(addr string) string {
	if !strings.HasPrefix(addr, "/") {
		addr = "/" + addr
	}
	trimmed := strings.TrimPrefix(addr, apiPathV1)
	if !strings.HasPrefix(trimmed, "/") {
		return "/" + trimmed
	}
	return trimmed
}

// Gets and returns a required environment variable. If the
// environment variable is empty, the default will be used. If
// both are empty then the process will exit.
func getRequiredENV(name string, def string) string {
	val := os.Getenv(name)
	if len(val) == 0 {
		if len(def) > 0 {
			return def
		}
		log.Fatalf("please set the %s environment variable", name)
	}
	return val
}
