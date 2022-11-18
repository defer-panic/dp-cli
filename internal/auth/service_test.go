package auth_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/defer-panic/dp-cli/internal/auth"
	"github.com/defer-panic/dp-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestService_GetAuthLink(t *testing.T) {
	t.Run("returns auth page link", func(t *testing.T) {
		var (
			ts = httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write([]byte(`{"link": "http://example.com"}`))
				}),
			)
			svc = auth.NewService(ts.URL)
		)

		link, err := svc.GetAuthLink(context.Background())
		require.NoError(t, err)
		assert.Equal(t, "http://example.com", link)
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when request fails", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusBadRequest)
						_, _ = w.Write([]byte(`{"message": "invalid request"}`))
					}),
				)
				svc = auth.NewService(ts.URL)
			)

			link, err := svc.GetAuthLink(context.Background())
			require.Error(t, err)
			assert.Equal(t, "", link)
		})

		t.Run("when server address is invalid", func(t *testing.T) {
			svc := auth.NewService("://")
			link, err := svc.GetAuthLink(context.Background())
			require.Error(t, err)
			assert.Equal(t, "", link)
		})

		t.Run("when response is invalid", func(t *testing.T) {
			var (
				ts = httptest.NewServer(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte(`}`))
					}),
				)
				svc = auth.NewService(ts.URL)
			)

			link, err := svc.GetAuthLink(context.Background())
			require.Error(t, err)
			assert.Equal(t, "", link)
		})
	})
}

func TestService_SaveToken(t *testing.T) {
	t.Run("saves token", func(t *testing.T) {
		t.Run("when config file exists", func(t *testing.T) {
			svc := auth.NewService("http://example.com")

			err := svc.SaveToken("token", "https://dfrp.cc", "testdata/config.json")
			require.NoError(t, err)

			cfg, err := config.Load("testdata/config.json")
			require.NoError(t, err)

			assert.Equal(t, "token", cfg.JWT)
			assert.Equal(t, "https://dfrp.cc", cfg.Server)
		})

		t.Run("when config file does not exist", func(t *testing.T) {
			t.Cleanup(func() {
				os.Remove("testdata/config_new.json")
			})

			svc := auth.NewService("http://example.com")

			err := svc.SaveToken("token", "https://dfrp.cc", "testdata/config_new.json")
			require.NoError(t, err)

			cfg, err := config.Load("testdata/config_new.json")
			require.NoError(t, err)

			assert.Equal(t, "token", cfg.JWT)
			assert.Equal(t, "https://dfrp.cc", cfg.Server)
		})
	})

	t.Run("returns error", func(t *testing.T) {
		t.Run("when there are not enough permissions", func(t *testing.T) {
			svc := auth.NewService("http://example.com")

			err := svc.SaveToken("token", "https://dfrp.cc", "/config.json")
			require.Error(t, err)
		})
	})
}
