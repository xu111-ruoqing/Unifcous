package scrapers

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/unifocus/backend/internal/domain"
)

// CyNcssScraper 互联网+大赛爬虫 (cy.ncss.cn)
type CyNcssScraper struct {
	*BaseScraper
	httpClient *http.Client
}

// NewCyNcssScraper 创建互联网+大赛爬虫
func NewCyNcssScraper() *CyNcssScraper {
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
	}

	return &CyNcssScraper{
		BaseScraper: NewBaseScraper("cy_ncss", userAgents),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return nil // 允许重定向
			},
		},
	}
}

// Scrape 爬取互联网+大赛网站
func (s *CyNcssScraper) Scrape(ctx context.Context, task *domain.CrawlTask) ([]RawOpportunity, error) {
	s.WaitForRateLimit()

	homeURL := task.TargetURL
	if homeURL == "" {
		homeURL = "https://cy.ncss.cn"
	}

	items, err := s.scrapeHome(ctx, homeURL)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// scrapeHome 从首页抓取大赛动态和省市高校动态
func (s *CyNcssScraper) scrapeHome(ctx context.Context, pageURL string) ([]RawOpportunity, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置请求头
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://cy.ncss.cn/")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %w", err)
	}

	items := make([]RawOpportunity, 0)
	now := time.Now()
	base, _ := url.Parse(pageURL)
	baseURL := base.Scheme + "://" + base.Host

	// A) 大赛动态
	doc.Find(".dynamics .dynamics-block").Each(func(i int, block *goquery.Selection) {
		title := strings.TrimSpace(block.Find(".dynamics-text-title").Text())
		desc := strings.TrimSpace(block.Find(".dynamics-text-content").Text())
		if desc == "" {
			desc = title
		}
		linkSel := block.Find("a[href^=\"/information/\"]").First()
		href, ok := linkSel.Attr("href")
		if !ok || href == "" {
			return
		}
		sourceURL := resolveURL(baseURL, href)

		raw := FetchPlainText(s.httpClient, sourceURL)

		items = append(items, RawOpportunity{
			Title:       title,
			Description: desc,
			RawText:     raw,
			SourceURL:   sourceURL,
			HTMLContent: "",
			ExtractedAt: now,
		})
	})

	// B) 省市高校动态
	doc.Find(".province .province-main .p-block").Each(func(i int, block *goquery.Selection) {
		title := strings.TrimSpace(block.Find(".p-block-text-content").Text())
		if title == "" {
			return
		}
		desc := title
		linkSel := block.Find("a[href^=\"/information/\"]").First()
		href, ok := linkSel.Attr("href")
		if !ok || href == "" {
			return
		}
		sourceURL := resolveURL(baseURL, href)

		raw := FetchPlainText(s.httpClient, sourceURL)

		items = append(items, RawOpportunity{
			Title:       title,
			Description: desc,
			RawText:     raw,
			SourceURL:   sourceURL,
			HTMLContent: "",
			ExtractedAt: now,
		})
	})

	return items, nil
}
