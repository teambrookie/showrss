package torrent

import (
	"strings"

	"github.com/qopher/go-torrentapi"
)

type Category struct {
	Type      string   `json:"type"`
	Optionnal bool     `json:"optionnal"`
	Keywords  []string `json:"keywords"`
}

type Config struct {
	Categories []Category
}

func chooseFilterFunc(catType string) func(string, torrentapi.TorrentResults) torrentapi.TorrentResults {
	if catType == "include" {
		return include
	}
	return exclude

}

func filterCat(category Category, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	results := torrents
	filter := chooseFilterFunc(category.Type)

	for _, keyword := range category.Keywords {
		results = filter(keyword, torrents)
		if results != nil {
			break
		}
	}
	if category.Optionnal && (results == nil) {
		return torrents
	}
	return results

}

func Filter(categories []Category, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	results := torrents
	for _, cat := range categories {
		tmp := filterCat(cat, results)
		results = tmp
	}
	return results
}

func include(keyword string, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var results torrentapi.TorrentResults
	keyword = strings.ToLower(keyword)
	for _, t := range torrents {
		var filename = strings.ToLower(t.Download)
		if strings.Contains(filename, keyword) {
			results = append(results, t)
		}
	}
	return results
}

func exclude(keyword string, torrents torrentapi.TorrentResults) torrentapi.TorrentResults {
	var results torrentapi.TorrentResults
	keyword = strings.ToLower(keyword)
	for _, t := range torrents {
		var filename = strings.ToLower(t.Download)
		if !strings.Contains(filename, keyword) {
			results = append(results, t)
		}
	}
	return results
}
