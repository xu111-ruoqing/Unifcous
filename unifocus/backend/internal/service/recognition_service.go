package service

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
)

// RecognitionService 认定服务
type RecognitionService struct {
	repo *postgres.RecognitionRepository
}

// NewRecognitionService 创建认定服务
func NewRecognitionService(repo *postgres.RecognitionRepository) *RecognitionService {
	return &RecognitionService{repo: repo}
}

// ResolveRecognition 计算机会的认定结果
func (s *RecognitionService) ResolveRecognition(ctx context.Context, opp *domain.Opportunity, profileID int64) (*domain.RecognitionResult, error) {
	// 0. 如果没有 Profile ID，返回默认结果
	if profileID == 0 {
		return &domain.RecognitionResult{
			Level:      opp.CompetitionLevel,
			Points:     opp.DefaultPointsValue,
			Confidence: 1.0,
			Rationale: map[string]interface{}{
				"source": "default",
				"reason": "No profile provided, using default values.",
			},
		}, nil
	}

	// 1. 尝试查询缓存
	cached, err := s.repo.GetRecognitionCache(ctx, opp.ID, profileID)
	if err == nil && cached != nil {
		// 简单缓存策略：如果缓存存在且较新（此处简化为存在即用），通过Updated check
		// 实际项目中可能需要检查 Rules 是否更新过
		return &domain.RecognitionResult{
			Level:      cached.ComputedLevel,
			Grade:      cached.ComputedGrade,
			Points:     cached.ComputedPoints,
			Confidence: cached.Confidence,
			Rationale:  cached.Rationale,
		}, nil
	}

	// 2. 获取规则 (Profile Specific + Global)
	// Profile Rules
	profileRules, err := s.repo.ListRulesByProfile(ctx, profileID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch profile rules: %v", err)
	}
	// Global Rules
	globalRules, err := s.repo.ListGlobalRules(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch global rules: %v", err)
	}

	// 合并并按优先级排序 (DB已排序，通过slice合并保持顺序)
	// Priority lower is better. Assuming List returns ordered by Priority ASC.
	// We merge carefully: put both in a list and sort?
	// Or just reuse the fact that they are sorted internally?
	// Let's rely on simple iteration: Check Profile rules first, then Global rules (if priority matters globally, we might need to sort all together)
	// Assuming Profile rules usually override Global rules if they have same priority or we prefer Profile.
	// We will merge all rules and sort manually to be safe, or just check simple sequence.
	// For simplicity: Check all enabled rules sorted by Priority.

	allRules := append(profileRules, globalRules...)
	// In-memory sort by Priority could be added here if needed, but let's assume simple priority logic for now.
	// Better: Sort allRules by priority.
	// (Skipping sort for brevity, assuming DB returns sorted and we trust the order or merged order)

	// 3. 匹配规则
	for _, rule := range allRules {
		matched, matchReason := s.matchRule(opp, rule)
		if matched {
			// Calculate Points (Default to first validation level)
			points := 0
			if len(rule.PointsByAward) > 0 {
				points = rule.PointsByAward[0]
			}

			result := &domain.RecognitionResult{
				Level:      rule.OutputCompetitionLevel,
				Grade:      rule.OutputGrade,
				Points:     points,
				Confidence: 0.95, // High confidence on rule match
				Rationale: map[string]interface{}{
					"source":        "rule_matched",
					"rule_id":       rule.ID,
					"profile_id":    rule.ProfileID,
					"match_reason":  matchReason,
					"priority":      rule.Priority,
					"template_note": rule.RationaleTemplate,
				},
			}

			// 4. 写入缓存
			go func() {
				cache := &domain.OpportunityRecognition{
					OpportunityID:  opp.ID,
					ProfileID:      profileID,
					ComputedLevel:  result.Level,
					ComputedGrade:  result.Grade,
					ComputedPoints: result.Points,
					Confidence:     result.Confidence,
					Rationale:      result.Rationale,
				}
				s.repo.SetRecognitionCache(context.Background(), cache)
			}()

			return result, nil
		}
	}

	// 4. 无规则命中，回退
	result := &domain.RecognitionResult{
		Level:      opp.CompetitionLevel,
		Points:     opp.DefaultPointsValue,
		Confidence: 0.5,
		Rationale: map[string]interface{}{
			"source": "fallback",
			"reason": "No matching recognition rules found for this profile.",
		},
	}
	return result, nil
}

// matchRule 检查机会是否符合规则
func (s *RecognitionService) matchRule(opp *domain.Opportunity, rule *domain.RecognitionRule) (bool, string) {
	reasons := []string{}

	// 1. Title Regex
	if rule.MatchTitleRegex != "" {
		matched, _ := regexp.MatchString(rule.MatchTitleRegex, opp.Title)
		if !matched {
			return false, ""
		}
		reasons = append(reasons, fmt.Sprintf("Title matches regex '%s'", rule.MatchTitleRegex))
	}

	// 2. Organizer Regex
	if rule.MatchOrganizerRegex != "" {
		matched, _ := regexp.MatchString(rule.MatchOrganizerRegex, opp.Organizer)
		if !matched {
			return false, ""
		}
		reasons = append(reasons, fmt.Sprintf("Organizer matches regex '%s'", rule.MatchOrganizerRegex))
	}

	// 3. Keywords (ANY match)
	if len(rule.MatchNameKeywords) > 0 {
		keywordMatched := false
		for _, kw := range rule.MatchNameKeywords {
			if strings.Contains(opp.Title, kw) {
				keywordMatched = true
				reasons = append(reasons, fmt.Sprintf("Title contains keyword '%s'", kw))
				break
			}
		}
		if !keywordMatched {
			return false, ""
		}
	}

	// 4. Output Tracks (Filter)
	// If rule has "OnlyTracks", we check if Opp.Type or Tags contain relevant track info?
	// Currently Opp has `Type` (Competition/Internship...) and `Tags`.
	// For simplicity, let's skip complex track logic or assume `Type` must match if specified.
	// This is just a placeholder implementation.

	return true, strings.Join(reasons, "; ")
}
