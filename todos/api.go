package todos

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func TodosRoute(s Service) chi.Router {
	r := chi.NewRouter()

	r.Get("/", getAllHandler(s))
	r.Post("/", createHandler(s))
	r.Put("/", updateHandler(s))

	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", getHandler(s))
		r.Delete("/", removeHandler(s))
	})

	return r
}

func noopHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func getAllHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		todos, err := s.GetAll(r.Context())
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, todos)
	}
}

func createHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := render.DecodeJSON(r.Body, &todo); err != nil {
			render.JSON(w, r, err)
			return
		}
		todo, err := s.Add(r.Context(), todo)
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, todo)
	}
}

func getHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		todo, err := s.Get(r.Context(), id)
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, todo)
	}
}


func removeHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		count, err := s.Remove(r.Context(), id)
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, count)
	}
}

func updateHandler(s Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var todo Todo
		if err := render.DecodeJSON(r.Body, &todo); err != nil {
			render.JSON(w, r, err)
			return
		}
		count, err := s.Update(r.Context(), todo)
		if err != nil {
			render.JSON(w, r, err)
			return
		}
		render.JSON(w, r, count)
	}
}
