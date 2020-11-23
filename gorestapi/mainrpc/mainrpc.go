package mainrpc

import (
	"github.com/go-chi/chi"
	"github.com/snowzach/gorestapi/gorestapi"
	"go.uber.org/zap"
)

// Server is the API web server
type Server struct {
	logger  *zap.SugaredLogger
	router  chi.Router
	grStore gorestapi.GRStore
}

// Setup will setup the API listener
func Setup(router chi.Router, grStore gorestapi.GRStore) error {

	s := &Server{
		logger:  zap.S().With("package", "thingrpc"),
		router:  router,
		grStore: grStore,
	}

	// Base Functions
	s.router.Post("/things", s.ThingSave())
	s.router.Get("/things/{id}", s.ThingGetByID())
	s.router.Delete("/things/{id}", s.ThingDeleteByID())
	s.router.Get("/things", s.ThingsFind())

	s.router.Post("/widgets", s.WidgetSave())
	s.router.Get("/widgets/{id}", s.WidgetGetByID())
	s.router.Delete("/widgets/{id}", s.WidgetDeleteByID())
	s.router.Get("/widgets", s.WidgetsFind())

	return nil

}
