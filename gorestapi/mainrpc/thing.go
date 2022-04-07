package mainrpc

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/snowzach/queryp"

	"github.com/snowzach/gorestapi/gorestapi"
	"github.com/snowzach/gorestapi/pkg/server/render"
	"github.com/snowzach/gorestapi/store"
)

// ThingSave saves a thing
//
// @ID ThingSave
// @Tags things
// @Summary Save thing
// @Description Save a thing
// @Param thing body gorestapi.ThingExample true "Thing"
// @Success 200 {object} gorestapi.Thing
// @Failure 400 {object} render.ErrResponse "Invalid Argument"
// @Failure 500 {object} render.ErrResponse "Internal Error"
// @Router /things [post]
func (s *Server) ThingSave() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		var thing = new(gorestapi.Thing)
		if err := render.DecodeJSON(r.Body, thing); err != nil {
			render.ErrInvalidRequest(w, err)
			return
		}

		err := s.grStore.ThingSave(ctx, thing)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				render.ErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpSave))
			} else {
				requestID := middleware.GetReqID(ctx)
				render.ErrInternalWithRequestID(w, requestID, nil)
				s.logger.Errorw("ThingSave error", "error", err, "request_id", requestID)
			}
			return
		}

		render.JSON(w, http.StatusOK, thing)
	}

}

// ThingGetByID saves a thing
//
// @ID ThingGetByID
// @Tags things
// @Summary Get thing
// @Description Get a thing
// @Param id path string true "ID"
// @Success 200 {object} gorestapi.Thing
// @Failure 400 {object} render.ErrResponse "Invalid Argument"
// @Failure 404 {object} render.ErrResponse "Not Found"
// @Failure 500 {object} render.ErrResponse "Internal Error"
// @Router /things/{id} [get]
func (s *Server) ThingGetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		thing, err := s.grStore.ThingGetByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				render.ErrResourceNotFound(w, "thing")
			} else if serr, ok := err.(*store.Error); ok {
				render.ErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpGet))
			} else {
				requestID := middleware.GetReqID(ctx)
				render.ErrInternalWithRequestID(w, requestID, nil)
				s.logger.Errorw("ThingGetByID error", "error", err, "request_id", requestID)
			}
			return
		}

		render.JSON(w, http.StatusOK, thing)
	}

}

// ThingDeleteByID saves a thing
//
// @ID ThingDeleteByID
// @Tags things
// @Summary Delete thing
// @Description Delete a thing
// @Param id path string true "ID"
// @Success 204 "Success"
// @Failure 400 {object} render.ErrResponse "Invalid Argument"
// @Failure 404 {object} render.ErrResponse "Not Found"
// @Failure 500 {object} render.ErrResponse "Internal Error"
// @Router /things/{id} [delete]
func (s *Server) ThingDeleteByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		id := chi.URLParam(r, "id")

		err := s.grStore.ThingDeleteByID(ctx, id)
		if err != nil {
			if err == store.ErrNotFound {
				render.ErrResourceNotFound(w, "thing")
			} else if serr, ok := err.(*store.Error); ok {
				render.ErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpDelete))
			} else {
				requestID := middleware.GetReqID(ctx)
				render.ErrInternalWithRequestID(w, requestID, nil)
				s.logger.Errorw("ThingDeleteByID error", "error", err, "request_id", requestID)
			}
			return
		}

		render.NoContent(w)

	}

}

// ThingsFind saves a thing
//
// @ID ThingsFind
// @Tags things
// @Summary Find things
// @Description Find things
// @Param id query string false "id"
// @Param name query string false "name"
// @Param description query string false "description"
// @Param offset query int false "offset"
// @Param limit query int false "limit"
// @Param sort query string false "query"
// @Success 200 {array} gorestapi.Thing
// @Failure 400 {object} render.ErrResponse "Invalid Argument"
// @Failure 500 {object} render.ErrResponse "Internal Error"
// @Router /things [get]
func (s *Server) ThingsFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		qp, err := queryp.ParseRawQuery(r.URL.RawQuery)
		if err != nil {
			render.ErrInvalidRequest(w, err)
		}

		things, count, err := s.grStore.ThingsFind(ctx, qp)
		if err != nil {
			if serr, ok := err.(*store.Error); ok {
				render.ErrInvalidRequest(w, serr.ErrorForOp(store.ErrorOpFind))
			} else {
				requestID := middleware.GetReqID(ctx)
				render.ErrInternalWithRequestID(w, requestID, nil)
				s.logger.Errorw("ThingsFind error", "error", err, "request_id", requestID)
			}
			return
		}

		render.JSON(w, http.StatusOK, store.Results{Count: count, Results: things})

	}

}
