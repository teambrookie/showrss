package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
	"github.com/teambrookie/showrss/handlers"
	"github.com/teambrookie/showrss/worker"

	"syscall"

	"cloud.google.com/go/firestore"
)

const version = "1.0.0"

func main() {
	//need for querying the Betaseries API
	apiKey := os.Getenv("BETASERIES_KEY")
	if apiKey == "" {
		log.Fatalln("BETASERIES_KEY must be set in env")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	//small wrapper around the Betaseries API
	episodeProvider := betaseries.Betaseries{APIKey: apiKey}

	log.Println("Starting showrss ...")
	log.Printf("HTTP service listening on port  %s", port)

	//Intialize Firestore client
	client, err := firestore.NewClient(context.Background(), "showrss-64e4b")
	if err != nil {
		log.Fatalf("Error when initializing the firestore client : %s\n", err)
	}
	log.Println("Firestore connection OK ...")

	//wrapper around firestore client for convenience
	datastore := &dao.Datastore{Store: client}

	// Worker stuff
	log.Println("Starting worker ...")
	//torrentSearchs is a channel that will receive episode that need to be searched for a magnet link ( on rarbg at the moment)
	torrentSearchs := make(chan dao.Episode, 1000)
	// when the app find a magnet link for an episode it is send on these channel to update every user with the new episode ( fuck I suck at commenting ...)
	updateEpisode := make(chan dao.Episode, 100)
	//channel use as a rate limiter for the refresh worker, size is not 1 to allow mannual refresh by sending a tick in the channel ( maybe I don't suck completely)
	limiter := make(chan time.Time, 10)

	// goroutine that add a tick every hour to the limiter so that the refresh worker works
	go func() {
		for t := range time.Tick(time.Hour * 1) {
			limiter <- t
		}
	}()

	go worker.TorrentSearch(torrentSearchs, updateEpisode, client)
	go worker.UpdateEpisode(datastore, updateEpisode)
	go worker.Refresh(limiter, torrentSearchs, datastore, episodeProvider)

	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/auth", handlers.AuthHandler(datastore, episodeProvider))
	mux.Handle("/refreshes", handlers.RefreshHandler(limiter))
	mux.Handle("/{user}/episodes", handlers.EpisodeHandler(datastore))
	mux.Handle("/{user}/rss", handlers.RSSHandler(datastore))

	httpServer := http.Server{}
	httpServer.Addr = ":" + port
	httpServer.Handler = handlers.LoggingHandler(mux)

	errChan := make(chan error, 10)
	go func() {
		errChan <- httpServer.ListenAndServe()
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	for {
		select {
		// if we receive an error , log and exit
		case err := <-errChan:
			if err != nil {
				log.Fatal(err)
			}
		// if we receive a system signal, shutdown the httpserver and exit
		case s := <-signalChan:
			log.Println(fmt.Sprintf("Captured %v. Exiting...", s))
			httpServer.Shutdown(context.Background())
			os.Exit(0)
		}
	}

}
