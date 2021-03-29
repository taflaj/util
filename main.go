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

var (
	myServer server.Server
	setup    *server.Setup
)

func init() {
	log.SetFlags(log.Flags() | log.Lmicroseconds)
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%#v", r)
	fmt.Fprint(w, "Ok")
}

func exit(w http.ResponseWriter, r *http.Request) {
	log.Printf("%#v", r)
	fmt.Fprint(w, "Exiting")
	go func() {
		time.Sleep(time.Second)
		myServer.Stop()
		os.Exit(0)
	}()
}

func main() {
	var handlers server.Handlers
	handlers = append(handlers, server.Handler{Pattern: "/", Handler: handler})
	handlers = append(handlers, server.Handler{Pattern: "/exit", Handler: exit})
	port := "8000"
	myServer, setup = server.NewServer(":"+port, &handlers)
	setup.OnStart = func() { log.Printf("%v listening on port %v", os.Getpid(), port) }
	setup.OnExit = func() { log.Print("Server is exiting now") }
	setup.OnFail = func(err error) { log.Print(err) }
	setup.OnInterrupt = func(sig os.Signal) { log.Printf("Received %v", sig) }
	log.Printf("%#v", myServer)
	myServer.Start()
}
