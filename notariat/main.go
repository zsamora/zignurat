/*
*  Notariat
*  Author: @zsamora
 */
package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/zsamora/utils"
)

var g errgroup.Group

func main() {
	notariatPort := ":" + utils.GetConfig("NOTARIAT_PORT")
	notariatServer := &http.Server{
		Addr:         notariatPort,
		Handler:      NotariatController(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := notariatServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
