package torrent

import torrentapi "github.com/qopher/go-torrentapi"
import "github.com/teambrookie/showrss/dao"
import "strconv"

func goodEnoughTorrent(results torrentapi.TorrentResults) string {
	for _, t := range results {
		if t.Seeders > 0 || t.Leechers > 0 {
			return t.Download
		}
	}
	return ""
}

func Search(episode dao.Episode) (string, error) {
	api, err := torrentapi.Init()
	if err != nil {
		return "", err
	}
	api.Format("json_extended")
	api.SearchTVDB(strconv.Itoa(episode.ShowID))
	api.SearchString(episode.Code + " 720p")
	results, err := api.Search()
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}
	return goodEnoughTorrent(results), nil
}
