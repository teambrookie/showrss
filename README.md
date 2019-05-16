# ShowRSS

[![Build Status](https://travis-ci.org/teambrookie/showrss.svg?branch=master)](https://travis-ci.org/teambrookie/showrss)
[![Go Report Card](https://goreportcard.com/badge/github.com/teambrookie/showrss)](https://goreportcard.com/report/github.com/teambrookie/showrss)

## Description

ShowRSS is a small app that let you ask Betaseries for your unseen episodes and then find the corresponding torrent using RARBG and expose them as an RSS feed.

### How to run it

```docker
docker run -p xxx:xxx -e PORT=xxx -e BETASERIES_KEY=xxx -e BETASERIES_SECRET=xxx teambrookie/showrss
```

### Additionnal environnement variables

- SHOWRSS_QUALITY=xxx (default is 720p)
- SHOWRSS_REFRESH_TIME=xxx (in minutes,default is 60)
- SHOWRSS_EP_PER_SHOW=xxx (default is 48)


### Using it

The authentification with Betaseries is using Oauth. You call the endpoint */auth*, it will redirect you to Betaseries to authentificate. And then if it's successful will redirect you to /rss/{user_token} so you can access your RSS feed.

The search for new unseen episodes and torrents is done automatically. It will start when you authentificate and then will occurs every 60 minutes ( if you haven't change the default)


### Endpoint

- /auth -> authentificate you with Betaseries
- /rss/{user_token} -> return your rss feed
- /info/{filename} -> return all the info given a filename

### Examples

```http
http://yourdomain.com/info/Modern.Family.S10E01.iNTERNAL.720p.WEB.x264-BAMBOOZLE
```

will return you the following
```json
{
  "name": "Modern.Family.S10E01",
  "season": 10,
  "episode": 1,
  "code": "S10E01",
  "show_id": 000,
  "magnet_link": "xxx",
  "filename": "Modern.Family.S10E01.iNTERNAL.720p.WEB.x264-BAMBOOZLE",
  "last_modified": "xxx"
}
```
