package domain

import (
	"context"
	"time"
)

// RecognitionProfile 认定档案
type RecognitionProfile struct {
	ID        int64     `json:"id" db:"id"`
	School    string    `json:"school" db:"school"`
	College   string    `json:"college,omitempty" db:"college"`
	Major     string    `json:"major,omitempty" db:"major"`
	EntryYear int       `json:"entry_year,omitempty" db:"entry_year"`
	IsDefault bool      `json:"is_default" db:"is_default"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RecognitionRule 认定规则
type RecognitionRule struct {
	ID        int64  `json:"id" db:"id"`
	ProfileID *int64 `json:"profile_id,omitempty" db:"profile_id"` // NULL means global default
	Enabled   bool   `json:"enabled" db:"enabled"`
	Priority  int    `json:"priority" db:"priority"` // Smaller value = higher priority

	// Match Conditions
	MatchTitleRegex     string   `json:"match_title_regex,omitempty" db:"match_title_regex"`
	MatchNameKeywords   []string `json:"match_name_keywords,omitempty" db:"match_name_keywords"`
	MatchOrganizerRegex string   `json:"match_organizer_regex,omitempty" db:"match_organizer_regex"`
	OnlyTracks          []string `json:"only_tracks,omitempty" db:"only_tracks"`

	// Output Results
	OutputCompetitionLevel string `json:"output_competition_level,omitempty" db:"output_competition_level"`
	OutputGrade            string `json:"output_grade,omitempty" db:"output_grade"`
	PointsByAward          []int  `json:"points_by_award,omitempty" db:"points_by_award"` // e.g., {100, 80, 60}
	RationaleTemplate      string `json:"rationale_template,omitempty" db:"rationale_template"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RecognitionResult 认定结果 (VO, Not a DB Table directly)
type RecognitionResult struct {
	Level      string                 `json:"level"`
	Grade      string                 `json:"grade"`
	Points     int                    `json:"points"`
	Confidence float64                `json:"confidence"`
	Rationale  map[string]interface{} `json:"rationale"`
}

// OpportunityRecognition 认定结果缓存 (Optional Cache Table)
type OpportunityRecognition struct {
	OpportunityID int64 `json:"opportunity_id" db:"opportunity_id"`
	ProfileID     int64 `json:"profile_id" db:"profile_id"`

	ComputedLevel  string  `json:"computed_level" db:"computed_level"`
	ComputedGrade  string  `json:"computed_grade" db:"computed_grade"`
	ComputedPoints int     `json:"computed_points" db:"computed_points"`
	Confidence     float64 `json:"confidence" db:"confidence"`
	Rationale      map[string]interface{} `json:"rationale" db:"rationale"`

	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RecognitionRepository 接口
type RecognitionRepository interface {
	// Profile Operations
	GetProfile(ctx context.Context, id int64) (*RecognitionProfile, error)
	ListProfiles(ctx context.Context) ([]*RecognitionProfile, error)
	CreateProfile(ctx context.Context, profile *RecognitionProfile) error

	// Rule Operations
	ListRulesByProfile(ctx context.Context, profileID int64) ([]*RecognitionRule, error)
	ListGlobalRules(ctx context.Context) ([]*RecognitionRule, error)
	CreateRule(ctx context.Context, rule *RecognitionRule) error
}
