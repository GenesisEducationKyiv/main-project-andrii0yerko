package app

import (
	"context"
	"log"
	"net/http"
	"time"
)

const readTimeout = 3 * time.Second

var ErrServerClosed = http.ErrServerClosed

type Server struct {
	handler *ExchangeRateHandler
	mux     *http.ServeMux
	server  *http.Server
}

func NewServer(handler *ExchangeRateHandler, addr string) *Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", AllowMethods(handler.GetRoot, http.MethodGet))
	mux.HandleFunc("/rate", AllowMethods(handler.GetRate, http.MethodGet))
	mux.HandleFunc("/subscribe", AllowMethods(handler.PostSubscribe, http.MethodPost))
	mux.HandleFunc("/sendEmails", AllowMethods(handler.PostSendEmails, http.MethodPost))

	return &Server{
		handler: handler,
		mux:     mux,
		server:  &http.Server{Addr: addr, Handler: mux, ReadHeaderTimeout: readTimeout},
	}
}

func (s *Server) Start() error {
	log.Printf("Running on http://%s\n", s.server.Addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	return s.server.Shutdown(context.TODO())
}

func AllowMethods(handler http.HandlerFunc, methods ...string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, method := range methods {
			if r.Method == method {
				handler(w, r)
				return
			}
		}
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
