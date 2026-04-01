package service

import (
	"context"

	"github.com/unifocus/backend/internal/crawler"
	"github.com/unifocus/backend/internal/domain"
)

// CrawlerService 定义爬虫服务接口
type CrawlerService interface {
	// 任务管理
	CreateTask(ctx context.Context, task *domain.CrawlTask) error
	GetTask(ctx context.Context, id int64) (*domain.CrawlTask, error)
	ListTasks(ctx context.Context, offset, limit int) ([]*domain.CrawlTask, error)
	UpdateTask(ctx context.Context, task *domain.CrawlTask) error
	DeleteTask(ctx context.Context, id int64) error

	// 调度控制
	StartScheduler()
	StopScheduler()
}

type crawlerService struct {
	repo      domain.CrawlTaskRepository
	scheduler *crawler.Scheduler
}

// NewCrawlerService 创建新的爬虫服务
func NewCrawlerService(repo domain.CrawlTaskRepository, scheduler *crawler.Scheduler) CrawlerService {
	return &crawlerService{
		repo:      repo,
		scheduler: scheduler,
	}
}

func (s *crawlerService) CreateTask(ctx context.Context, task *domain.CrawlTask) error {
	// 设置默认状态
	if task.Status == "" {
		task.Status = domain.CrawlStatusPending
	}
	// 计算下一次运行时间
	if task.NextCrawlAt == nil {
		next := task.CalculateNextCrawl()
		task.NextCrawlAt = &next
	}

	return s.repo.Create(ctx, task)
}

func (s *crawlerService) GetTask(ctx context.Context, id int64) (*domain.CrawlTask, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *crawlerService) ListTasks(ctx context.Context, offset, limit int) ([]*domain.CrawlTask, error) {
	return s.repo.List(ctx, offset, limit)
}

func (s *crawlerService) UpdateTask(ctx context.Context, task *domain.CrawlTask) error {
	return s.repo.Update(ctx, task)
}

func (s *crawlerService) DeleteTask(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}

func (s *crawlerService) StartScheduler() {
	if s.scheduler != nil {
		s.scheduler.Start()
	}
}

func (s *crawlerService) StopScheduler() {
	if s.scheduler != nil {
		s.scheduler.Stop()
	}
}
