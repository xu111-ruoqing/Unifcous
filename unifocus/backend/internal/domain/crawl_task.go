package domain

import (
	"context"
	"time"
)

// CrawlStatus 爬虫任务状态
type CrawlStatus string

const (
	CrawlStatusPending CrawlStatus = "pending"
	CrawlStatusRunning CrawlStatus = "running"
	CrawlStatusSuccess CrawlStatus = "success"
	CrawlStatusFailed  CrawlStatus = "failed"
	CrawlStatusStopped CrawlStatus = "stopped" // 手动停止
)

// SelectorConfig 页面选择器配置
type SelectorConfig struct {
	TitleSelector   string `json:"title"`
	LinkSelector    string `json:"link"`
	DateSelector    string `json:"date"`
	ContentSelector string `json:"content"`
	ItemSelector    string `json:"item"` // 列表项容器
}

// CrawlTask 爬虫任务实体
type CrawlTask struct {
	ID             int64           `json:"id"`
	TargetURL      string          `json:"target_url"`
	SiteName       string          `json:"site_name"`
	SelectorConfig *SelectorConfig `json:"selector_config"` // 存储为 JSONB
	Frequency      string          `json:"frequency"`       // hourly, daily, weekly, or cron expression

	LastCrawledAt *time.Time  `json:"last_crawled_at"`
	NextCrawlAt   *time.Time  `json:"next_crawl_at"`
	Status        CrawlStatus `json:"status"`
	ErrorMessage  string      `json:"error_message,omitempty"`

	CreatedAt time.Time `json:"created_at"`
}

// TableName 指定数据库表名
func (CrawlTask) TableName() string {
	return "crawl_tasks"
}

// IsDue 检查任务是否到期需要运行
func (t *CrawlTask) IsDue() bool {
	if t.NextCrawlAt == nil {
		return true
	}
	return time.Now().After(*t.NextCrawlAt)
}

// CalculateNextCrawl 根据频率计算下一次运行时间
func (t *CrawlTask) CalculateNextCrawl() time.Time {
	now := time.Now()
	switch t.Frequency {
	case "hourly":
		return now.Add(time.Hour)
	case "daily":
		return now.Add(24 * time.Hour)
	case "weekly":
		return now.Add(7 * 24 * time.Hour)
	default:
		// 默认为每天
		return now.Add(24 * time.Hour)
	}
}

// CrawlTaskRepository 爬虫任务仓库接口
type CrawlTaskRepository interface {
	Create(ctx context.Context, task *CrawlTask) error
	GetByID(ctx context.Context, id int64) (*CrawlTask, error)
	List(ctx context.Context, offset, limit int) ([]*CrawlTask, error)
	Update(ctx context.Context, task *CrawlTask) error
	Delete(ctx context.Context, id int64) error
	GetDueTasks(ctx context.Context, limit int) ([]*CrawlTask, error) // 获取到期的任务
}
