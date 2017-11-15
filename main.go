package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
	"github.com/teambrookie/showrss/handlers"
	"github.com/teambrookie/showrss/torrent"

	"flag"
	"syscall"

	"cloud.google.com/go/firestore"

	"strconv"
)

const version = "1.0.0"

func worker(jobs <-chan dao.Episode, updateEpisode chan<- dao.Episode, client *firestore.Client) {
	time.Sleep(2 * time.Second)
	for episode := range jobs {
		time.Sleep(2 * time.Second)
		torrentLink, err := torrent.Search(strconv.Itoa(episode.ShowID), episode.Code, "720p")
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", episode.Name, err)
			continue
		}
		if torrentLink == "" {
			log.Printf("Cannot find : %s", episode.Name)
			continue
		}
		episode.MagnetLink = torrentLink
		episode.LastModified = time.Now()
		batch := client.Batch()
		oldRef := client.Collection("notFoundTorrents").Doc(episode.Name)
		newRef := client.Collection("foundTorrents").Doc(episode.Name)
		batch.Set(newRef, episode)
		batch.Delete(oldRef)
		_, err = batch.Commit(context.Background())
		updateEpisode <- episode
		if err != nil {
			log.Printf("Error saving %s to DB ...\n", episode.Name)
		}

	}
}

func updateEpisodeUsers(datastore *dao.Datastore, episodes <-chan dao.Episode) {
	for episode := range episodes {
		users, err := datastore.GetAllUsers()
		if err != nil {
			log.Println("Error retrieving user from firestore")
		}
		for _, user := range users {
			err := datastore.UpdateUserEpisode(user, episode)
			if err != nil && !strings.Contains(err.Error(), " code = NotFound") {
				log.Println(err)
			}
		}
	}
}

func main() {

	var httpAddr = flag.String("http", "0.0.0.0:8000", "HTTP service address")
	flag.Parse()

	apiKey := os.Getenv("BETASERIES_KEY")
	if apiKey == "" {
		log.Fatalln("BETASERIES_KEY must be set in env")
	}

	episodeProvider := betaseries.Betaseries{APIKey: apiKey}

	log.Println("Starting server ...")
	log.Printf("HTTP service listening on %s", *httpAddr)
	log.Println("Connecting to db ...")

	//Intialize Firestore client
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, "showrss-64e4b")
	if err != nil {
		log.Fatalf("Error when initializing the firestore client : %s\n", err)
	}
	log.Println("Firestore client initialized ...")

	datastore := dao.Datastore{client}

	// Worker stuff
	log.Println("Starting worker ...")
	jobs := make(chan dao.Episode, 1000)
	updateEpisode := make(chan dao.Episode, 100)
	go worker(jobs, updateEpisode, client)
	go updateEpisodeUsers(&datastore, updateEpisode)
	errChan := make(chan error, 10)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/auth", handlers.AuthHandler(&datastore, episodeProvider))
	mux.Handle("/refresh", handlers.RefreshHandler(&datastore, episodeProvider, jobs))
	mux.Handle("/episodes", handlers.EpisodeHandler(&datastore))
	mux.Handle("/rss", handlers.RSSHandler(&datastore))

	httpServer := http.Server{}
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
			httpServer.Shutdown(context.Background())
			os.Exit(0)
		}
	}

}
