# ShowRSS

[![Build Status](https://travis-ci.org/teambrookie/showrss.svg?branch=master)](https://travis-ci.org/teambrookie/showrss)
[![Go Report Card](https://goreportcard.com/badge/github.com/teambrookie/showrss)](https://goreportcard.com/report/github.com/teambrookie/showrss)

## Description

ShowRSS is a small API that let you ask Betaseries for your unseen episodes and then find the corresponding torrent using RARBG and expose them as an RSS feed.

### Using it

Pretty straight forward, first you need to authentificate to obtain an user token for the Betaseries API.
You can obtain one using the /auth endpoint.
```
curl -X POST --data "username=xxx&password=xxx" http://localhost:8000/auth
```
Note: the username and password are send using x-www-urlencoded

From there your username and token are saved in the database and you can start using the API.

The main endpoint is **/{user}/rss** that expose the rss feed corresponding of the torrent for your unseen episode.
```
http://localhost:8000/fabienfoerster/rss
```

Note : every hour ShowRSS will refresh your unseen episode and try to find the corresponding torrent. But if you like you can force the refresh using the **/refreshes** endpoint.
```
http://localhost:8000/refreshes
```

Note : if it tickles your fancy your can see all your unseen episodes using the **/{user}/episodes** endpoint

## Running

```
docker run -p 8000:8000 -e BETASERIES_KEY=xxx teambrookie/showrss
```
