package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/unifocus/backend/internal/config"
)

type NLPClient struct {
	baseURL    string
	httpClient *http.Client
}

type VectorizeRequest struct {
	Text string `json:"text"`
}

type VectorizeResponse struct {
	Vector []float32 `json:"vector"`
	Dim    int       `json:"dim"`
}

func NewNLPClient(cfg *config.NLPServiceConfig) *NLPClient {
	return &NLPClient{
		baseURL: cfg.URL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// VectorizeText calls Python service to get embeddings
func (c *NLPClient) VectorizeText(ctx context.Context, text string) ([]float32, error) {
	url := fmt.Sprintf("%s/api/v1/vectorize/text", c.baseURL)
	reqBody := VectorizeRequest{Text: text}
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call NLP service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NLP service returned status: %d", resp.StatusCode)
	}

	var respData VectorizeResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return nil, err
	}

	return respData.Vector, nil
}

// ExtractTextFromPDF (Not implemented per new requirement)
func (c *NLPClient) ExtractTextFromPDF(ctx context.Context, fileData []byte) (string, error) {
	return "", fmt.Errorf("pdf extraction deprecated")
}

// ExtractSkills (Placeholder)
func (c *NLPClient) ExtractSkills(ctx context.Context, text string) ([]string, error) {
	// TODO: Call python extract skills endpoint
	return []string{}, nil
}
