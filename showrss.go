package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"showrss/handlers"

	"flag"

	"syscall"

	"github.com/braintree/manners"
)

const version = "1.0.0"

func main() {
	var httpAddr = flag.String("http", "0.0.0.0:8000", "HTTP service address")
	flag.Parse()

	log.Println("Starting server ...")
	log.Printf("HTTP service listening on %s", *httpAddr)

	errChan := make(chan error, 10)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	httpServer := manners.NewServer()
	httpServer.Addr = *httpAddr
	httpServer.Handler = handlers.LoggingHandler(mux)

	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.BlockingClose()
			os.Exit(0)
		}
	}

}
