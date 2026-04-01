package crawler

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/unifocus/backend/internal/crawler/scrapers"
	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/pkg/logger"
)

// Pipeline 定义数据处理管道接口
type Pipeline interface {
	Process(ctx context.Context, task *domain.CrawlTask, items []scrapers.RawOpportunity) error
}

// Scheduler 爬虫调度器
type Scheduler struct {
	repo           domain.CrawlTaskRepository
	pipeline       Pipeline
	scraperFactory map[string]scrapers.Scraper

	ctx       context.Context
	cancel    context.CancelFunc
	wg        sync.WaitGroup
	ticker    *time.Ticker
	isRunning bool
	mu        sync.Mutex
}

// NewScheduler 创建新的调度器
func NewScheduler(repo domain.CrawlTaskRepository, pipeline Pipeline) *Scheduler {
	return &Scheduler{
		repo:           repo,
		pipeline:       pipeline,
		scraperFactory: make(map[string]scrapers.Scraper),
	}
}

// RegisterScraper 注册爬虫
func (s *Scheduler) RegisterScraper(scraper scrapers.Scraper) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scraperFactory[scraper.Name()] = scraper
}

// SetPipeline 设置处理管道
func (s *Scheduler) SetPipeline(p Pipeline) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pipeline = p
}

// Start 启动调度器
func (s *Scheduler) Start() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return
	}

	logger.Info("Starting crawler scheduler...")
	ctx, cancel := context.WithCancel(context.Background())
	s.ctx = ctx
	s.cancel = cancel
	s.isRunning = true
	s.ticker = time.NewTicker(1 * time.Minute) // 每分钟检查一次

	s.wg.Add(1)
	go s.run()
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return
	}

	logger.Info("Stopping crawler scheduler...")
	s.ticker.Stop()
	s.cancel()
	s.wg.Wait()
	s.isRunning = false
}

// run 调度循环
func (s *Scheduler) run() {
	defer s.wg.Done()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			s.processDueTasks()
		}
	}
}

// processDueTasks 处理到期任务
func (s *Scheduler) processDueTasks() {
	// 获取到期任务
	tasks, err := s.repo.GetDueTasks(s.ctx, 10) // 每次处理10个
	if err != nil {
		logger.Errorf("Failed to get due crawl tasks: %v", err)
		return
	}

	if len(tasks) == 0 {
		return
	}

	logger.Infof("Found %d due crawl tasks", len(tasks))

	for _, task := range tasks {
		s.wg.Add(1)
		go func(t *domain.CrawlTask) {
			defer s.wg.Done()
			s.executeTask(t)
		}(task)
	}
}

// executeTask 执行单个爬取任务
func (s *Scheduler) executeTask(task *domain.CrawlTask) {
	logger.Infof("Executing crawl task: %s (%s)", task.TargetURL, task.SiteName)

	// 更新状态为运行中
	task.Status = domain.CrawlStatusRunning
	if err := s.repo.Update(s.ctx, task); err != nil {
		logger.Errorf("Failed to update task status to running: %v", err)
		return
	}

	// 选择爬虫
	scraperName := task.SiteName
	scraper, ok := s.scraperFactory[scraperName]
	if !ok {
		// 回退到默认爬虫
		scraperName = "static"
		scraper, ok = s.scraperFactory[scraperName]
	}

	var err error
	var items []scrapers.RawOpportunity

	if !ok {
		err = fmt.Errorf("scraper not found: %s", scraperName)
	} else {
		items, err = scraper.Scrape(s.ctx, task)
	}

	// 处理结果
	now := time.Now()
	task.LastCrawledAt = &now
	next := task.CalculateNextCrawl()
	task.NextCrawlAt = &next

	if err != nil {
		logger.Errorf("Crawl task failed: %v", err)
		task.Status = domain.CrawlStatusFailed
		task.ErrorMessage = err.Error()
	} else {
		logger.Infof("Crawl task success: %d items found", len(items))
		task.Status = domain.CrawlStatusSuccess
		task.ErrorMessage = ""

		// 只有成功才走Pipeline
		if s.pipeline != nil {
			if pErr := s.pipeline.Process(s.ctx, task, items); pErr != nil {
				logger.Errorf("Pipeline processing failed: %v", pErr)
				// 管道失败是否算任务失败？暂不算，因为爬取是成功的
			}
		}
	}

	if updateErr := s.repo.Update(s.ctx, task); updateErr != nil {
		logger.Errorf("Failed to update task status after execution: %v", updateErr)
	}
}
