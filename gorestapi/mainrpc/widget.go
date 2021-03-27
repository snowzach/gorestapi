package mainrpc

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/snowzach/queryp"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/server"
	"github.com/snowzach/gorestapi/store"
)

// WidgetSave saves a widget
//
// @ID WidgetSave
// @Tags widgets
// @Summary Save widget
// @Description Save a widget
// @Param widget body gorestapi.WidgetExample true "Widget"
// @Success 200 {object} gorestapi.Widget
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /widgets [post]
func (s *Server) WidgetSave() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var widget = new(gorestapi.Widget)
		if err := server.DecodeJSON(r.Body, widget); err != nil {
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

		server.RenderJSON(w, http.StatusOK, widget)
	}

}

// WidgetGetByID saves a widget
//
// @ID WidgetGetByID
// @Tags widgets
// @Summary Get widget
// @Description Get a widget
// @Param id path string true "ID"
// @Success 200 {object} gorestapi.Widget
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 404 {object} server.ErrResponse "Not Found"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /widgets/{id} [get]
func (s *Server) WidgetGetByID() http.HandlerFunc {
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

// WidgetDeleteByID saves a widget
//
// @ID WidgetDeleteByID
// @Tags widgets
// @Summary Delete widget
// @Description Delete a widget
// @Param id path string true "ID"
// @Success 204 "Success"
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 404 {object} server.ErrResponse "Not Found"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /widgets/{id} [delete]
func (s *Server) WidgetDeleteByID() http.HandlerFunc {
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

// WidgetsFind saves a widget
//
// @ID WidgetsFind
// @Tags widgets
// @Summary Find widgets
// @Description Find widgets
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param offset query int false "offset"
// @Param limit query int false "limit"
// @Param sort query string false "query"
// @Success 200 {array} gorestapi.Widget
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 500 {object} server.ErrResponse "Internal Error"
// @Router /widgets [get]
func (s *Server) WidgetsFind() http.HandlerFunc {
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
