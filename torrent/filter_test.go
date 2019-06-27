package torrent

import (
	"reflect"
	"testing"

	"github.com/qopher/go-torrentapi"
)

func TestInclude(t *testing.T) {
	var torrents = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "A B C"},
		torrentapi.TorrentResult{Download: "A B D"},
		torrentapi.TorrentResult{Download: "A C D"},
		torrentapi.TorrentResult{Download: "B C D"},
	}
	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "A B C"},
		torrentapi.TorrentResult{Download: "A B D"},
		torrentapi.TorrentResult{Download: "B C D"},
	}
	results := include("B", torrents)
	if !reflect.DeepEqual(results, expected) {
		t.Logf("%v", results)
		t.Logf("%v", expected)
		t.Fatalf("include fail")
	}
}

func TestExclude(t *testing.T) {
	var torrents = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "A B C"},
		torrentapi.TorrentResult{Download: "A B D"},
		torrentapi.TorrentResult{Download: "A C D"},
		torrentapi.TorrentResult{Download: "B C D"},
	}
	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "A C D"},
	}
	results := exclude("B", torrents)
	if !reflect.DeepEqual(results, expected) {
		t.Fatalf("exclude fail")
	}
}

func TestFilterCat(t *testing.T) {
	cat := Category{Type: "include", Optionnal: false, Keywords: []string{"mp4", "mkv", "avi"}}
	var torrents = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "torrent.mkv"},
		torrentapi.TorrentResult{Download: "anotherone.avi"},
		torrentapi.TorrentResult{Download: "yolo.mkv"},
	}

	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "torrent.mkv"},
		torrentapi.TorrentResult{Download: "yolo.mkv"},
	}

	results := filterCat(cat, torrents)
	if !reflect.DeepEqual(results, expected) {
		t.Log(results)
		t.Fatalf("filter cat fail")

	}
}

func TestFilterCatOptionnal(t *testing.T) {
	cat := Category{Type: "include", Optionnal: true, Keywords: []string{"mp3", "movie", "x265"}}
	var torrents = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "torrent.mkv"},
		torrentapi.TorrentResult{Download: "anotherone.avi"},
		torrentapi.TorrentResult{Download: "yolo.mkv"},
	}

	expected := torrents

	results := filterCat(cat, torrents)
	if !reflect.DeepEqual(results, expected) {
		t.Log(results)
		t.Fatalf("filter cat fail")

	}
}

func TestFilter(t *testing.T) {
	cats := []Category{}
	cats = append(cats, Category{Type: "include", Optionnal: false, Keywords: []string{"mkv", "mp4"}})
	cats = append(cats, Category{Type: "exclude", Optionnal: true, Keywords: []string{"serie"}})
	cats = append(cats, Category{Type: "include", Optionnal: false, Keywords: []string{"720p", "1080p", "4K"}})

	var torrents = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "serie.One.720p.mp4"},
		torrentapi.TorrentResult{Download: "serie.Dos.4K.mkv"},
		torrentapi.TorrentResult{Download: "serie.Trois.720p.mkv"},
		torrentapi.TorrentResult{Download: "serie.Four.1080p.mkv"},
		torrentapi.TorrentResult{Download: "serie.Cinco.720p.mkv"},
		torrentapi.TorrentResult{Download: "serie.Six.720p.mp4"},
	}

	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "serie.Trois.720p.mkv"},
		torrentapi.TorrentResult{Download: "serie.Cinco.720p.mkv"},
	}

	results := Filter(cats, torrents)
	if !reflect.DeepEqual(expected, results) {
		t.Log(results)
		t.Fatalf("Filter fail")
	}
}
