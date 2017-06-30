# ShowRSS

[![Build Status](https://travis-ci.org/teambrookie/showrss.svg?branch=master)](https://travis-ci.org/teambrookie/showrss)
[![Go Report Card](https://goreportcard.com/badge/github.com/teambrookie/showrss)](https://goreportcard.com/report/github.com/teambrookie/showrss)

## Description

ShowRSS is a small app that let you ask Betaseries for your unseen episodes and then find the corresponding torrent using RARBG and expose them as an RSS feed.

### Using it

First of all you need a Betaseries Token, you obtain it using the /auth endpoint like this
```
curl -X POST --data "username=xxx&password=xxx" http://localhost:8000/auth
```
Note: the username and password are send using x-www-urlencoded

Then they are the /refresh endpoint, it's role is to refresh the unseen episode and to refresh the torrent. They are use like this :
```
http://yourdomain/refresh?action=episode&token=xxx
http://yourdomain/refresh?action=torrent
```

And finally what really interest us is the /rss endpoint
```
http://yourdomain/rss?token=xxx
```

###Testing

```
curl -X POST --data "username=xxx&password=xxx" http://localhost:8000/auth
```

##Running

```
docker run -p 8000:8000 -e BETASERIES_KEY=xxx showrss
```
