// Copyright © 2017-present Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"sync/atomic"
)

// TcProxy is the proxy http server.
type TcProxy struct {
	routes    atomic.Value
	Log       *log.Logger
	templ     *template.Template
	cardsFile string
}

// NewTcProxy creates a new twitter card proxy.
func NewTcProxy(cardsFile string) *TcProxy {
	templ := template.Must(template.New("").Parse(pageTemplate))
	return &TcProxy{
		templ:     templ,
		cardsFile: cardsFile,
		Log:       log.New(os.Stderr, "", log.LstdFlags)}
}

func (p *TcProxy) getTweet(path string) (twitterCard, bool) {
	r := p.routes.Load().(routes)
	t, found := r[path]
	return t, found
}

func (p *TcProxy) Load() error {
	tweets, err := readTwitterCards(p.cardsFile)
	if err != nil {
		return err
	}

	r := make(routes)
	for _, tweet := range tweets {
		tweetPath := "/" + tweet.ID
		r[tweetPath] = tweet
	}
	p.routes.Store(r)

	return nil
}

func (p *TcProxy) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(w, "Not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := cleanPath(req.URL.Path)

	tweet, found := p.getTweet(path)

	if !found {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	ua := strings.ToLower(req.UserAgent())

	// Show the Twitter Card to the Twittter Bot only!
	if !strings.Contains(ua, "twitterbot") {
		http.Redirect(w, req, tweet.Target, 301)
		return
	}

	if err := p.templ.Execute(w, tweet); err != nil {
		p.Log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type routes map[string]twitterCard

// Borrowed from the net/http package.
func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)
	// path.Clean removes trailing slash except for root;
	// put the trailing slash back if necessary.
	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}

func readTwitterCards(filename string) ([]twitterCard, error) {
	var tweets []twitterCard

	log.Printf("Read twitter cards from %q", filename)

	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return tweets, err
	}

	if err := json.Unmarshal(f, &tweets); err != nil {
		return tweets, err
	}

	return tweets, nil
}

type twitterCard struct {
	ID          string `json:"id"`
	Site        string `json:"site"`
	Creator     string `json:"creator"`
	Image       string `json:"image"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Target      string `json:"target"`
}

const pageTemplate = `
<!doctype html>
<html>
    <head>
		<link rel="stylesheet" href="//cdn.rawgit.com/milligram/milligram/master/dist/milligram.min.css">
		<meta charset="utf-8"/>
		<meta name="twitter:card" content="summary_large_image">
		<meta name="twitter:site" content="{{ .Site }}">
		<meta name="twitter:creator" content="{{ .Creator }}">
		<meta name="twitter:title" content="{{ .Title }}">
		<meta name="twitter:description" content="{{ .Description }}">
		<meta name="twitter:image" content="{{ .Image }}">
        <title>{{ .Title }}</title>
    </head>
    <body>
		<div class="container">
			<div class="row">
				<div class="column column-50">
				<h1>{{ .Title }}</h1>
				<h2><a class="button" href="{{ .Target }}">Go!</a></h2>
				<a href="{{ .Target }}">
					<img src="{{ .Image }}" width="300px">
				</a>
				</div>
			</div>
    </body>
</html>
`
