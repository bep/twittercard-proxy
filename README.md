# Twitter Card Proxy

[![*Nix Build Status](https://travis-ci.org/bep/twittercard-proxy.svg)](https://travis-ci.org/bep/twittercard-proxy)
[![CircleCI](https://circleci.com/gh/bep/twittercard-proxy.svg?style=svg)](https://circleci.com/gh/bep/twittercard-proxy)
[![Windows Build status](https://ci.appveyor.com/api/projects/status/v9bbybn1n6y2k5xc?svg=true)](https://ci.appveyor.com/project/bep/twittercard-proxy)
[![GoDoc](https://godoc.org/github.com/bep/twittercard-proxy?status.svg)](https://godoc.org/github.com/bep/twittercard-proxy)
[![Go Report Card](https://goreportcard.com/badge/github.com/bep/twittercard-proxy)](https://goreportcard.com/report/github.com/bep/twittercard-proxy)

![Hugo 0.24 Tweet](https://s9.postimg.org/hvepyc1vz/hugo-024-tweet.png "Hugo 0.24 Tweet")

The above screenshot is the single motivation behind this tool: To get [nicer looking](https://twitter.com/GoHugoIO/status/877500564405444608) Twitter cards when announcing GitHub-releases. GitHub provides some Twitter metadata, but it is for your account only, and just linking to the release page gets you a [rather blend](https://twitter.com/GoHugoIO/status/875629224228306944) Twitter card.

## Install

Download a binary from [releases](https://github.com/bep/twittercard-proxy/releases).

**twittercard-proxy** is a Go application, so you can also easiliy compile it yourself or intall it is via `go get`:

```bash
 go get -v github.com/bep/twittercard-proxy
```
 
## Use

```bash
▶ ./twittercard-proxy -h
Usage of ./twittercard-proxy:
  -f string
    	The JSON filename with twitter cards (default "./twittercards.json")
  -http string
    	The HTTP listen address (default "0.0.0.0:1414"
```

To add or modify Twitter cards, just edit the `JSON` file and send a `SIGHUP` signal to the server process:

```bash
kill -s SIGHUP <process-id>
```

## File format

The `id` below will become the path and `target` is the proxy target. 

The rest are Twitter Card properties, see [https://dev.twitter.com/cards/types/summary-large-image](https://dev.twitter.com/cards/types/summary-large-image).

```json
[
	{
		"id": "hugo-0.24",
		"site": "@GoHugoIO",
		"creator": "@bepsays",
		"image": "https://s11.postimg.org/jsd2aq1er/hugo-24-poster.png",
		"title": "The Revival of the Archetypes!",
		"description": "In the new Hugo 0.24, archetype files are Go templates with all funcs and the full .Site available, for all content formats.",
		"target": "https://github.com/gohugoio/hugo/releases/tag/v0.24"
	}
]
```
