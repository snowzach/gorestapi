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

// WidgetSave saves a widget
func (s *Server) WidgetSave() http.HandlerFunc {

	// swagger:operation POST /api/widgets WidgetSave
	//
	// Create/Save Widget
	//
	// Creates or saves a widget. Omit the ID to auto generate.
	// Pass an existing ID to update.
	//
	// ---
	// tags:
	// - WIDGETS
	// parameters:
	// - name: widget
	//   in: body
	//   description: Widget to Save/Update
	//   required: true
	//   type: object
	//   schema:
	//     "$ref": "#/definitions/gorestapi_WidgetExample"
	// responses:
	//   '200':
	//     description: User Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/gorestapi_Widget"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var widget = new(gorestapi.Widget)
		if err := render.DecodeJSON(r.Body, widget); err != nil {
			render.Render(w, r, server.ErrInvalidRequest(err))
			return
		}

		err := s.grStore.WidgetSave(ctx, widget)
		if err != nil {
			s.logger.Warnf("WidgetSave error: %v", err)
			render.Render(w, r, server.ErrInvalidRequest(fmt.Errorf("could not save widget: %v", err)))
			return
		}

		render.JSON(w, r, widget)
	}

}

// WidgetGetByID returns the widget
func (s *Server) WidgetGetByID() http.HandlerFunc {

	// swagger:operation GET /api/widgets/{id} WidgetGetByID
	//
	// Get a Widget
	//
	// Fetches a Widget
	//
	// ---
	// tags:
	// - WIDGETS
	// parameters:
	// - name: id
	//   in: path
	//   description: Widget ID to fetch
	//   type: string
	//   required: true
	// responses:
	//   '200':
	//     description: Widget Object
	//     type: object
	//     schema:
	//       "$ref": "#/definitions/gorestapi_Widget"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		widget, err := s.grStore.WidgetGetByID(ctx, id)
		if err == store.ErrNotFound {
			render.Render(w, r, server.ErrNotFound)
			return
		} else if err != nil {
			s.logger.Errorf("WidgetGetByID error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.JSON(w, r, widget)
	}

}

// WidgetDeleteByID deletes a widget
func (s *Server) WidgetDeleteByID() http.HandlerFunc {

	// swagger:operation DELETE /api/widgets/{id} WidgetDeleteByID
	//
	// Delete a Widget
	//
	// Deletes a Widget
	//
	// ---
	// tags:
	// - WIDGETS
	// parameters:
	// - name: id
	//   in: path
	//   description: Widget ID to delete
	//   type: string
	//   required: true
	// responses:
	//   '204':
	//     description: No Content
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.grStore.WidgetDeleteByID(ctx, id)
		if err == store.ErrNotFound {
			render.Render(w, r, server.ErrNotFound)
			return
		} else if err != nil {
			s.logger.Errorf("WidgetDeleteByID error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.NoContent(w, r)

	}

}

// WidgetsFind finds widgets
func (s *Server) WidgetsFind() http.HandlerFunc {

	// swagger:operation GET /api/widgets WidgetsFind
	//
	// Find Widgets
	//
	// Gets a list of widgets
	//
	// ---
	// tags:
	// - WIDGETS
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
	//     description: Widget Objects
	//     schema:
	//       type: array
	//       items:
	//         "$ref": "#/definitions/gorestapi_Widget"
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		fqp := store.ParseURLValuesToFindQueryParameters(r.URL.Query())

		widgets, count, err := s.grStore.WidgetsFind(ctx, fqp)
		if err != nil {
			s.logger.Errorf("WidgetsFind error: %v", err)
			render.Render(w, r, server.ErrInternal(nil))
			return
		}

		render.JSON(w, r, store.Results{Count: count, Results: widgets})

	}

}
