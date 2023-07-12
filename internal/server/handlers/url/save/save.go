package save

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"

	"github.com/igorakimy/url-shortener/internal/lib/api/response"
	"github.com/igorakimy/url-shortener/internal/lib/logger/sl"
	"github.com/igorakimy/url-shortener/internal/lib/random"
	"github.com/igorakimy/url-shortener/internal/storage"
)

const (
	aliasLength = 6
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	response.Response
	Alias string `json:"alias,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.28.2 --name=URLSaver --with-expecter=true
type URLSaver interface {
	SaveURL(urlToSave string, alias string) (int64, error)
	GetURL(alias string) (string, error)
	DeleteURL(alias string) error
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"

		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", sl.Err(err))
			render.JSON(w, r, response.Error("failed to decode request"))
			return
		}

		log.Info("request body decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, response.ValidationError(validateErr))
			return
		}

		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLength)
		}

		id, err := urlSaver.SaveURL(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, response.Error("url already exists"))
				return
			}
			log.Error("failed to add url", slog.String("url", req.URL))
			render.JSON(w, r, response.Error("failed to add url"))
			return
		}

		log.Info("url successfully added", slog.Int64("id", id))

		render.JSON(w, r, Response{
			Response: response.OK(),
			Alias:    alias,
		})
	}
}
