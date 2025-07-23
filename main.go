package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

func main() {
	listenPort := os.Getenv("PORT")
	if listenPort == "" {
		listenPort = ":8080"
	}

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	log.Println("Listening at " + listenPort)
	if err := http.ListenAndServe(listenPort, &handler{hostname}); err != nil {
		panic(err)
	}
}

type handler struct {
	hostname string
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s][%s] %s", realIPFromRequest(r), r.Method, r.URL.Path)

	query := map[string]string{}
	header := map[string]string{}
	for k := range r.URL.Query() {
		query[k] = r.URL.Query().Get(k)
	}
	for k := range r.Header {
		header[k] = r.Header.Get(k)
	}

	request := map[string]any{
		"method":  r.Method,
		"uri":     r.RequestURI,
		"path":    r.URL.Path,
		"query":   query,
		"headers": header,
	}
	if r.URL.Query().Get("raw") == "1" {
		request["raw"] = map[string]any{
			"query":   r.URL.Query(),
			"headers": r.Header,
		}
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]any{
		"hostname": h.hostname,
		"request":  request,
	}); err != nil {
		panic(err)
	}
}

func realIPFromRequest(r *http.Request) string {
	const (
		headerXForwardedFor = "X-Forwarded-For"
		headerXRealIP       = "X-Real-IP"
	)

	if ip := r.Header.Get(headerXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := r.Header.Get(headerXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ra
}
