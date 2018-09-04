package mediarss

import (
	"bytes"
	"html/template"
	"time"

	"github.com/gorilla/feeds"
	"github.com/teambrookie/MediaRSS/mediarss/db"
)

func rss(title string, medias []db.Media) (string, error) {
	feed := &feeds.Feed{
		Title:   title,
		Link:    &feeds.Link{Href: "https://github.com/teambrookie/mediarss"},
		Created: time.Now(),
	}
	for _, m := range medias {
		if m.Magnet == "" {
			continue
		}
		description, err := rssItemDescrition(m)
		if err != nil {
			return "", err
		}
		item := &feeds.Item{
			Title:       m.Name,
			Link:        &feeds.Link{Href: m.Magnet},
			Description: description,
			Created:     m.LastUpdate,
		}
		feed.Add(item)
	}
	return feed.ToRss()
}

func rssItemDescrition(media db.Media) (string, error) {
	t, err := template.ParseFiles("rss_item.tmpl")
	if err != nil {
		return "", err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, media); err != nil {
		return "", err
	}
	return tpl.String(), nil
}
