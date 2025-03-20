package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"github.com/sanchey92/jwt-example/internal/config"
	"github.com/sanchey92/jwt-example/internal/logger"
	"github.com/sanchey92/jwt-example/pkg/closer"
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

func (a *App) Run() {
	log := logger.GetLogger()

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal("server failed", zap.Error(err))
		}
	}()

	log.Info("Server started", zap.String("port", a.config.Port))

	closer.Add(func() error {
		log.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := a.httpServer.Shutdown(ctx); err != nil {
			return err
		}

		log.Info("Server Stopped")
		return nil
	})

	closer.Wait()
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
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
	a.config = config.MustLoadConfig()
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init()
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from server"))
	})

	a.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Port),
		Handler: r,
	}

	return nil
}
