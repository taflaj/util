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
	myServer = server.NewServer(":"+port, &handlers)
	myServer.SetOnStart(func() { log.Printf("%v listening on port %v", os.Getpid(), port) })
	myServer.SetOnExit(func() { log.Print("Server is exiting now") })
	myServer.SetOnFail(func(err error) { log.Print(err) })
	myServer.SetOnInterrupt(func(sig os.Signal) { log.Printf("Received %v", sig) })
	log.Printf("%#v", myServer)
	myServer.Start()
}
