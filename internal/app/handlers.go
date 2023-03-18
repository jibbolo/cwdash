package app

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/jibbolo/cwdash/internal/widget"
)

const emptyPNG = `iVBORw0KGgoAAAANSUhEUgAAAAEAAAABAQMAAAAl21bKAAAAA1BMVEUAAACnej3aAAAAAXRSTlMAQObYZgAAAApJREFUCNdjYAAAAAIAAeIhvDMAAAAASUVORK5CYII=`

var (
	dashboardTmpl = template.Must(template.ParseFiles("templates/dashboard.html"))
	indexTmpl     = template.Must(template.ParseFiles("templates/index.html"))
)

func (a *App) indexFunc(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	if err := indexTmpl.Execute(w, struct {
		BUILD string
		List  []string
	}{a.build, a.dm.DashboardList()}); err != nil {
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
	clen, err := a.dm.RenderGraph(&buf, name, index)
	if err != nil {
		log.Printf("can't render widget: %v\n", err)
		renderEmptyPNG(w)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=300") // 30 days
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
		if dashboard == nil {
			if err := a.dm.RefreshBody(name); err != nil {
				http.Error(w, http.StatusText(404), 404)
				return
			}
			dashboard = a.dm.GetDashboard(name)
		}
		tmpl := dashboardTmpl
		w.Header().Set("Content-Type", "text/html")

		body := struct {
			Grid       bool
			Name       string
			Widgets    []widget.Widget
			LastUpdate time.Time
		}{grid, name, dashboard.Widgets, time.Now()}

		if err := tmpl.Execute(w, body); err != nil {
			log.Printf("can't render template: %v\n", err)
			http.Error(w, http.StatusText(500), 500)
			return
		}
	}
}
