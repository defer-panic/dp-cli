package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

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

	if httpResp.StatusCode >= http.StatusBadRequest {
		return "", fmt.Errorf("failed to get login page: %s", httpResp.Status)
	}

	var resp struct {
		Link string `json:"link"`
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&resp); err != nil {
		return "", fmt.Errorf("failed to decode login page response: %w", err)
	}

	return resp.Link, nil
}

func (s *Service) SaveToken(token, server, configPath string) error {
	cfg := &config.Config{
		Server: server,
		JWT:  token,
	}
	
	if err := cfg.Save(configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	return nil
}
