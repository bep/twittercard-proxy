// Copyright © 2017-present Bjørn Erik Pedersen <bjorn.erik.pedersen@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var (
		httpAddr  string
		cardsFile string
	)

	flag.StringVar(&httpAddr, "http", "0.0.0.0:1414", "The HTTP listen address")
	flag.StringVar(&cardsFile, "f", "./twittercards.json", "The JSON filename with twitter cards")

	flag.Parse()

	p := newTcProxy(cardsFile)
	if err := p.load(); err != nil {
		p.log.Fatal(err)
	}

	server := http.Server{
		Addr:    httpAddr,
		Handler: p,
	}

	go func() {
		p.log.Fatal(server.ListenAndServe())
	}()

	p.log.Printf("HTTP listener on %s...", httpAddr)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	for {
		s := <-signalChan
		switch s {
		case syscall.SIGHUP:
			p.log.Printf("Captured SIGHUP, reloading twitter cards from %q...", p.cardsFile)
			if err := p.load(); err != nil {
				p.log.Println("ERROR: Failed to reload twitter cards:", err)
			}
		default:
			p.log.Printf("Captured %v. Exiting...", s)
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			server.Shutdown(shutdownCtx)

			<-shutdownCtx.Done()
			p.log.Println(shutdownCtx.Err())
			os.Exit(0)
		}
	}
}
