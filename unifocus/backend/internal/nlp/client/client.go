package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// NLPClient NLP服务客户端
type NLPClient struct {
	baseURL    string
	httpClient *http.Client
}

// NewNLPClient 创建NLP客户端
func NewNLPClient(baseURL string) *NLPClient {
	return &NLPClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// VectorizeRequest 向量化请求
type VectorizeRequest struct {
	Text string `json:"text"`
}

// VectorizeResponse 向量化响应
type VectorizeResponse struct {
	Vector []float32 `json:"vector"`
	Error  string    `json:"error,omitempty"`
}

// ExtractRequest 信息提取请求
type ExtractRequest struct {
	Text string `json:"text"`
}

// ExtractResponse 信息提取响应
type ExtractResponse struct {
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	StartDate   string   `json:"start_date,omitempty"`
	Deadline    string   `json:"deadline,omitempty"`
	EventDate   string   `json:"event_date,omitempty"`
	Location    string   `json:"location,omitempty"`
	Organizer   string   `json:"organizer,omitempty"`
	Description string   `json:"description,omitempty"`
	Error       string   `json:"error,omitempty"`
}

// Vectorize 文本向量化
func (c *NLPClient) Vectorize(ctx context.Context, text string) ([]float32, error) {
	reqBody := VectorizeRequest{Text: text}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/vectorize", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result VectorizeResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("nlp service error: %s", result.Error)
	}

	return result.Vector, nil
}

// Extract 信息提取
func (c *NLPClient) Extract(ctx context.Context, text string) (*ExtractResponse, error) {
	reqBody := ExtractRequest{Text: text}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/extract", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result ExtractResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("nlp service error: %s", result.Error)
	}

	return &result, nil
}

// HealthCheck 健康检查
func (c *NLPClient) HealthCheck(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/health", nil)
	if err != nil {
		return fmt.Errorf("create request failed: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return nil
}
