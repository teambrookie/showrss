package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/teambrookie/showrss/betaseries"
	"github.com/teambrookie/showrss/dao"
	"github.com/teambrookie/showrss/handlers"
	"github.com/teambrookie/showrss/torrent"
	"golang.org/x/oauth2"

	"flag"

	"syscall"

	"strconv"
)

const version = "1.0.0"

func handleNewAuth(newAuth <-chan string, users map[string]bool, refreshLimiter chan<- time.Time) {
	for token := range newAuth {
		if exists := users[token]; !exists {
			users[token] = true
			log.Printf("New user token : %s\n", token)
			refreshLimiter <- time.Now()
		}
	}
}

func searchWorker(jobs <-chan dao.Episode, store dao.EpisodeStore, quality string) {
	for episode := range jobs {
		time.Sleep(2 * time.Second)
		log.Println("Processing : " + episode.Name)
		torrentLink, err := torrent.Search(strconv.Itoa(episode.ShowID), episode.Code, quality)
		log.Println("Result : " + torrentLink)
		if err != nil {
			log.Printf("Error processing %s : %s ...\n", episode.Name, err)
			continue
		}
		if torrentLink == "" {
			continue
		}
		episode.MagnetLink = torrentLink
		episode.LastModified = time.Now()
		err = store.UpdateEpisode(episode)
		if err != nil {
			log.Printf("Error saving %s to DB ...\n", episode.Name)
		}

	}
}

func refresh(limiter <-chan time.Time, users map[string]bool, db dao.EpisodeStore, betaseries betaseries.EpisodeProvider, episodeToSearch chan<- dao.Episode) {
	for {
		<-limiter
		log.Println("Refresh started")
		for user := range users {
			log.Printf("Refresing for user %s\n", user)
			episodes, err := betaseries.Episodes(user)
			if err != nil {
				log.Printf("Error retriving episodes for user %s : %s\n", user, err)
			}
			for _, ep := range episodes {
				err := db.AddEpisode(ep)
				if err != nil {
					log.Printf("Error adding episodes to database: %s", err)
				}
			}
		}
		log.Println("Passing not found episodes to the search worker")
		notFounds, err := db.GetAllNotFoundEpisode()
		if err != nil {
			log.Printf("Error retriving unfound episodes from db : %s\n", err)
		}
		for _, episode := range notFounds {
			episodeToSearch <- episode
		}

	}
}

func main() {

	//Opitional flag for passing the http server address and the db name
	var dbAddr = flag.String("db", "showrss.db", "DB address")
	flag.Parse()

	//API key and secret for Betaseries are retrieve from the environnement variables
	apiKey := os.Getenv("BETASERIES_KEY")
	if apiKey == "" {
		log.Fatalln("BETASERIES_KEY must be set in env")
	}

	apiSecret := os.Getenv("BETASERIES_SECRET")
	if apiSecret == "" {
		log.Fatalln("BETASERIES_SECRET must be set in env")
	}

	// The quality can be specified using an environnement variable
	quality := os.Getenv("SHOWRSS_QUALITY")
	if quality == "" {
		quality = "720p"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "7777"
	}

	//The refreshInterval is expressed in minutes
	refreshInterval, err := strconv.Atoi(os.Getenv("SHOWRSS_REFRESH_TIME"))
	if refreshInterval == 0 || err != nil {
		refreshInterval = 60
	}

	limitPerShow, err := strconv.Atoi(os.Getenv("SHOWRSS_EP_PER_SHOW"))
	if limitPerShow == 0 || err != nil {
		limitPerShow = 100
	}

	//workaround for Heroku
	// must enable runtime-dyno-metadata
	//with heroku labs:enable runtime-dyno-metadata -a <app name>
	hostname := os.Getenv("HEROKU_APP_NAME")
	host := fmt.Sprintf("http://%s.herokuapp.com", hostname)

	if hostname == "" {
		hostname, _ = os.Hostname()
		host = fmt.Sprintf("http://%s:%s", hostname, port)
	}

	redirectURL, err := url.Parse(fmt.Sprintf("%s/auth_callback", host))
	if err != nil {
		log.Fatalf("Error parsing redirectURL : %s", err)
	}
	redirectURLString := redirectURL.String()

	// Configuration for the Oauth authentification with Betaseries
	conf := oauth2.Config{
		ClientID:     apiKey,
		ClientSecret: apiSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://www.betaseries.com/authorize",
			TokenURL: "https://api.betaseries.com/oauth/access_token",
		},
		RedirectURL: redirectURLString,
	}

	episodeProvider := betaseries.Betaseries{APIKey: apiKey, LimitPerShow: limitPerShow}

	log.Println("Starting SHOWRSS ...")
	log.Printf("Refresh interval set to %v minutes\n", refreshInterval)
	log.Printf("Betaserie episodes per show limit set to %v \n", limitPerShow)
	log.Printf("Default quality : %s", quality)

	//DB stuff
	store, err := dao.InitDB(*dbAddr)
	if err != nil {
		log.Fatalln("Error connecting to DB")
	}

	err = store.CreateBucket("episodes")
	if err != nil {
		log.Fatalln("Error when creating bucket")
	}

	// Worker stuff
	// A channel is used to pass the episode that we need to search
	episodeToSearch := make(chan dao.Episode, 1000)
	//searchWorker read the episode to search from the channel and if it found them save them in the db
	go searchWorker(episodeToSearch, store, quality)

	refreshLimiter := make(chan time.Time, 10)
	go func() {
		for t := range time.Tick(time.Duration(refreshInterval) * time.Hour) {
			refreshLimiter <- t
		}
	}()

	// we use a map to store the users because why not (we only store the token for each user si that we can refresh the unseen episodes from Betaseries)
	users := make(map[string]bool)
	newAuthChan := make(chan string, 10)
	go handleNewAuth(newAuthChan, users, refreshLimiter)

	go refresh(refreshLimiter, users, store, episodeProvider, episodeToSearch)

	errChan := make(chan error, 10)

	mux := mux.NewRouter()
	mux.HandleFunc("/", handlers.HelloHandler)
	mux.Handle("/auth", handlers.OauthHandler(conf))
	mux.Handle("/auth_callback", handlers.AuthCallbackHandler(conf, newAuthChan))
	mux.Handle("/episodes", handlers.EpisodeHandler(store))
	mux.Handle("/rss/{user}", handlers.RSSHandler(store, episodeProvider))

	httpServer := http.Server{}
	httpServer.Addr = ":" + port
	httpServer.Handler = handlers.LoggingHandler(mux)

	log.Printf("HTTP service listening on %s", host)

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
