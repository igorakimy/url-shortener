package delete

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"github.com/igorakimy/url-shortener/internal/lib/api/response"
	"github.com/igorakimy/url-shortener/internal/lib/logger/sl"
)

// URLDeleter is interface for deleting a URL
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLDeleter --with-expecter=true
type URLDeleter interface {
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.delete.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "alias")
		if alias == "" {
			log.Info("alias is empty")
			render.JSON(w, r, response.Error("alias is required"))
			return
		}

		if err := urlDeleter.DeleteURL(alias); err != nil {
			log.Error("failed to delete url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
			return
		}

		log.Info("url deleted", slog.String("url", alias))

		render.JSON(w, r, response.OK())
	}
}
