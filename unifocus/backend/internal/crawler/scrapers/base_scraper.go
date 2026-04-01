package scrapers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/unifocus/backend/internal/domain"
)

// RawOpportunity 爬虫提取的原始机会数据
type RawOpportunity struct {
	Title       string
	Description string
	SourceURL   string
	HTMLContent string
	RawText     string
	ExtractedAt time.Time
}

// Scraper 爬虫接口
type Scraper interface {
	// Scrape 执行爬取任务，返回原始机会数据列表
	Scrape(ctx context.Context, task *domain.CrawlTask) ([]RawOpportunity, error)

	// Name 返回爬虫名称
	Name() string
}

// BaseScraper 基础爬虫，提供通用功能
type BaseScraper struct {
	name        string
	userAgents  []string
	rateLimiter *RateLimiter
}

// NewBaseScraper 创建基础爬虫
func NewBaseScraper(name string, userAgents []string) *BaseScraper {
	return &BaseScraper{
		name:        name,
		userAgents:  userAgents,
		rateLimiter: NewRateLimiter(2.0, 5), // 默认每秒2个请求，突发5个
	}
}

// Name 返回爬虫名称
func (b *BaseScraper) Name() string {
	return b.name
}

// WaitForRateLimit 等待速率限制
func (b *BaseScraper) WaitForRateLimit() {
	b.rateLimiter.Wait()
}

// GetRandomUserAgent 获取随机User-Agent
func (b *BaseScraper) GetRandomUserAgent() string {
	if len(b.userAgents) == 0 {
		return "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"
	}
	// 简单实现：使用时间戳取模
	index := int(time.Now().UnixNano()) % len(b.userAgents)
	return b.userAgents[index]
}

// RateLimiter 速率限制器
type RateLimiter struct {
	requestsPerSecond float64
	burst             int
	lastRequest       time.Time
	requestCount      int
}

// NewRateLimiter 创建速率限制器
func NewRateLimiter(requestsPerSecond float64, burst int) *RateLimiter {
	return &RateLimiter{
		requestsPerSecond: requestsPerSecond,
		burst:             burst,
		lastRequest:       time.Now(),
	}
}

// Wait 等待直到可以发送下一个请求
func (r *RateLimiter) Wait() {
	now := time.Now()
	elapsed := now.Sub(r.lastRequest)

	// 计算应该等待的时间
	interval := time.Duration(float64(time.Second) / r.requestsPerSecond)

	if elapsed < interval {
		time.Sleep(interval - elapsed)
	}

	r.lastRequest = time.Now()
}

// FetchPlainText 获取页面纯文本（移除script/style并折叠空白）
func FetchPlainText(client *http.Client, pageURL string) string {
	req, err := http.NewRequest("GET", pageURL, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return ""
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return ""
	}
	doc.Find("script, style").Remove()
	text := strings.TrimSpace(doc.Text())
	return strings.Join(strings.Fields(text), " ")
}
