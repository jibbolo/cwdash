package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jibbolo/cwdash/internal/manager"
)

const loadTimeout = 5 * time.Second

type App struct {
	version string
	dm      *manager.DashboardManager
	router  *chi.Mux
}

func New(version string) (*App, error) {

	a := &App{
		version: version,
		dm:      manager.New(),
	}

	a.router = chi.NewRouter()
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middleware.Heartbeat("/health-check"))

	a.router.Get("/", a.indexFunc)
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}/{number:\\d+}", a.widgetFunc)
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}", a.dashboardFunc(false))
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}/grid", a.dashboardFunc(true))

	if err := a.init(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Handler() http.Handler {
	return a.router
}

func (a *App) init() error {
	ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
	defer cancel()

	err := a.dm.RefreshDashboards(ctx)
	if err != nil {
		return fmt.Errorf("can't refresh dashboard list: %w", err)
	}
	return nil
}
