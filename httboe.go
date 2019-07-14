package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
)

var srvName = "httboe"

func mainHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter {
			w,
			http.StatusOK,
		}
		w.Header().Set("Server", srvName)
		next.ServeHTTP(lrw, r)
		log.Printf("[%s]: %s %d\n", r.Method, r.URL, lrw.statusCode)
	})
}

type loggingResponseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
    lrw.statusCode = code
    lrw.ResponseWriter.WriteHeader(code)
}

func main() {

	httpPort := flag.Int("p", 80, "Port to serve on (Plain HTTP)")

	flag.Parse()
	srvMux := http.NewServeMux()

	srv := &webdav.Handler{
		Prefix: "/webdav/",
		FileSystem: webdav.Dir("./webdav"),
		LockSystem: webdav.NewMemLS(),		
	}

	// keep last '/'
	srvMux.Handle("/webdav/", srv)

	srvMux.Handle("/", http.FileServer(http.Dir("./public")))

	// Create a server listening
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", *httpPort),
		Handler: mainHandler(srvMux),
	}

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
