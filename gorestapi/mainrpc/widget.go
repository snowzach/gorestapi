package mainrpc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/snowzach/queryp"

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
			server.RenderErrInvalidRequest(w, err)
			return
		}

		err := s.grStore.WidgetSave(ctx, widget)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpSave))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("WidgetSave error", "error", err, "error_id", errID)
			}
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
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "widget")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpGet))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("WidgetGetByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, widget)
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
		if err != nil {
			if err == store.ErrNotFound {
				server.RenderErrResourceNotFound(w, "widget")
			} else if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpDelete))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("WidgetDeleteByID error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderNoContent(w)
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

		qp, err := queryp.ParseRawQuery(r.URL.RawQuery)
		if err != nil {
			server.RenderErrInvalidRequest(w, err)
		}

		widgets, count, err := s.grStore.WidgetsFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				server.RenderErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpFind))
			} else {
				errID := server.RenderErrInternalWithID(w, nil)
				s.logger.Errorw("WidgetsFind error", "error", err, "error_id", errID)
			}
			return
		}

		server.RenderJSON(w, http.StatusOK, store.Results{Count: count, Results: widgets})
	}
}
