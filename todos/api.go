package todos

import (
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func Route(s Service) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	r.Route("/todos", func(r chi.Router) {
		r.Get("/", getAllHandler(s))
		r.Post("/", createHandler(s))
		r.Put("/", updateHandler(s))

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", getHandler(s))
			r.Delete("/", removeHandler(s))
		})
	})

	return r
}

type ErrorMessage struct {
	Message string `json:"message"`
}

type UpdateMessage struct {
	Updated int `json:"updated"`
}

func getAllHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		todos, err := s.GetAll(ctx)
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		render.JSON(w, r, todos)
	}
}

func createHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var todo Todo
		if err := render.DecodeJSON(r.Body, &todo); err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		todo, err := s.Add(ctx, todo)
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		render.JSON(w, r, todo)
	}
}

func getHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		todo, err := s.Get(ctx, id)
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		render.JSON(w, r, todo)
	}
}

func removeHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		updated, err := s.Remove(ctx, id)
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		render.JSON(w, r, UpdateMessage{Updated: updated})
	}
}

func updateHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		var todo Todo
		if err := render.DecodeJSON(r.Body, &todo); err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		updated, err := s.Update(ctx, todo)
		if err != nil {
			log.Error().
				Stack().
				Err(err).
				Str("request_id", middleware.GetReqID(ctx)).
				Msg("Internal error")
			render.JSON(w, r, ErrorMessage{Message: err.Error()})
			return
		}

		render.JSON(w, r, UpdateMessage{Updated: updated})
	}
}
