package torrent

import (
	"reflect"
	"testing"

	"github.com/qopher/go-torrentapi"
)

func TestExclude3DMovies(t *testing.T) {
	var movies = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "The.Mommy.3D.mkv"},
		torrentapi.TorrentResult{Download: "Peach.Perfect.mkv"},
	}

	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "Peach.Perfect.mkv"},
	}

	if !reflect.DeepEqual(exclude3DMovies(movies), expected) {
		t.Fatalf("Exclude3DMovies doesn't work")
	}

}

func TestExcludeNoSeeder(t *testing.T) {
	var movies = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "The.Mommy.3D.mkv", Seeders: 0},
		torrentapi.TorrentResult{Download: "Peach.Perfect.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "Training.Drogo.mkv", Seeders: 0},
	}

	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "Peach.Perfect.mkv", Seeders: 12},
	}

	if !reflect.DeepEqual(excludeNoSeeder(movies), expected) {
		t.Fatalf("ExcludeNoSeeder doesn't work")
	}
}

func TestFilterAudioQuality(t *testing.T) {
	var movies = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "The.Mommy.DTS-HD.mkv"},
		torrentapi.TorrentResult{Download: "Peach.Perfect.TrueHD.7.1Atmos.mkv"},
	}

	var expected = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "The.Mommy.DTS-HD.mkv"},
	}

	if !reflect.DeepEqual(filteraudioQuality("DTS-HD", movies), expected) {
		t.Fatalf("FilterAudioQuality doesn't work")
	}
}

func TestFilterMovies(t *testing.T) {
	var movies = torrentapi.TorrentResults{
		torrentapi.TorrentResult{Download: "The.Mommy.DTS-HD.mkv", Seeders: 4},
		torrentapi.TorrentResult{Download: "Lords.of.the.Pimp.TrueHD.7.1Atmos.mkv", Seeders: 2},
		torrentapi.TorrentResult{Download: "Lords.of.the.Rims.DTS-HD.mkv", Seeders: 1},
		torrentapi.TorrentResult{Download: "Lords.of.the.Dolphin.TrueHD.7.1Atmos.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "BaeWatch.DTS-HD.MA.7.1.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "Ex-Catina.TrueHD.7.1Atmos.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "King.Arthour.DTS-HD.MA.7.1.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "The.Square.TrueHD.7.1Atmos.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "Fist.Club.3D.DTS-HD.MA.7.1.EXTENDED.mkv", Seeders: 12},
		torrentapi.TorrentResult{Download: "Peach.Perfect.DTS.EXTENDED.mkv", Seeders: 0},
	}

	var expected = "BaeWatch.DTS-HD.MA.7.1.mkv"

	if filterMovies(movies) != expected {
		t.Fatalf("FilterMovies doesn't work")
	}
}
