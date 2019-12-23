package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/goroute/logger"
	"github.com/goroute/route"
)

func main() {
	mux := route.NewServeMux()
	mux.Use(logger.New(logger.Format(logger.FormatTypeText)))

	mux.GET("/ok", func(c route.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	mux.GET("/cont", func(c route.Context) error {
		return c.String(http.StatusContinue, "continue")
	})

	mux.GET("/warn", func(c route.Context) error {
		return c.String(http.StatusNotFound, "not found")
	})

	mux.GET("/err", func(c route.Context) error {
		return errors.New("something went wrong")
	})

	log.Fatal(http.ListenAndServe(":9000", mux))
}
