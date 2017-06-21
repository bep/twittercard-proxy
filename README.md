# Twitter Card Proxy

## Install

**twittercard-proxy** is a Go application. The currently easiest way to intall it is via `go get`:

```bash
 go get -v github.com/bep/twittercard-proxy
```
 
## Use

```bash
Usage of twittercard-proxy:
  -f string
    	The JSON filename with twitter cards. (default "twittercards.json")
  -http string
    	The HTTP listen address. (default "0.0.0.0:1414")
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