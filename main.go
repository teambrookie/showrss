package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"showrss/dao"
	"showrss/handlers"
	"showrss/torrent"
	"time"

	"flag"

	"syscall"

	"github.com/braintree/manners"
)

const version = "1.0.0"

func worker(jobs <-chan dao.Episode, db *dao.DB) {
	for episode := range jobs {
		log.Println("Processing : " + episode.Name)
		torrentLink, err := torrent.Search(episode.Name)
		log.Println("Result : " + torrentLink)
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", episode.Name, err)
			continue
		}
		if torrentLink == "" {
			continue
		}
		episode.MagnetLink = torrentLink
		err = db.AddEpisode(episode)
		if err != nil {
			log.Printf("Error saving %s to DB ...\n", episode.Name)
		}
		time.Sleep(2 * time.Second)
	}
}

func main() {
	var httpAddr = flag.String("http", "0.0.0.0:8000", "HTTP service address")
	var dbAddr = flag.String("db", "showrss.db", "DB address")
	flag.Parse()

	log.Println("Starting server ...")
	log.Printf("HTTP service listening on %s", *httpAddr)
	log.Println("Connecting to db ...")

	//DB stuff
	db, err := dao.InitDB(*dbAddr)
	if err != nil {
		log.Fatalln("Error connecting to DB")
	}

	err = db.CreateBucket("episodes")
	if err != nil {
		log.Fatalln("Error when creating bucket")
	}

	// Worker stuff
	log.Println("Starting worker ...")
	jobs := make(chan dao.Episode, 100)
	go worker(jobs, db)

	errChan := make(chan error, 10)

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.HandleFunc("/auth", handlers.AuthHandler)
	mux.Handle("/refresh", handlers.RefreshHandler(db, jobs))
	mux.Handle("/episodes", handlers.EpisodeHandler(db))
	mux.Handle("/rss", handlers.RSSHandler(db))

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
