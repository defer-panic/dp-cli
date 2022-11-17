package url

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	. "github.com/samber/mo"
)

type Service struct {
	server string
	token  string
}

func NewService(server, token string) *Service {
	return &Service{
		server: server,
		token:  token,
	}
}

func (s *Service) Shorten(ctx context.Context, url string, identifier Option[string]) (string, error) {
	var (
		req = shortenRequest{
			URL:        url,
			Identifier: identifier.OrEmpty(),
		}
		resp shortenResponse
	)

	if err := s.doRequest(ctx, http.MethodPost, "/api/shorten", req, &resp); err != nil {
		return "", fmt.Errorf("failed to shorten URL: %w", err)
	}

	return resp.ShortURL, nil
}

func (s *Service) Stats(ctx context.Context, identifier string) (*Stats, error) {
	var resp Stats

	if err := s.doRequest(ctx, http.MethodGet, fmt.Sprintf("/api/stats/%s", identifier), nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	return &resp, nil
}

func (s *Service) doRequest(ctx context.Context, method, path string, body, response any) error {
	fullURL, err := url.JoinPath(s.server, path)
	if err != nil {
		return fmt.Errorf("failed to build URL: %w", err)
	}

	encodedBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to encode body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, fullURL, bytes.NewReader(encodedBody))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.token))

	httpResp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to do request: %w", err)
	}

	switch httpResp.StatusCode {
	case http.StatusOK:
		if err := json.NewDecoder(httpResp.Body).Decode(response); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	case http.StatusUnauthorized:
		// TODO: use var errors
		return fmt.Errorf("unauthorized")
	case http.StatusConflict:
		return fmt.Errorf("identifier already exists")
	case http.StatusNotFound:
		return fmt.Errorf("identifier not found")
	default:
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}

	return nil
}

type Stats struct {
	Identifier  string `json:"identifier"`
	OriginalURL string `json:"original_url"`
	CreatedBy   string `json:"created_by"`
	Visits      int64  `json:"visits"`
}

type shortenRequest struct {
	URL        string `json:"url"`
	Identifier string `json:"identifier,omitempty"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}
