package torrent

import (
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/qopher/go-torrentapi"
	"github.com/teambrookie/showrss/dao"
)

type Torrent struct {
	Name        string
	Filename    string
	MagnetLink  string
	TorrentURL  string
	TorrentType string
	ShowID      int
	Season      int
}

func goodEnoughTorrent(results torrentapi.TorrentResults) torrentapi.TorrentResult {
	for _, t := range results {
		if t.Seeders > 0 || t.Leechers > 0 {
			return t
		}
	}
	return torrentapi.TorrentResult{}
}

func getFilename(magnetLink string) string {
	regex := "dn=([^&%]+)"
	r, _ := regexp.Compile(regex)
	return r.FindStringSubmatch(magnetLink)[1]
}

func searchEpisode(episode dao.Episode, quality string) (Torrent, error) {
	api, err := torrentapi.New("SHOWRSS")
	if err != nil {
		return Torrent{}, err
	}
	searchString := fmt.Sprintf("%s %s", episode.Code, quality)
	api.Format("json_extended")
	api.SearchTVDB(strconv.Itoa(episode.ShowID))
	api.SearchString(searchString)
	results, err := api.Search()
	if err != nil {
		return Torrent{}, err
	}

	if len(results) == 0 {
		return Torrent{}, nil
	}
	tr := goodEnoughTorrent(results)
	torrent := Torrent{}
	if (tr != torrentapi.TorrentResult{}) {
		torrent.Name = tr.Title
		torrent.Filename = getFilename(tr.Download)
		torrent.MagnetLink = tr.Download
		torrent.TorrentURL = fmt.Sprintf("http://itorrents.org/torrent/%s.torrent", extractHashFromMagnet(tr.Download))
		torrent.TorrentType = "episode"
		torrent.ShowID = episode.ShowID
		torrent.Season = episode.Season
	}
	return torrent, nil

}

func extractHashFromMagnet(magnet string) string {
	r, _ := regexp.Compile("urn:btih:([^&]+)")
	return strings.ToUpper(r.FindStringSubmatch(magnet)[1])
}

func matchEpisode(name string) bool {
	regex := ".+S[0-9]+E.+"
	match, _ := regexp.MatchString(regex, name)
	return match
}

func filterFullSeasonTorrent(torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var fullSeasonTorrents torrentapi.TorrentResults
	for _, t := range torrents {
		if !matchEpisode(t.Title) {
			fullSeasonTorrents = append(fullSeasonTorrents, t)
		}
	}
	return fullSeasonTorrents
}

func searchSeason(episode dao.Episode, quality string) (Torrent, error) {
	api, err := torrentapi.New("SHOWRSS")
	if err != nil {
		return Torrent{}, err
	}
	searchString := fmt.Sprintf("S%02d %s", episode.Season, quality)
	api.Format("json_extended")
	api.SearchTVDB(strconv.Itoa(episode.ShowID))
	api.SearchString(searchString)
	results, err := api.Search()
	if err != nil {
		return Torrent{}, err
	}
	if len(results) == 0 {
		return Torrent{}, nil
	}

	torrents := filterFullSeasonTorrent(results)
	tr := goodEnoughTorrent(torrents)
	torrent := Torrent{}
	if (tr != torrentapi.TorrentResult{}) {
		torrent.Name = tr.Title
		torrent.Filename = getFilename(tr.Download)
		torrent.MagnetLink = tr.Download
		torrent.TorrentURL = fmt.Sprintf("http://itorrents.org/torrent/%s.torrent", extractHashFromMagnet(tr.Download))
		torrent.TorrentType = "season"
		torrent.ShowID = episode.ShowID
		torrent.Season = episode.Season
	}
	return torrent, nil

}

func Search(episode dao.Episode, quality string) (Torrent, error) {

	log.Println("Searching Season")
	result, err := searchSeason(episode, quality)
	if err != nil {
		return Torrent{}, nil
	}
	if (result != Torrent{}) {
		return result, nil
	}

	log.Println("Searching Episode")
	result, err = searchEpisode(episode, quality)
	if err != nil {
		return Torrent{}, err
	}
	return result, err

}
