package service

import (
	"context"
	"math"
	"time"

	"github.com/unifocus/backend/internal/domain"
)

// ScoringConfig 评分模型超参数（论文式1）
// S = 100 · g · σ(α·T + β·R + γ·L − δ·C)
type ScoringConfig struct {
	Alpha float64 // 时间可行性权重 T
	Beta  float64 // 收益权重 R
	Gamma float64 // 地域偏好权重 L
	Delta float64 // 准备成本权重 C
}

// DefaultScoringConfig 默认权重（论文表3）
var DefaultScoringConfig = ScoringConfig{
	Alpha: 0.35,
	Beta:  0.30,
	Gamma: 0.15,
	Delta: 0.20,
}

// ScoreResult 评分结果
type ScoreResult struct {
	Score       float64            `json:"score"`        // 最终得分 [0, 100]
	Passed      bool               `json:"passed"`       // 是否通过硬性门槛 g
	Components  ScoreComponents    `json:"components"`   // 各分量明细
	Explanation map[string]interface{} `json:"explanation"` // 解释性信息
}

// ScoreComponents 评分分量
type ScoreComponents struct {
	TimeFeasibility float64 `json:"time_feasibility"` // T [0,1]
	Reward          float64 `json:"reward"`           // R [0,1]
	LocationPref    float64 `json:"location_pref"`    // L [0,1]
	PrepCost        float64 `json:"prep_cost"`        // C [0,1]，越高代价越大
	LogitRaw        float64 `json:"logit_raw"`        // α·T + β·R + γ·L − δ·C
}

// ScoringService 可达性评分服务（论文第 2.4 节）
type ScoringService struct {
	cfg ScoringConfig
}

// NewScoringService 创建评分服务
func NewScoringService() *ScoringService {
	return &ScoringService{cfg: DefaultScoringConfig}
}

// NewScoringServiceWithConfig 使用自定义权重
func NewScoringServiceWithConfig(cfg ScoringConfig) *ScoringService {
	return &ScoringService{cfg: cfg}
}

// Score 计算单个机会对用户的可达性评分
//
// 参数：
//   opp     - 机会实体
//   profile - 用户画像（可为 nil，降级为纯内容评分）
//   user    - 用户基本信息（用于地域偏好、年级门槛）
func (s *ScoringService) Score(
	ctx context.Context,
	opp *domain.Opportunity,
	profile *domain.UserProfile,
	user *domain.User,
) *ScoreResult {

	// ── 硬性门槛 g ──────────────────────────────────────────
	// g = 0 表示用户不符合基本资格，直接屏蔽（论文 2.4.1）
	g := s.hardGate(opp, user)

	// ── 分量计算 ─────────────────────────────────────────────
	T := s.timeFeasibility(opp)
	R := s.reward(opp)
	L := s.locationPreference(opp, user)
	C := s.prepCost(opp, profile)

	// ── Logit 原始值 ─────────────────────────────────────────
	// α·T + β·R + γ·L − δ·C
	logit := s.cfg.Alpha*T + s.cfg.Beta*R + s.cfg.Gamma*L - s.cfg.Delta*C

	// ── Sigmoid 压缩到 (0,1) ─────────────────────────────────
	sigma := sigmoid(logit)

	// ── 最终得分 S = 100 · g · σ ─────────────────────────────
	score := 100.0 * float64(g) * sigma

	return &ScoreResult{
		Score:  roundTo(score, 2),
		Passed: g == 1,
		Components: ScoreComponents{
			TimeFeasibility: roundTo(T, 4),
			Reward:          roundTo(R, 4),
			LocationPref:    roundTo(L, 4),
			PrepCost:        roundTo(C, 4),
			LogitRaw:        roundTo(logit, 4),
		},
		Explanation: s.buildExplanation(opp, user, g, T, R, L, C),
	}
}

