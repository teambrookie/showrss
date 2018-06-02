package worker

import (
	"strconv"
	"testing"
	"time"

	"github.com/teambrookie/MediaRSS/movierss/dao"
)

type FakeMovieStore struct {
	movies []dao.Movie
}

func (store *FakeMovieStore) GetMovie(id string) (dao.Movie, error) {
	for _, m := range store.movies {
		if strconv.Itoa(m.Ids.Trakt) == id {
			return m, nil
		}
	}
	return dao.Movie{}, nil
}

func (store *FakeMovieStore) DeleteMovie(id string) error {
	for i, m := range store.movies {
		if strconv.Itoa(m.Ids.Trakt) == id {
			store.movies = append(store.movies[:i], store.movies[i+1:]...)
			return nil
		}
	}
	return nil
}

func (store *FakeMovieStore) AddMovie(mov dao.Movie) error {
	return nil
}

func (store *FakeMovieStore) UpdateMovie(mov dao.Movie) error {
	return nil
}

func (store *FakeMovieStore) GetAllMovies() ([]dao.Movie, error) {
	return store.movies, nil
}

func (store *FakeMovieStore) GetAllFoundMovies() ([]dao.Movie, error) {
	return store.GetAllMovies()
}

func (store *FakeMovieStore) GetAllNotFoundMovies() ([]dao.Movie, error) {
	return store.GetAllMovies()
}

func TestCleanupLogic(t *testing.T) {

	var movies = []dao.Movie{
		dao.Movie{Ids: dao.Ids{Trakt: 0}, MagnetLink: "The.Mommy.mkv", LastModified: time.Now().Add(-time.Hour * 24 * 8), LastAccess: time.Now()},
		dao.Movie{Ids: dao.Ids{Trakt: 1}, MagnetLink: "BaeWatch.mkv", LastModified: time.Now().Add(-time.Hour * 24 * 8), LastAccess: time.Now().Add(-time.Hour * 72)},
		dao.Movie{Ids: dao.Ids{Trakt: 2}, MagnetLink: "", LastModified: time.Now(), LastAccess: time.Now()},
	}

	var fakeStore = FakeMovieStore{movies: movies}

	pipeline := make(chan dao.Movie, 10)

	clean(pipeline, &fakeStore)

	mov, _ := fakeStore.GetMovie("2")

	if (mov == dao.Movie{}) {
		t.Fatalf("Movie with empty magnet link should still be in the store")
	}

	mov, _ = fakeStore.GetMovie("1")

	if (mov != dao.Movie{}) {
		t.Fatalf("Movie last modified more than a week ago and not accessed in the last 12h should be deleted from the store")
	}

	mov = <-pipeline

	expected, _ := fakeStore.GetMovie("0")
	if mov != expected {
		t.Fatalf("Movie last modified more than a week ago and accessed in the last 12h should go in the pipeline")
	}

}
