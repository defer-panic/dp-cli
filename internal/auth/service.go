package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/defer-panic/dp-cli/internal/config"
)

type Service struct {
	server string
}

func NewService(server string) *Service {
	return &Service{
		server: server,
	}
}

func (s *Service) GetAuthLink(ctx context.Context) (string, error) {
	loginPageURL, err := url.JoinPath(s.server, "/auth/oauth/github/link")
	if err != nil {
		return "", fmt.Errorf("failed to build login page URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, loginPageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to get login page: %w", err)
	}

	var resp struct {
		Link string `json:"link"`
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return "", fmt.Errorf("failed to decode login page response: %w", err)
	}

	return resp.Link, nil
}

func (s *Service) SaveToken(token, server string) error {
	cfg, err := config.Load(config.DefaultPath)
	if err != nil {
		if os.IsNotExist(err) {
			if err := os.MkdirAll(path.Dir(config.DefaultPath), 0755); err != nil {
				return fmt.Errorf("failed to create config directory: %w", err)
			}
		} else {
			return err
		}
	}

	cfg.JWT = token
	cfg.Server = server
	
	if err := cfg.Save(config.DefaultPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
