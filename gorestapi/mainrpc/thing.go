package mainrpc

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/server"
	"github.com/snowzach/gorestapi/store"
)

// ThingSave saves a thing
func (s *Server) ThingSave() http.HandlerFunc {

	// swagger:operation POST /api/things ThingSave
	//
	// Create/Save Thing
	//
	// Creates or saves a thing. Omit the ID to auto generate.
	// Pass an existing ID to update.
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: thing
	//   in: body
	//   description: Thing to Save/Update
	//   required: true
	//   type: object
	//   schema:
	//     "$ref": "#/definitions/gorestapi_ThingExample"
	// responses:
	//   '200':
	//     description: User Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/gorestapi_Thing"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var thing = new(gorestapi.Thing)
		if err := render.DecodeJSON(r.Body, thing); err != nil {
			render.Render(w, r, server.ErrInvalidRequest(err))
			return
		}

		err := s.mainStore.ThingSave(ctx, thing)
		if err != nil {
			s.logger.Warnf("ThingSave error: %v", err)
			render.Render(w, r, server.ErrInvalidRequest(fmt.Errorf("could not save thing: %v", err)))
			return
		}

		render.JSON(w, r, thing)
	}

}

// ThingGetByID returns the thing
func (s *Server) ThingGetByID() http.HandlerFunc {

	// swagger:operation GET /api/things/{id} ThingGetByID
	//
	// Get a Thing
	//
	// Fetches a Thing
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: id
	//   in: path
	//   description: Thing ID to fetch
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: Thing Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/gorestapi_Thing"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		thing, err := s.mainStore.ThingGetByID(ctx, id)
		if err == store.ErrNotFound {
			render.Render(w, r, server.ErrNotFound)
			return
		} else if err != nil {
			s.logger.Errorf("ThingGetByID error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.JSON(w, r, thing)
	}

}

// ThingDeleteByID deletes a thing
func (s *Server) ThingDeleteByID() http.HandlerFunc {

	// swagger:operation DELETE /api/things/{id} ThingDeleteByID
	//
	// Delete a Thing
	//
	// Deletes a Thing
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: id
	//   in: path
	//   description: Thing ID to delete
	//   type: string
	//   required: true
	// responses:
	//   '204':
	//     description: No Content
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.mainStore.ThingDeleteByID(ctx, id)
		if err == store.ErrNotFound {
			render.Render(w, r, server.ErrNotFound)
			return
		} else if err != nil {
			s.logger.Errorf("ThingDeleteByID error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.NoContent(w, r)

	}

}

// ThingsFind finds things
func (s *Server) ThingsFind() http.HandlerFunc {

	// swagger:operation GET /api/things ThingsFind
	//
	// Find Things
	//
	// Gets a list of things
	//
	// ---
	// tags:
	// - THINGS
	// parameters:
	// - name: limit
	//   in: query
	//   description: Number of records to return
	//   type: int
	//   required: false
	// - name: offset
	//   in: query
	//   description: Offset of records to return
	//   type: int
	//   required: false
	// - name: id
	//   in: query
	//   description: Filter id
	//   type: string
	//   required: false
	// - name: name
	//   in: query
	//   description: Filter name
	//   type: string
	//   required: false
	// responses:
	//   '200':
	//     description: Thing Objects
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/gorestapi_Thing"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		fqp := store.ParseURLValuesToFindQueryParameters(r.URL.Query())

		things, count, err := s.mainStore.ThingsFind(ctx, fqp)
		if err != nil {
			s.logger.Errorf("ThingsFind error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.JSON(w, r, store.Results{Count: count, Results: things})

	}

}
