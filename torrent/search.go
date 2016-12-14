package torrent

import torrentapi "github.com/qopher/go-torrentapi"

func Search(name string) (string, error) {
	searchName := name + " 720p"
	api, err := torrentapi.Init()
	if err != nil {
		return "", err
	}
	api.Format("json_extended")
	api.SearchString(searchName)
	results, err := api.Search()
	if err != nil {
		return "", nil
	}
	if len(results) == 0 {
		return "", nil
	}
	return results[0].Download, nil
}
