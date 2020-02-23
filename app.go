package main

import (
	"context"
	"time"

	"github.com/Systemnick/user-service/config"
	"github.com/Systemnick/user-service/middleware"
	"github.com/Systemnick/user-service/users/delivery/http"
	"github.com/Systemnick/user-service/users/repository/pgx-users"
	"github.com/Systemnick/user-service/users/use-case"
	"github.com/qiangxue/fasthttp-routing"
	"github.com/rs/zerolog"
	"github.com/valyala/fasthttp"
)

type Application struct {
	config       *config.Config
	logger       *zerolog.Logger
	router       *routing.Router
	server       *fasthttp.Server
	middleware   routing.Handler
	usersHandler *http.UsersHandler
}

func NewApplication(c *config.Config, l *zerolog.Logger) (*Application, error) {
	a := &Application{}
	a.config = c
	a.logger = l

	repo, err := pgx_users.New(c.DatabaseUrl, "public", "users", 60*time.Second)
	if err != nil {
		l.Fatal().Err(err).Msg("Failed to init users repository")
	}

	a.middleware = middleware.New()

	router := routing.New()
	v1 := router.Group("/v1")
	v1.Use(middleware.New())

	a.router = router

	uc := use_case.New(repo, 60*time.Second)
	a.usersHandler = http.NewHandler(v1, uc, a.middleware)

	return a, nil
}

func (a *Application) Run() error {
	s := &fasthttp.Server{
		Handler:     a.router.HandleRequest,
		ReadTimeout: serviceStoppingTimeout,
	}
	a.server = s

	return s.ListenAndServe(a.config.HttpEndpoint)
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
		a.logger.Info().Msg("Server shutdown gracefully")
	}
	return nil
}
