package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/lib/pq"

	"github.com/unifocus/backend/internal/domain"
)

// RecognitionRepository 认定策略仓库实现
type RecognitionRepository struct {
	db *DB
}

// NewRecognitionRepository 创建认定策略仓库实例
func NewRecognitionRepository(db *DB) *RecognitionRepository {
	return &RecognitionRepository{db: db}
}

// GetProfile 根据ID获取认定档案
func (r *RecognitionRepository) GetProfile(ctx context.Context, id int64) (*domain.RecognitionProfile, error) {
	query := `
		SELECT id, school, college, major, entry_year, is_default, created_at, updated_at
		FROM recognition_profiles
		WHERE id = $1
	`
	var profile domain.RecognitionProfile
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&profile.ID,
		&profile.School,
		&profile.College,
		&profile.Major,
		&profile.EntryYear,
		&profile.IsDefault,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found
		}
		return nil, err
	}
	return &profile, nil
}

// ListProfiles 列出所有认定档案
func (r *RecognitionRepository) ListProfiles(ctx context.Context) ([]*domain.RecognitionProfile, error) {
	query := `
		SELECT id, school, college, major, entry_year, is_default, created_at, updated_at
		FROM recognition_profiles
		ORDER BY id ASC
	`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var profiles []*domain.RecognitionProfile
	for rows.Next() {
		var p domain.RecognitionProfile
		if err := rows.Scan(
			&p.ID,
			&p.School,
			&p.College,
			&p.Major,
			&p.EntryYear,
			&p.IsDefault,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		profiles = append(profiles, &p)
	}
	return profiles, nil
}

// CreateProfile 创建认定档案
func (r *RecognitionRepository) CreateProfile(ctx context.Context, profile *domain.RecognitionProfile) error {
	query := `
		INSERT INTO recognition_profiles (school, college, major, entry_year, is_default)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`
	return r.db.QueryRowContext(ctx, query,
		profile.School,
		profile.College,
		profile.Major,
		profile.EntryYear,
		profile.IsDefault,
	).Scan(&profile.ID, &profile.CreatedAt, &profile.UpdatedAt)
}

// ListRulesByProfile 获取指定Profile的认定规则
func (r *RecognitionRepository) ListRulesByProfile(ctx context.Context, profileID int64) ([]*domain.RecognitionRule, error) {
	query := `
		SELECT id, profile_id, enabled, priority, 
		       match_title_regex, match_name_keywords, match_organizer_regex, only_tracks,
		       output_competition_level, output_grade, points_by_award, rationale_template,
		       created_at, updated_at
		FROM recognition_rules
		WHERE profile_id = $1 AND enabled = true
		ORDER BY priority ASC
	`
	return r.scanRules(ctx, query, profileID)
}

// ListGlobalRules 获取全局默认认定规则
func (r *RecognitionRepository) ListGlobalRules(ctx context.Context) ([]*domain.RecognitionRule, error) {
	query := `
		SELECT id, profile_id, enabled, priority, 
		       match_title_regex, match_name_keywords, match_organizer_regex, only_tracks,
		       output_competition_level, output_grade, points_by_award, rationale_template,
		       created_at, updated_at
		FROM recognition_rules
		WHERE profile_id IS NULL AND enabled = true
		ORDER BY priority ASC
	`
	return r.scanRules(ctx, query)
}

// scanRules 辅助函数，扫描规则列表
func (r *RecognitionRepository) scanRules(ctx context.Context, query string, args ...interface{}) ([]*domain.RecognitionRule, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*domain.RecognitionRule
	for rows.Next() {
		var rule domain.RecognitionRule
		var pointsByAward []int64 // pq.Int64Array

		err := rows.Scan(
			&rule.ID,
			&rule.ProfileID,
			&rule.Enabled,
			&rule.Priority,
			&rule.MatchTitleRegex,
			pq.Array(&rule.MatchNameKeywords),
			&rule.MatchOrganizerRegex,
			pq.Array(&rule.OnlyTracks),
			&rule.OutputCompetitionLevel,
			&rule.OutputGrade,
			pq.Array(&pointsByAward),
			&rule.RationaleTemplate,
			&rule.CreatedAt,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Convert []int64 to []int
		if len(pointsByAward) > 0 {
			rule.PointsByAward = make([]int, len(pointsByAward))
			for i, v := range pointsByAward {
				rule.PointsByAward[i] = int(v)
			}
		}

		rules = append(rules, &rule)
	}
	return rules, nil
}

// CreateRule 创建认定规则
func (r *RecognitionRepository) CreateRule(ctx context.Context, rule *domain.RecognitionRule) error {
	query := `
		INSERT INTO recognition_rules (
			profile_id, enabled, priority, 
			match_title_regex, match_name_keywords, match_organizer_regex, only_tracks,
			output_competition_level, output_grade, points_by_award, rationale_template
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id, created_at, updated_at
	`

	// Convert []int to []int64 for pq.Array
	var pointsByAward []int64
	for _, p := range rule.PointsByAward {
		pointsByAward = append(pointsByAward, int64(p))
	}

	return r.db.QueryRowContext(ctx, query,
		rule.ProfileID,
		rule.Enabled,
		rule.Priority,
		rule.MatchTitleRegex,
		pq.Array(rule.MatchNameKeywords),
		rule.MatchOrganizerRegex,
		pq.Array(rule.OnlyTracks),
		rule.OutputCompetitionLevel,
		rule.OutputGrade,
		pq.Array(pointsByAward),
		rule.RationaleTemplate,
	).Scan(&rule.ID, &rule.CreatedAt, &rule.UpdatedAt)
}

// GetRecognitionCache 获取缓存的认定结果
func (r *RecognitionRepository) GetRecognitionCache(ctx context.Context, opportunityID, profileID int64) (*domain.OpportunityRecognition, error) {
	query := `
		SELECT computed_level, computed_grade, computed_points, confidence, rationale, updated_at
		FROM opportunity_recognitions
		WHERE opportunity_id = $1 AND profile_id = $2
	`
	var rec domain.OpportunityRecognition
	rec.OpportunityID = opportunityID
	rec.ProfileID = profileID

	var rationaleBytes []byte

	err := r.db.QueryRowContext(ctx, query, opportunityID, profileID).Scan(
		&rec.ComputedLevel,
		&rec.ComputedGrade,
		&rec.ComputedPoints,
		&rec.Confidence,
		&rationaleBytes,
		&rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(rationaleBytes, &rec.Rationale); err != nil {
		return nil, err
	}

	return &rec, nil
}

// SetRecognitionCache 设置/更新认定结果缓存
func (r *RecognitionRepository) SetRecognitionCache(ctx context.Context, rec *domain.OpportunityRecognition) error {
	query := `
		INSERT INTO opportunity_recognitions (
			opportunity_id, profile_id, 
			computed_level, computed_grade, computed_points, confidence, rationale, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		ON CONFLICT (opportunity_id, profile_id) 
		DO UPDATE SET
			computed_level = EXCLUDED.computed_level,
			computed_grade = EXCLUDED.computed_grade,
			computed_points = EXCLUDED.computed_points,
			confidence = EXCLUDED.confidence,
			rationale = EXCLUDED.rationale,
			updated_at = NOW()
	`
	rationaleBytes, err := json.Marshal(rec.Rationale)
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, query,
		rec.OpportunityID,
		rec.ProfileID,
		rec.ComputedLevel,
		rec.ComputedGrade,
		rec.ComputedPoints,
		rec.Confidence,
		rationaleBytes,
	)
	return err
}