// ScoreBatch 批量评分并排序（用于列表接口）
func (s *ScoringService) ScoreBatch(
	ctx context.Context,
	opps []*domain.Opportunity,
	profile *domain.UserProfile,
	user *domain.User,
) []ScoredOpportunity {
	results := make([]ScoredOpportunity, 0, len(opps))
	for _, opp := range opps {
		sr := s.Score(ctx, opp, profile, user)
		results = append(results, ScoredOpportunity{
			Opportunity: opp,
			ScoreResult: sr,
		})
	}
	// 按得分降序排列
	sortScoredOpps(results)
	return results
}

// ScoredOpportunity 带评分的机会
type ScoredOpportunity struct {
	*domain.Opportunity
	ScoreResult *ScoreResult `json:"score_result"`
}

// ── 内部实现 ──────────────────────────────────────────────────────────────────

// hardGate 硬性门槛（返回 0 或 1）
// 当前实现的门槛规则：
//  1. 已过截止日期 → 0
//  2. 用户年级不满足最低要求（若机会有明确年级限制） → 0
func (s *ScoringService) hardGate(opp *domain.Opportunity, user *domain.User) int {
	now := time.Now()

	// 门槛1：截止日期
	if opp.Deadline != nil && opp.Deadline.Before(now) {
		return 0
	}
	if opp.RegistrationDeadline != nil && opp.RegistrationDeadline.Before(now) {
		return 0
	}

	// 门槛2：专业匹配（若 TargetMajors 非空且非"不限"，检查用户专业）
	if user != nil && len(opp.TargetMajors) > 0 {
		hasAll := false
		for _, m := range opp.TargetMajors {
			if m == "不限" || m == "all" {
				hasAll = true
				break
			}
		}
		if !hasAll {
			majorMatched := false
			for _, m := range opp.TargetMajors {
				if containsAny(user.Major, m) {
					majorMatched = true
					break
				}
			}
			if !majorMatched {
				return 0
			}
		}
	}

	return 1
}

// timeFeasibility T：时间可行性 [0, 1]
// 逻辑：距截止日期越近 → 分越低；时间充裕 → 接近 1
func (s *ScoringService) timeFeasibility(opp *domain.Opportunity) float64 {
	now := time.Now()
	deadline := opp.Deadline
	if deadline == nil {
		deadline = opp.RegistrationDeadline
	}
	if deadline == nil {
		// 无截止日期信息，给中等分
		return 0.5
	}

	daysLeft := deadline.Sub(now).Hours() / 24
	if daysLeft <= 0 {
		return 0.0
	}
	// 超过 60 天：满分；30~60 天：线性递减至 0.7；7~30 天：继续递减；<7 天：低分
	switch {
	case daysLeft >= 60:
		return 1.0
	case daysLeft >= 30:
		return 0.7 + 0.3*(daysLeft-30)/30
	case daysLeft >= 7:
		return 0.3 + 0.4*(daysLeft-7)/23
	default:
		return 0.1 + 0.2*(daysLeft/7)
	}
}

// reward R：收益评分 [0, 1]
// 基于 CompetitionLevel（竞赛级别）和 DefaultPointsValue
func (s *ScoringService) reward(opp *domain.Opportunity) float64 {
	base := 0.3
	switch opp.CompetitionLevel {
	case "国家级", "国际级":
		base = 1.0
	case "省级":
		base = 0.75
	case "校级":
		base = 0.5
	case "院级":
		base = 0.35
	}

	// 用积分值微调（归一化到 0-0.1 的加成）
	if opp.DefaultPointsValue > 0 {
		bonus := math.Min(float64(opp.DefaultPointsValue)/100.0*0.1, 0.1)
		base = math.Min(base+bonus, 1.0)
	}

	// 官方认证加成
	if opp.IsOfficial {
		base = math.Min(base+0.05, 1.0)
	}

	return base
}

// locationPreference L：地域偏好 [0, 1]
// 线上 / 全国 → 高分；需要出行 → 扣分
func (s *ScoringService) locationPreference(opp *domain.Opportunity, user *domain.User) float64 {
	loc := opp.Location
	if loc == "" {
		loc = opp.LocationRaw
	}
	if loc == "" {
		return 0.5
	}

	// 线上或全国
	if containsAny(loc, "线上", "网络", "online", "全国", "不限") {
		return 1.0
	}

	// 用户所在城市匹配
	if user != nil && user.School != "" {
		if containsAny(loc, user.School) {
			return 0.9
		}
	}

	// 省内
	return 0.4
}

