package app

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sanchey92/jwt-example/internal/config"
)

type App struct {
	config     *config.Config
	httpServer *http.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run() error {
	fmt.Println("Server running on port: " + a.config.Port)

	if err := a.httpServer.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		//...
		a.initHTTPServer,
	}

	for _, fn := range inits {
		if err := fn(ctx); err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if a.config == nil {
		a.config = config.MustLoadConfig()
	}
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from server"))
	})

	a.httpServer = &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: r,
	}

	return nil
}
