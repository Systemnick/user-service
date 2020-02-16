package main

import (
	"context"
	"fmt"

	"github.com/Systemnick/user-service/config"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type Application struct {
	config *config.Config
	logger *zerolog.Logger
	router *routing.Router
	server *fasthttp.Server
}

func NewApplication(c *config.Config, l *zerolog.Logger) (*Application, error) {
	a := &Application{}
	a.config = c
	a.logger = l

	router := routing.New()
	v1 := router.Group("/v1")
	v1.Use(authorization)
	v1.Get("/users", list).Post(h2)
	v1.Put("/users/<id>", h3).Delete(h4)

	router.Get("/v1", func(c *routing.Context) error {
		fmt.Fprintf(c, "Hello, world!\n")
		return nil
	})

	return a, nil
}

func (a *Application) Run() error {
	s := &fasthttp.Server{
		Handler:     a.router.HandleRequest,
		ReadTimeout: serviceStoppingTimeout,
	}
	a.server = s

	return s.ListenAndServe(":80")
}

func (a *Application) Stop(ctx context.Context) error {
	shutdownChan := make(chan bool)

	go func(ch chan bool) {
		err := a.server.Shutdown()
		if err != nil {
			a.logger.Error().Err(err).Msg("Server shutdown error")
		}
		ch <- true
	}(shutdownChan)

	select {
	case <-ctx.Done():
		a.logger.Info().Msg("Server shutdown by timeout")
	case <-shutdownChan:
		a.logger.Info().Msg("Server successfully shutdown")
	}
	return nil
}
