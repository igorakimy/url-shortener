package redirect_test

import (
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"

	"github.com/igorakimy/url-shortener/internal/lib/api"
	"github.com/igorakimy/url-shortener/internal/lib/logger/handlers/slogdiscard"
	"github.com/igorakimy/url-shortener/internal/server/handlers/url/redirect"
	"github.com/igorakimy/url-shortener/internal/server/handlers/url/redirect/mocks"
)

func TestRedirectHandler(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:  "Success",
			alias: "test_alias",
			url:   "https://www.google.com",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			urlGetterMock := mocks.NewURLGetter(t)

			if tc.respError == "" || tc.mockError != nil {
				urlGetterMock.EXPECT().
					GetURL(tc.alias).
					Return(tc.url, tc.mockError).
					Once()
			}

			r := chi.NewRouter()
			r.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetterMock))

			ts := httptest.NewServer(r)
			defer ts.Close()

			redirectedToURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)
			require.NoError(t, err)

			require.Equal(t, tc.url, redirectedToURL)
		})
	}
}
