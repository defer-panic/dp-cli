package url_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/defer-panic/dp-cli/internal/url"
	. "github.com/samber/mo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_Shorten(t *testing.T) {
	t.Run("returns shortened URL", func(t *testing.T) {
		t.Run("when no custom identifier passed", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(`{"short_url": "https://dfrp.cc/abc"}`))
					}),
				)
				svc = url.NewService(ts.URL, "")
			)

			shortURL, err := svc.Shorten(context.Background(), "https://example.com", None[string]())
			require.NoError(t, err)
			assert.Equal(t, "https://dfrp.cc/abc", shortURL)
		})

		t.Run("when custom identifier passed", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						var req struct {
							ShortURL   string `json:"short_url"`
							Identifier string `json:"identifier"`
						}

						require.NoError(t, json.NewDecoder(r.Body).Decode(&req))

						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(fmt.Sprintf(`{"short_url": "https://dfrp.cc/%s"}`, req.Identifier)))
					}),
				)
				svc = url.NewService(ts.URL, "")
			)

			shortURL, err := svc.Shorten(context.Background(), "https://example.com", Some("myAwesomeIdentifier"))
			require.NoError(t, err)
			assert.Equal(t, "https://dfrp.cc/myAwesomeIdentifier", shortURL)
		})
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when request fails", func(t *testing.T) {
			svc := url.NewService("http://localhost:12345", "")

			shortURL, err := svc.Shorten(context.Background(), "https://example.com", None[string]())
			require.Error(t, err)
			assert.Empty(t, shortURL)
		})

		t.Run("when authorization fails", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusUnauthorized)
						_, _ = w.Write([]byte(`{"error": "unable to authorize"}`))
					}),
				)
				svc = url.NewService(ts.URL, "")
			)

			shortURL, err := svc.Shorten(context.Background(), "https://example.com", Some("myAwesomeIdentifier"))
			require.ErrorIs(t, err, url.ErrUnauthorized)
			assert.Empty(t, shortURL)
		})

		t.Run("when identifier is already taken", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusConflict)
						_, _ = w.Write([]byte(`{"error": "identifier already exists"}`))
					}),
				)
				svc = url.NewService(ts.URL, "")
			)

			shortURL, err := svc.Shorten(context.Background(), "https://example.com", Some("myAwesomeIdentifier"))
			require.ErrorIs(t, err, url.ErrIdentifierExists)
			assert.Empty(t, shortURL)
		})
	})
}

func TestService_Stats(t *testing.T) {
	t.Run("returns shortening info", func(t *testing.T) {
		var (
			ts = httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"identifier": "abc", "created_by": "bob", "original_url": "https://example.com", "visits": 42}`))
				}),
			)
			svc = url.NewService(ts.URL, "")
		)

		info, err := svc.Stats(context.Background(), "https://dfrp.cc/abc")
		require.NoError(t, err)
		assert.Equal(
			t,
			&url.Stats{
				Identifier:  "abc",
				OriginalURL: "https://example.com",
				CreatedBy:   "bob",
				Visits:      42,
			},
			info,
		)
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when shortening not found", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusNotFound)
						_, _ = w.Write([]byte(`{"error": "shortening not found"}`))
					}),
				)
				svc = url.NewService(ts.URL, "")
			)

			info, err := svc.Stats(context.Background(), "https://dfrp.cc/abc")
			require.ErrorIs(t, err, url.ErrNotFound)
			assert.Nil(t, info)
		})
	})
}