// prepCost C：准备成本 [0, 1]（越高代价越大，评分公式中用减法）
// 基于 Requirements 的复杂度和竞赛类型
func (s *ScoringService) prepCost(opp *domain.Opportunity, profile *domain.UserProfile) float64 {
	cost := 0.5 // 基准

	// 技术类竞赛准备成本更高
	switch opp.Type {
	case "竞赛":
		cost = 0.7
	case "实习":
		cost = 0.6
	case "讲座":
		cost = 0.1
	case "奖学金":
		cost = 0.4
	case "科研":
		cost = 0.8
	}

	// 国际级别额外提高成本
	if opp.CompetitionLevel == "国际级" {
		cost = math.Min(cost+0.15, 1.0)
	}

	// 用户有相关技能 → 降低成本
	if profile != nil && len(profile.Skills) > 0 {
		skillMatch := 0
		for _, tag := range opp.Tags {
			for _, skill := range profile.Skills {
				if containsAny(skill, tag) || containsAny(tag, skill) {
					skillMatch++
				}
			}
		}
		if skillMatch > 0 {
			reduction := math.Min(float64(skillMatch)*0.05, 0.25)
			cost = math.Max(cost-reduction, 0.0)
		}
	}

	return cost
}

// buildExplanation 构建解释性信息（用于前端展示）
func (s *ScoringService) buildExplanation(
	opp *domain.Opportunity,
	user *domain.User,
	g int,
	T, R, L, C float64,
) map[string]interface{} {
	now := time.Now()
	deadline := opp.Deadline
	if deadline == nil {
		deadline = opp.RegistrationDeadline
	}

	exp := map[string]interface{}{
		"model":       "S = 100 · g · σ(αT + βR + γL − δC)",
		"weights":     map[string]float64{"α": s.cfg.Alpha, "β": s.cfg.Beta, "γ": s.cfg.Gamma, "δ": s.cfg.Delta},
		"gate_passed": g == 1,
	}

	if deadline != nil {
		daysLeft := int(deadline.Sub(now).Hours() / 24)
		exp["days_until_deadline"] = daysLeft
	}

	if g == 0 {
		exp["gate_reason"] = s.gateReason(opp, user)
	}

	return exp
}

func (s *ScoringService) gateReason(opp *domain.Opportunity, user *domain.User) string {
	now := time.Now()
	if opp.Deadline != nil && opp.Deadline.Before(now) {
		return "报名截止日期已过"
	}
	if opp.RegistrationDeadline != nil && opp.RegistrationDeadline.Before(now) {
		return "注册截止日期已过"
	}
	if user != nil && len(opp.TargetMajors) > 0 {
		return "专业不匹配"
	}
	return "未满足资格要求"
}

// ── 工具函数 ─────────────────────────────────────────────────────────────────

func sigmoid(x float64) float64 {
	return 1.0 / (1.0 + math.Exp(-x))
}

func roundTo(v float64, decimals int) float64 {
	factor := math.Pow(10, float64(decimals))
	return math.Round(v*factor) / factor
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if len(sub) == 0 {
			continue
		}
		if len(s) == 0 {
			continue
		}
		if len(s) >= len(sub) {
			for i := 0; i <= len(s)-len(sub); i++ {
				if s[i:i+len(sub)] == sub {
					return true
				}
			}
		}
	}
	return false
}

func sortScoredOpps(opps []ScoredOpportunity) {
	// 简单插入排序（列表通常不超过 200 条）
	for i := 1; i < len(opps); i++ {
		key := opps[i]
		j := i - 1
		for j >= 0 && opps[j].ScoreResult.Score < key.ScoreResult.Score {
			opps[j+1] = opps[j]
			j--
		}
		opps[j+1] = key
	}
}
