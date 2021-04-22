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
	GetServer() *http.Server
	SetCloseTimeout(time.Duration)
	SetOnExit(func())
	SetOnFail(func(error))
	SetOnInterrupt(func(os.Signal))
	SetOnStart(func())
	Start()
	Stop()
}

type Handler struct {
	Pattern string
	Handler func(http.ResponseWriter, *http.Request)
}

type Handlers []Handler

type setup struct {
	server       *http.Server
	closeTimeout time.Duration
	onStart      func()
	onExit       func()
	onFail       func(error)
	onInterrupt  func(os.Signal)
}

func NewServer(address string, handlers *Handlers) Server {
	mux := http.NewServeMux()
	for _, handler := range *handlers {
		mux.HandleFunc(handler.Pattern, handler.Handler)
	}
	server := setup{}
	server.server = &http.Server{Addr: address, Handler: mux}
	return &server
}

func (server *setup) Start() {
	s := server.server
	if s.IdleTimeout == 0 {
		s.IdleTimeout = 60 * time.Second
	}
	if s.ReadTimeout == 0 {
		s.ReadTimeout = 15 * time.Second
	}
	if s.WriteTimeout == 0 {
		s.WriteTimeout = 15 * time.Second
	}
	if server.closeTimeout == 0 {
		server.closeTimeout = 15 * time.Second
	}
	go func() {
		if server.onStart != nil {
			server.onStart()
		}
		if err := server.server.ListenAndServe(); err != nil {
			if server.onFail != nil {
				server.onFail(err)
			}
		}
	}()
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGUSR1)
	sig := <-stop
	if server.onInterrupt != nil {
		server.onInterrupt(sig)
	}
	server.Stop()
}

func (server *setup) Stop() {
	if server.onExit != nil {
		server.onExit()
	}
	ctx, cancel := context.WithTimeout(context.Background(), server.closeTimeout)
	defer cancel()
	if err := server.server.Shutdown(ctx); err != nil && server.onFail != nil {
		server.onFail(err)
	}
}

func (server *setup) GetServer() *http.Server {
	return server.server
}

func (server *setup) SetCloseTimeout(t time.Duration) {
	server.closeTimeout = t
}

func (server *setup) SetOnStart(f func()) {
	server.onStart = f
}

func (server *setup) SetOnExit(f func()) {
	server.onExit = f
}

func (server *setup) SetOnFail(f func(error)) {
	server.onFail = f
}

func (server *setup) SetOnInterrupt(f func(os.Signal)) {
	server.onInterrupt = f
}
