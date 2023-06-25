package http

import (
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
	mux.HandleFunc("/", handler.GetRoot)
	mux.HandleFunc("/rate", handler.GetRate)
	mux.HandleFunc("/subscribe", handler.PostSubscribe)
	mux.HandleFunc("/sendEmails", handler.PostSendEmails)

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
