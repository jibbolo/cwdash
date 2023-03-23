package app

import (
	"bytes"
	"embed"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jibbolo/cwdash/internal/manager/widget"
)

//go:embed all:templates
var templates embed.FS

const emptyPNG = `iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII=`

var (
	dashboardTmpl = template.Must(template.ParseFS(templates, "templates/dashboard.html"))
	indexTmpl     = template.Must(template.ParseFS(templates, "templates/index.html"))
)

func (a *App) indexFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := indexTmpl.Execute(w, struct {
		BUILD string
		List  []string
	}{a.version, a.dm.DashboardList()}); err != nil {
		log.Printf("can't render template: %v\n", err)
		http.Error(w, http.StatusText(500), 500)
		return
	}
}

func renderEmptyPNG(w http.ResponseWriter) {
	pngBody, _ := base64.RawStdEncoding.DecodeString(emptyPNG)
	w.Write(pngBody)
}

func (a *App) widgetFunc(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	number := chi.URLParam(r, "number")
	index, _ := strconv.Atoi(number)

	var buf bytes.Buffer
	clen, err := a.dm.RenderGraph(r.Context(), &buf, name, index)
	if err != nil {
		log.Printf("can't render widget: %v\n", err)
		renderEmptyPNG(w)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=300")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", clen))

	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("can't write buf to response writer: %v\n", err)
		return
	}
}

func (a *App) dashboardFunc(grid bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		name := chi.URLParam(r, "name")
		dashboard := a.dm.GetDashboard(name)

		w.Header().Set("Content-Type", "text/html")

		body := struct {
			Grid       bool
			Name       string
			Widgets    []widget.Widget
			LastUpdate time.Time
		}{grid, name, dashboard.Widgets, time.Now()}

		if err := dashboardTmpl.Execute(w, body); err != nil {
			log.Printf("can't render template: %v\n", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}
