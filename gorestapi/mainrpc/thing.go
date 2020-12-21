package mainrpc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/snowzach/queryp"

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
		if err := server.DecodeJSON(r.Body, thing); err != nil {
			server.RenderErrInvalidRequest(w, err)
			return
		}

		err := s.grStore.ThingSave(ctx, thing)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpSave))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("ThingSave error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, thing)
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

		thing, err := s.grStore.ThingGetByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "thing")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpGet))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("ThingGetByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, thing)
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

		err := s.grStore.ThingDeleteByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "thing")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpDelete))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("ThingDeleteByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderNoContent(w)

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

		qp, err := queryp.ParseRawQuery(r.URL.RawQuery)
		if err != nil {
			server.RenderErrInvalidRequest(w, err)
		}

		things, count, err := s.grStore.ThingsFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpFind))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("ThingsFind error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, store.Results{Count: count, Results: things})

	}

}
