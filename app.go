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
	return fasthttp.ListenAndServe(":80", a.router.HandleRequest)
}

func (a *Application) Stop(context context.Context) error {

	return nil
}
