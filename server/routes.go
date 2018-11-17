package server

// SetupRoutes configures all the routes for this service
func (s *Server) SetupRoutes() {

	// Register our routes
	s.router.Get("/version", s.GetVersion())

	// Base Functions
	s.router.Get("/things", s.ThingFind())
	s.router.Get("/things/{id}", s.ThingGet())
	s.router.Post("/things", s.ThingSave())
	s.router.Delete("/things/{id}", s.ThingDelete())

}
