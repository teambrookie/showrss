package torrent

import (
	"fmt"

	"github.com/qopher/go-torrentapi"
)

func goodEnoughTorrent(results torrentapi.TorrentResults) string {
	for _, t := range results {
		if t.Seeders > 0 || t.Leechers > 0 {
			return t.Download
		}
	}
	return ""
}

func Search(showID string, episodeCode string, quality string) (string, error) {
	api, err := torrentapi.New("SHOWRSS")
	if err != nil {
		return "", err
	}
	searchString := fmt.Sprintf("%s %s", episodeCode, quality)
	api.Format("json_extended")
	api.SearchTVDB(showID)
	api.SearchString(searchString)
	results, err := api.Search()
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}
	return goodEnoughTorrent(results), nil
}
