// server.go

package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server interface {
	Start()
	Stop()
}

type Handler struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

type Handlers []Handler

type Setup struct {
	Server       *http.Server
	CloseTimeout time.Duration
	OnStart      func()
	OnExit       func()
	OnFail       func(error)
	OnInterrupt  func(os.Signal)
}

func NewServer(address string, handlers *Handlers) (Server, *Setup) {
	mux := http.NewServeMux()
	for _, handler := range *handlers {
		mux.HandleFunc(handler.Pattern, handler.Handler)
	}
	server := Setup{}
	server.Server = &http.Server{Addr: address, Handler: mux}
	return &server, &server
}

func (server *Setup) Start() {
	s := server.Server
	if s.IdleTimeout == 0 {
		s.IdleTimeout = 60 * time.Second
	}
	if s.ReadTimeout == 0 {
		s.ReadTimeout = 15 * time.Second
	}
	if s.WriteTimeout == 0 {
		s.WriteTimeout = 15 * time.Second
	}
	if server.CloseTimeout == 0 {
		server.CloseTimeout = 15 * time.Second
	}
	go func() {
		if server.OnStart != nil {
			server.OnStart()
		}
		if err := server.Server.ListenAndServe(); err != nil {
			if server.OnFail != nil {
				server.OnFail(err)
			}
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	sig := <-stop
	if server.OnInterrupt != nil {
		server.OnInterrupt(sig)
	}
	server.Stop()
}

func (server *Setup) Stop() {
	if server.OnExit != nil {
		server.OnExit()
	}
	ctx, cancel := context.WithTimeout(context.Background(), server.CloseTimeout)
	defer cancel()
	if err := server.Server.Shutdown(ctx); err != nil && server.OnFail != nil {
		server.OnFail(err)
	}
}
