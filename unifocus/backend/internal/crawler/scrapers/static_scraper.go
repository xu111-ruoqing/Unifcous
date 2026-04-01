package scrapers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/unifocus/backend/internal/domain"
)

// StaticScraper 静态页面爬虫（使用Colly）
type StaticScraper struct {
	*BaseScraper
	httpClient *http.Client
}

// NewStaticScraper 创建静态页面爬虫
func NewStaticScraper(userAgents []string) *StaticScraper {
	return &StaticScraper{
		BaseScraper: NewBaseScraper("static", userAgents),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Scrape 爬取静态页面
func (s *StaticScraper) Scrape(ctx context.Context, task *domain.CrawlTask) ([]RawOpportunity, error) {
	s.WaitForRateLimit()

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, "GET", task.TargetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置User-Agent
	req.Header.Set("User-Agent", s.GetRandomUserAgent())
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 使用goquery解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	// 提取机会数据（根据selector配置）
	opportunities := []RawOpportunity{}

	// 从selector_config中获取选择器配置
	// 这里简化处理，实际应该从task.SelectorConfig中读取
	selectorConfig := map[string]string{
		"item":    ".opportunity-item, .news-item, .notice-item",
		"title":   "h3, .title, a",
		"link":    "a",
		"date":    ".date, .time",
		"content": ".content, .description, p",
	}

	doc.Find(selectorConfig["item"]).Each(func(i int, s *goquery.Selection) {
		title := strings.TrimSpace(s.Find(selectorConfig["title"]).First().Text())
		link, _ := s.Find(selectorConfig["link"]).First().Attr("href")
		content := strings.TrimSpace(s.Find(selectorConfig["content"]).First().Text())

		if title == "" {
			return
		}

		// 处理相对链接
		if link != "" && !strings.HasPrefix(link, "http") {
			link = resolveURL(task.TargetURL, link)
		}

		opp := RawOpportunity{
			Title:       title,
			Description: content,
			SourceURL:   link,
			ExtractedAt: time.Now(),
		}

		opportunities = append(opportunities, opp)
	})

	return opportunities, nil
}

// resolveURL 解析相对URL为绝对URL
func resolveURL(baseURL, relativeURL string) string {
	if strings.HasPrefix(relativeURL, "http://") || strings.HasPrefix(relativeURL, "https://") {
		return relativeURL
	}

	if strings.HasPrefix(relativeURL, "/") {
		// 从baseURL提取协议和域名
		parts := strings.Split(baseURL, "/")
		if len(parts) >= 3 {
			return parts[0] + "//" + parts[2] + relativeURL
		}
	}

	// 相对路径
	lastSlash := strings.LastIndex(baseURL, "/")
	if lastSlash != -1 {
		return baseURL[:lastSlash+1] + relativeURL
	}

	return baseURL + "/" + relativeURL
}
