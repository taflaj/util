// main.go

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/taflaj/util/server"
)

var myServer server.Server

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

type genericHandler struct{}

func (h *genericHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%#v", r)
	fmt.Fprintln(w, "Ok")
}

type exitHandler struct{}

func (h *exitHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%#v", r)
	fmt.Fprintln(w, "Exiting")
	go func() {
		time.Sleep(time.Second)
		myServer.Stop()
		os.Exit(0)
	}()
}

func main() {
	var handlers server.Handlers
	handlers = append(handlers, server.ApiHandler{Pattern: "/", CallHandler: &genericHandler{}})
	handlers = append(handlers, server.ApiHandler{Pattern: "/exit", CallHandler: &exitHandler{}})
	port := "8000"
	myServer = server.NewServer(":"+port, &handlers).
		SetOnStart(func() {
			log.Printf("%v listening on port %v", os.Getpid(), port)
		}).
		SetOnExit(func() {
			log.Print("Server is exiting now")
		}).
		SetOnFail(func(err error) {
			log.Print(err)
		}).
		SetOnInterrupt(func(sig os.Signal) {
			log.Printf("Received %v", sig)
		})
	log.Printf("%#v", myServer)
	myServer.StartAndWait()
}
