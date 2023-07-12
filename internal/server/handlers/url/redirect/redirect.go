package redirect

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"golang.org/x/exp/slog"

	"github.com/igorakimy/url-shortener/internal/lib/api/response"
	"github.com/igorakimy/url-shortener/internal/lib/logger/sl"
	"github.com/igorakimy/url-shortener/internal/storage"
)

// URLGetter is an interface for getting a URL by alias
//
//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLGetter --with-expecter=true
type URLGetter interface {
	GetURL(alias string) (string, error)
}

// New is a http handler for redirecting from url to alias
func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.redirect.New"

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

		resURL, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not found", "alias", alias)
			render.JSON(w, r, response.Error("url not found"))
			return
		}
		if err != nil {
			log.Error("failed to get url", sl.Err(err))
			render.JSON(w, r, response.Error("internal error"))
		}

		log.Info("url found", slog.String("url", resURL))

		// redirect to found url
		http.Redirect(w, r, resURL, http.StatusFound)
	}
}
