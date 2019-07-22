package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/net/webdav"
	"httboe/httboe"
)

var srvName = "httboe"

func mainHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &loggingResponseWriter{
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

	confFile := flag.String("c", "main.conf", "Conf file")
	flag.Parse()

	log.Printf("Reading conf file %s\n", *confFile)

	conf := &httboe.Conf{}
	err := conf.Load(*confFile)

	if err != nil {
		log.Fatalf("Error parsing file %s: %s", *confFile, err.Error())
	}

	srvMux := http.NewServeMux()

	for _, loc := range conf.Server.Location {
		if loc.Type == "webdav" {
			srv := &webdav.Handler{
				Prefix:     loc.Path,
				FileSystem: webdav.Dir(loc.Root),
				LockSystem: webdav.NewMemLS(),
			}
			srvMux.Handle(loc.Path, srv)
		} else {
			srvMux.Handle(loc.Path, http.FileServer(http.Dir(loc.Root)))
		}
	}

	// Create a server listening
	s := &http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Server.Port),
		Handler: mainHandler(srvMux),
	}

	log.Printf("Listen on port %d", conf.Server.Port)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

}
