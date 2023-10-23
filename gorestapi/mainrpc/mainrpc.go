package mainrpc

import (
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/snowzach/golib/log"

	"github.com/snowzach/gorestapi/gorestapi"
)

// Server is the API web server
type Server struct {
	logger  *slog.Logger
	router  chi.Router
	grStore gorestapi.GRStore
}

// Setup will setup the API listener
func Setup(router chi.Router, grStore gorestapi.GRStore) error {

	s := &Server{
		logger:  log.Logger.With("context", "mainrpc"),
		router:  router,
		grStore: grStore,
	}

	// Base Functions
	s.router.Route("/api", func(r chi.Router) {
		r.Post("/things", s.ThingSave())
		r.Get("/things/{id}", s.ThingGetByID())
		r.Delete("/things/{id}", s.ThingDeleteByID())
		r.Get("/things", s.ThingsFind())

		r.Post("/widgets", s.WidgetSave())
		r.Get("/widgets/{id}", s.WidgetGetByID())
		r.Delete("/widgets/{id}", s.WidgetDeleteByID())
		r.Get("/widgets", s.WidgetsFind())
	})

	return nil

}
