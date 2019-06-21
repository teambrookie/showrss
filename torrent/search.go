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

func bestTorrent(torrents torrentapi.TorrentResults) torrentapi.TorrentResult {
	bt := torrentapi.TorrentResult{}
	for _, t := range torrents {
		if (bt == torrentapi.TorrentResult{}) {
			bt = t
			continue
		}
		if (t.Seeders / (1 + t.Leechers)) > (bt.Seeders / (1 + bt.Leechers)) {
			bt = t
		}
	}
	return bt
}

func filterOutDeadTorrents(torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var res torrentapi.TorrentResults
	for _, t := range torrents {
		if t.Seeders > 0 || t.Leechers > 0 {
			res = append(res, t)
		}
	}
	return res
}

func getFilename(magnetLink string) string {
	regex := "dn=([^&%]+)"
	r, _ := regexp.Compile(regex)
	return r.FindStringSubmatch(magnetLink)[1]
}

func searchEpisode(episode dao.Episode, config Config) (Torrent, error) {
	api, err := torrentapi.New("SHOWRSS")
	if err != nil {
		return Torrent{}, err
	}
	searchString := episode.Code
	api.Format("json_extended")
	api.SearchTVDB(strconv.Itoa(episode.ShowID))
	api.SearchString(searchString)
	torrents, err := api.Search()
	if err != nil {
		return Torrent{}, err
	}

	if len(torrents) == 0 {
		return Torrent{}, nil
	}
	torrents = filterOutDeadTorrents(torrents)
	torrents = Filter(config.Categories, torrents)
	tr := bestTorrent(torrents)
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

func searchSeason(episode dao.Episode, config Config) (Torrent, error) {
	api, err := torrentapi.New("SHOWRSS")
	if err != nil {
		return Torrent{}, err
	}
	searchString := fmt.Sprintf("S%02d", episode.Season)
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
	torrents = filterOutDeadTorrents(torrents)
	torrents = Filter(config.Categories, torrents)
	tr := bestTorrent(torrents)
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

func Search(episode dao.Episode, config Config) (Torrent, error) {

	log.Println("Searching Season")
	result, err := searchSeason(episode, config)
	if err != nil {
		return Torrent{}, nil
	}
	if (result != Torrent{}) {
		return result, nil
	}

	log.Println("Searching Episode")
	result, err = searchEpisode(episode, config)
	if err != nil {
		return Torrent{}, err
	}
	return result, err

}
