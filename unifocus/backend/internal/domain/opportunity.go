package domain

import "time"

// Opportunity represents a competition or academic opportunity.
type Opportunity struct {
	ID                   int64      `json:"id" db:"id"`
	Title                string     `json:"title" db:"title"`
	Type                 string     `json:"type" db:"type"`
	Description          string     `json:"description" db:"description"`
	SourceURL            string     `json:"source_url" db:"source_url"`
	SourceType           string     `json:"source_type" db:"source_type"`
	CompetitionLevel     string     `json:"competition_level" db:"competition_level"`
	CertificationType    string     `json:"certification_type" db:"certification_type"`
	Organizer            string     `json:"organizer" db:"organizer"`
	OrganizerType        string     `json:"organizer_type" db:"organizer_type"`
	AwardLevel           string     `json:"award_level" db:"award_level"`
	DefaultPointsValue   int        `json:"default_points_value" db:"default_points_value"`
	IsOfficial           bool       `json:"is_official" db:"is_official"`
	StartDate            *time.Time `json:"start_date,omitempty" db:"start_date"`
	Deadline             *time.Time `json:"deadline,omitempty" db:"deadline"`
	EventDate            *time.Time `json:"event_date,omitempty" db:"event_date"`
	EventStartAt         *time.Time `json:"event_start_at,omitempty" db:"event_start_at"`
	EventEndAt           *time.Time `json:"event_end_at,omitempty" db:"event_end_at"`
	PublishedAt          *time.Time `json:"published_at,omitempty" db:"published_at"`
	Location             string     `json:"location" db:"location"`
	TimeInfoRaw          string     `json:"time_info_raw" db:"time_info_raw"`
	LocationRaw          string     `json:"location_raw" db:"location_raw"`
	NoticeType           string     `json:"notice_type" db:"notice_type"`
	RegistrationStartAt  *time.Time `json:"registration_start_at,omitempty" db:"registration_start_at"`
	RegistrationDeadline *time.Time `json:"registration_deadline,omitempty" db:"registration_deadline"`
	RawText              string     `json:"raw_text,omitempty" db:"raw_text"`
	Requirements         string     `json:"requirements,omitempty" db:"requirements"`
	EligibilityRules     string     `json:"eligibility_rules,omitempty" db:"eligibility_rules"`
	TargetMajors         []string   `json:"target_majors,omitempty" db:"target_majors"`
	Tags                 []string   `json:"tags,omitempty" db:"tags"`
	Attachments          []string   `json:"attachments,omitempty" db:"attachments"`
	DescriptionVector    []float32  `json:"description_vector,omitempty" db:"description_vector"`
	IsActive             bool       `json:"is_active" db:"is_active"`
	ViewCount            int        `json:"view_count" db:"view_count"`
	SaveCount            int        `json:"save_count" db:"save_count"`
	CreatedAt            time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at" db:"updated_at"`
	Recognition          *RecognitionResult `json:"recognition,omitempty"`
}

// CreateOpportunityRequest is the input for creating/updating an opportunity.
type CreateOpportunityRequest struct {
	Title        string     `json:"title"`
	Type         string     `json:"type"`
	Description  string     `json:"description"`
	SourceURL    string     `json:"source_url"`
	Organizer    string     `json:"organizer"`
	StartDate    *time.Time `json:"start_date,omitempty"`
	Deadline     *time.Time `json:"deadline,omitempty"`
	EventDate    *time.Time `json:"event_date,omitempty"`
	Location     string     `json:"location"`
	Requirements string     `json:"requirements"`
	TargetMajors []string   `json:"target_majors"`
	Tags         []string   `json:"tags"`
}

// OpportunityFilter defines query filters for listing opportunities.
type OpportunityFilter struct {
	Type             string
	CompetitionLevel string
	Major            string
	Tags             []string
	DeadlineAfter    *time.Time
	DeadlineBefore   *time.Time
	IsActive         *bool
	ProfileID        int64
	Limit            int
	Offset           int
}
