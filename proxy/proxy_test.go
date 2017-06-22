// Copyright © 2017-present Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTcProxy(t *testing.T) {
	assert := require.New(t)

	p := NewTcProxy("mytweets.json")

	assert.NotNil(p)
	assert.NotNil(p.templ)
	assert.NotNil(p.Log)
	assert.Equal("mytweets.json", p.cardsFile)
}

func TestTcProxyLoad(t *testing.T) {
	assert := require.New(t)

	p := NewTcProxy("testdata/test.json")
	assert.NotNil(p)
	assert.NoError(p.Load())
	assert.Len(p.routes.Load().(routes), 2)
	tweet, found := p.getTweet("/c1")
	assert.True(found)
	assert.Equal("Twitter Card #1", tweet.Title)

	// Error cases
	p = NewTcProxy("testdata/notfound.json")
	assert.Error(p.Load())

	p = NewTcProxy("testdata/invalid.json")
	assert.Error(p.Load())
}

func TestTcProxyServeHTTP(t *testing.T) {
	assert := require.New(t)

	p := NewTcProxy("testdata/test.json")
	assert.NoError(p.Load())

	req, err := http.NewRequest("GET", "/c1", nil)
	assert.NoError(err)
	rr := httptest.NewRecorder()
	p.ServeHTTP(rr, req)
	assert.Equal(http.StatusMovedPermanently, rr.Code)

	// Try the same request as a "twitterBot"
	req.Header.Set("User-Agent", "I am a twitterBot!")
	rr = httptest.NewRecorder()
	p.ServeHTTP(rr, req)
	assert.Equal(http.StatusOK, rr.Code)
	body := rr.Body.String()
	assert.Contains(body, "<h1>Twitter Card #1</h1>")
	assert.Contains(body, `<meta name="twitter:site" content="@GoHugoIO">`)
	assert.Contains(body, `<meta name="twitter:creator" content="@bepsays">`)
	assert.Contains(body, `<meta name="twitter:title" content="Twitter Card #1">`)
	assert.Contains(body, `<meta name="twitter:description" content="Some text">`)
	assert.Contains(body, `<meta name="twitter:image" content="https://s11.postimg.org/jsd2aq1er/hugo-24-poster.png">`)

	// Error cases
	req, _ = http.NewRequest("POST", "/c1", nil)
	rr = httptest.NewRecorder()
	p.ServeHTTP(rr, req)
	assert.Equal(http.StatusMethodNotAllowed, rr.Code)

	req, _ = http.NewRequest("GET", "/notfound", nil)
	rr = httptest.NewRecorder()
	p.ServeHTTP(rr, req)
	assert.Equal(http.StatusNotFound, rr.Code)

	req, _ = http.NewRequest("GET", "", nil)
	rr = httptest.NewRecorder()
	p.ServeHTTP(rr, req)
	assert.Equal(http.StatusNotFound, rr.Code)

}
