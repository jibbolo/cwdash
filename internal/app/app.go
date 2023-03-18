package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/jibbolo/cwdash/internal/manager"
)

const loadTimeout = 5 * time.Second

type App struct {
	build, port string
	dm          *manager.DashboardManager
	router      *chi.Mux
}

func New(build, port string) *App {

	a := &App{
		build: build,
		port:  port,
		dm:    manager.New(),
	}
	a.router = chi.NewRouter()
	a.router.Use(middleware.Logger)
	a.router.Use(middleware.Recoverer)
	a.router.Use(middleware.Heartbeat("/health-check"))

	a.router.Get("/", a.indexFunc)
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}/{number:\\d+}", a.widgetFunc)
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}", a.dashboardFunc(false))
	a.router.Get("/dashboard/{name:[A-Za-z0-9\\-]+}/grid", a.dashboardFunc(true))

	return a
}

// Run initation the imageGenerator and starts the webserver
func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), loadTimeout)
	defer cancel()

	err := a.dm.RefreshDashboardList(ctx)
	if err != nil {
		return fmt.Errorf("can't refresh dashboard list: %w", err)
	}

	srv := &http.Server{
		Handler:      a.router,
		Addr:         ":" + a.port,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	log.Println("Listening... http://0.0.0.0:" + a.port)
	return srv.ListenAndServe()
}
