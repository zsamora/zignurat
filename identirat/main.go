/*
*  Identirat
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
	identiratPort := ":" + utils.GetConfig("IDENTIRAT_PORT")
	identiratServer := &http.Server{
		Addr:         identiratPort,
		Handler:      IdentiratController(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		err := identiratServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("FATAL")
			log.Fatal(err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
