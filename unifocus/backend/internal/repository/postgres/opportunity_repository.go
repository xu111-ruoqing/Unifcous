package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/unifocus/backend/internal/domain"
)

// OpportunityRepository handles opportunity data access operations
type OpportunityRepository struct {
	db *DB
}

// NewOpportunityRepository creates a new opportunity repository
func NewOpportunityRepository(db *DB) *OpportunityRepository {
	return &OpportunityRepository{db: db}
}

// UpsertCompetition inserts/updates a competition master record
func (r *OpportunityRepository) UpsertCompetition(ctx context.Context, name, url, level string) error {
	query := `
		INSERT INTO competitions (name, official_url, level, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		ON CONFLICT (official_url) DO UPDATE
		SET name = EXCLUDED.name, level = EXCLUDED.level, updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, name, url, level)
	return err
}

// CreateMinimal inserts core fields only, leaving extended JSON columns NULL/default.
func (r *OpportunityRepository) CreateMinimal(ctx context.Context, opp *domain.Opportunity) (int64, error) {
	if opp.NoticeType == "" {
		opp.NoticeType = "notice"
	}
	query := `
		INSERT INTO opportunities (
			title, type, source_url, notice_type, competition_level,
			registration_start_at, registration_deadline, event_start_at, event_end_at, raw_text,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	if err := r.db.QueryRowContext(ctx, query,
		opp.Title,
		opp.Type,
		opp.SourceURL,
		opp.NoticeType,
		opp.CompetitionLevel,
		opp.RegistrationStartAt,
		opp.RegistrationDeadline,
		opp.EventStartAt,
		opp.EventEndAt,
		sql.NullString{String: opp.RawText, Valid: opp.RawText != ""},
	).Scan(&opp.ID, &opp.CreatedAt, &opp.UpdatedAt); err != nil {
		return 0, err
	}
	return opp.ID, nil
}

// Create creates a new opportunity
func (r *OpportunityRepository) Create(ctx context.Context, opp *domain.Opportunity) error {
	query := `
		INSERT INTO opportunities (
			title, type, description, source_url, source_type,
			competition_level, certification_type, organizer, organizer_type, award_level, default_points_value, is_official,
			start_date, deadline, event_date, event_start_at, event_end_at, published_at, location, time_info_raw, location_raw,
			notice_type, registration_start_at, registration_deadline, raw_text,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32)
		RETURNING id, created_at, updated_at
	`

	attachmentsJSON, err := json.Marshal(opp.Attachments)
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		opp.Title,
		opp.Type,
		opp.Description,
		opp.SourceURL,
		opp.SourceType,
		opp.CompetitionLevel,
		opp.CertificationType,
		opp.Organizer,
		opp.OrganizerType,
		opp.AwardLevel,
		opp.DefaultPointsValue,
		opp.IsOfficial,
		opp.StartDate,
		opp.Deadline,
		opp.EventDate,
		opp.EventStartAt,
		opp.EventEndAt,
		opp.PublishedAt,
		opp.Location,
		opp.TimeInfoRaw,
		opp.LocationRaw,
		opp.NoticeType,
		opp.RegistrationStartAt,
		opp.RegistrationDeadline,
		sql.NullString{String: opp.RawText, Valid: opp.RawText != ""},
		opp.Requirements,
		opp.EligibilityRules,
		pq.Array(opp.TargetMajors),
		pq.Array(opp.Tags),
		attachmentsJSON,
		pq.Array(opp.DescriptionVector),
		opp.IsActive,
	).Scan(&opp.ID, &opp.CreatedAt, &opp.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByID retrieves an opportunity by ID
func (r *OpportunityRepository) GetByID(ctx context.Context, id int64) (*domain.Opportunity, error) {
	query := `
		SELECT id, title, type, COALESCE(description, '') AS description, source_url, COALESCE(source_type, '') AS source_type,
			COALESCE(competition_level, '') AS competition_level, COALESCE(certification_type, '') AS certification_type,
			COALESCE(organizer, '') AS organizer, COALESCE(organizer_type, '') AS organizer_type, COALESCE(award_level, '') AS award_level,
			default_points_value, is_official,
			start_date, deadline, event_date, event_start_at, event_end_at, published_at,
			COALESCE(location, '') AS location, COALESCE(time_info_raw, '') AS time_info_raw, COALESCE(location_raw, '') AS location_raw,
			COALESCE(notice_type, '') AS notice_type, registration_start_at, registration_deadline, COALESCE(raw_text, '') AS raw_text,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active, view_count, save_count,
			created_at, updated_at
		FROM opportunities
		WHERE id = $1
	`

	opp := &domain.Opportunity{}
	var targetMajors, tags []string
	var descriptionVector []float32
	var attachmentsJSON []byte
	var noticeType sql.NullString
	var regStart, regDeadline sql.NullTime
	var rawText sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&opp.ID,
		&opp.Title,
		&opp.Type,
		&opp.Description,
		&opp.SourceURL,
		&opp.SourceType,
		&opp.CompetitionLevel,
		&opp.CertificationType,
		&opp.Organizer,
		&opp.OrganizerType,
		&opp.AwardLevel,
		&opp.DefaultPointsValue,
		&opp.IsOfficial,
		&opp.StartDate,
		&opp.Deadline,
		&opp.EventDate,
		&opp.EventStartAt,
		&opp.EventEndAt,
		&opp.PublishedAt,
		&opp.Location,
		&opp.TimeInfoRaw,
		&opp.LocationRaw,
		&noticeType,
		&regStart,
		&regDeadline,
		&rawText,
		&opp.Requirements,
		&opp.EligibilityRules,
		pq.Array(&targetMajors),
		pq.Array(&tags),
		&attachmentsJSON,
		pq.Array(&descriptionVector),
		&opp.IsActive,
		&opp.ViewCount,
		&opp.SaveCount,
		&opp.CreatedAt,
		&opp.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("opportunity not found")
		}
		return nil, err
	}

	opp.TargetMajors = targetMajors
	opp.Tags = tags
	opp.DescriptionVector = descriptionVector
	if noticeType.Valid {
		opp.NoticeType = noticeType.String
	}
	if regStart.Valid {
		t := regStart.Time
		opp.RegistrationStartAt = &t
	}
	if regDeadline.Valid {
		t := regDeadline.Time
		opp.RegistrationDeadline = &t
	}
	if rawText.Valid {
		opp.RawText = rawText.String
	}
	if len(attachmentsJSON) > 0 {
		if err := json.Unmarshal(attachmentsJSON, &opp.Attachments); err != nil {
			return nil, fmt.Errorf("failed to unmarshal attachments: %w", err)
		}
	}

	return opp, nil
}

// List retrieves opportunities with filtering and pagination
func (r *OpportunityRepository) List(ctx context.Context, filter *domain.OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	var conditions []string
	var args []interface{}
	argPos := 1

	// Build WHERE clause
	if filter.Type != "" {
		conditions = append(conditions, fmt.Sprintf("type = $%d", argPos))
		args = append(args, filter.Type)
		argPos++
	}

	if filter.CompetitionLevel != "" {
		conditions = append(conditions, fmt.Sprintf("competition_level = $%d", argPos))
		args = append(args, filter.CompetitionLevel)
		argPos++
	}

	if filter.Major != "" {
		conditions = append(conditions, fmt.Sprintf("$%d = ANY(target_majors)", argPos))
		args = append(args, filter.Major)
		argPos++
	}

	if filter.DeadlineAfter != nil {
		conditions = append(conditions, fmt.Sprintf("deadline >= $%d", argPos))
		args = append(args, *filter.DeadlineAfter)
		argPos++
	}

	if filter.DeadlineBefore != nil {
		conditions = append(conditions, fmt.Sprintf("deadline <= $%d", argPos))
		args = append(args, *filter.DeadlineBefore)
		argPos++
	}

	if filter.IsActive != nil {
		conditions = append(conditions, fmt.Sprintf("is_active = $%d", argPos))
		args = append(args, *filter.IsActive)
		argPos++
	}

	if len(filter.Tags) > 0 {
		conditions = append(conditions, fmt.Sprintf("tags && $%d", argPos))
		args = append(args, pq.Array(filter.Tags))
		argPos++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Count total
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM opportunities %s", whereClause)
	var total int64
	err := r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Set default pagination
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	// Select opportunities
	query := fmt.Sprintf(`
		SELECT id, title, type, COALESCE(description, '') AS description, source_url, COALESCE(source_type, '') AS source_type,
			COALESCE(competition_level, '') AS competition_level, COALESCE(certification_type, '') AS certification_type,
			COALESCE(organizer, '') AS organizer, COALESCE(organizer_type, '') AS organizer_type, COALESCE(award_level, '') AS award_level, default_points_value, is_official,
			start_date, deadline, event_date, event_start_at, event_end_at, published_at, COALESCE(location, '') AS location, COALESCE(time_info_raw, '') AS time_info_raw, COALESCE(location_raw, '') AS location_raw,
			COALESCE(notice_type, '') AS notice_type, registration_start_at, registration_deadline, COALESCE(raw_text, '') AS raw_text,
			requirements, eligibility_rules, target_majors,
			tags, attachments, description_vector, is_active, view_count, save_count,
			created_at, updated_at
		FROM opportunities
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argPos, argPos+1)

	args = append(args, filter.Limit, filter.Offset)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var opportunities []*domain.Opportunity
	for rows.Next() {
		opp := &domain.Opportunity{}
		var targetMajors, tags []string
		var descriptionVector []float32
		var attachmentsJSON []byte
		var noticeType sql.NullString
		var regStart, regDeadline sql.NullTime
		var rawText sql.NullString

		err := rows.Scan(
			&opp.ID,
			&opp.Title,
			&opp.Type,
			&opp.Description,
			&opp.SourceURL,
			&opp.SourceType,
			&opp.CompetitionLevel,
			&opp.CertificationType,
			&opp.Organizer,
			&opp.OrganizerType,
			&opp.AwardLevel,
			&opp.DefaultPointsValue,
			&opp.IsOfficial,
			&opp.StartDate,
			&opp.Deadline,
			&opp.EventDate,
			&opp.EventStartAt,
			&opp.EventEndAt,
			&opp.PublishedAt,
			&opp.Location,
			&opp.TimeInfoRaw,
			&opp.LocationRaw,
			&noticeType,
			&regStart,
			&regDeadline,
			&rawText,
			&opp.Requirements,
			&opp.EligibilityRules,
			pq.Array(&targetMajors),
			pq.Array(&tags),
			&attachmentsJSON,
			pq.Array(&descriptionVector),
			&opp.IsActive,
			&opp.ViewCount,
			&opp.SaveCount,
			&opp.CreatedAt,
			&opp.UpdatedAt,
		)

		if err != nil {
			return nil, 0, err
		}

		opp.TargetMajors = targetMajors
		opp.Tags = tags
		opp.DescriptionVector = descriptionVector
		if noticeType.Valid {
			opp.NoticeType = noticeType.String
		}
		if regStart.Valid {
			t := regStart.Time
			opp.RegistrationStartAt = &t
		}
		if regDeadline.Valid {
			t := regDeadline.Time
			opp.RegistrationDeadline = &t
		}
		if rawText.Valid {
			opp.RawText = rawText.String
		}
		if len(attachmentsJSON) > 0 {
			if err := json.Unmarshal(attachmentsJSON, &opp.Attachments); err != nil {
				return nil, 0, fmt.Errorf("failed to unmarshal attachments: %w", err)
			}
		}

		opportunities = append(opportunities, opp)
	}

	return opportunities, total, nil
}

// GetAll retrieves opportunities with pagination (no filters)
func (r *OpportunityRepository) GetAll(ctx context.Context, limit, offset int) ([]*domain.Opportunity, error) {
	filter := &domain.OpportunityFilter{
		Limit:  limit,
		Offset: offset,
	}

	opportunities, _, err := r.List(ctx, filter)
	return opportunities, err
}

// Update updates an opportunity
func (r *OpportunityRepository) Update(ctx context.Context, opp *domain.Opportunity) error {
	query := `
		UPDATE opportunities
		SET title = $1, type = $2, description = $3, source_url = $4, source_type = $5,
			competition_level = $6, certification_type = $7, organizer = $8, organizer_type = $9, 
			award_level = $10, default_points_value = $11, is_official = $12,
			start_date = $13, deadline = $14, event_date = $15, event_start_at = $16, event_end_at = $17, published_at = $18, location = $19, time_info_raw = $20, location_raw = $21,
			notice_type = $22, registration_start_at = $23, registration_deadline = $24, raw_text = $25,
			requirements = $26, eligibility_rules = $27, target_majors = $28,
			tags = $29, attachments = $30, description_vector = $31, is_active = $32,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $33
		RETURNING updated_at
	`

	attachmentsJSON, err := json.Marshal(opp.Attachments)
	if err != nil {
		return fmt.Errorf("failed to marshal attachments: %w", err)
	}

	err = r.db.QueryRowContext(ctx, query,
		opp.Title,
		opp.Type,
		opp.Description,
		opp.SourceURL,
		opp.SourceType,
		opp.CompetitionLevel,
		opp.CertificationType,
		opp.Organizer,
		opp.OrganizerType,
		opp.AwardLevel,
		opp.DefaultPointsValue,
		opp.IsOfficial,
		opp.StartDate,
		opp.Deadline,
		opp.EventDate,
		opp.EventStartAt,
		opp.EventEndAt,
		opp.PublishedAt,
		opp.Location,
		opp.TimeInfoRaw,
		opp.LocationRaw,
		opp.NoticeType,
		opp.RegistrationStartAt,
		opp.RegistrationDeadline,
		sql.NullString{String: opp.RawText, Valid: opp.RawText != ""},
		opp.Requirements,
		opp.EligibilityRules,
		pq.Array(opp.TargetMajors),
		pq.Array(opp.Tags),
		attachmentsJSON,
		pq.Array(opp.DescriptionVector),
		opp.IsActive,
		opp.ID,
	).Scan(&opp.UpdatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("opportunity not found")
		}
		return err
	}

	return nil
}

// Delete soft deletes an opportunity (sets is_active = false)
func (r *OpportunityRepository) Delete(ctx context.Context, id int64) error {
	query := `
		UPDATE opportunities
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("opportunity not found")
	}

	return nil
}

// IncrementViewCount increments the view count for an opportunity
func (r *OpportunityRepository) IncrementViewCount(ctx context.Context, id int64) error {
	query := `UPDATE opportunities SET view_count = view_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// IncrementSaveCount increments the save count for an opportunity
func (r *OpportunityRepository) IncrementSaveCount(ctx context.Context, id int64) error {
	query := `UPDATE opportunities SET save_count = save_count + 1 WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

// ExistsBySourceURL checks if an opportunity with the same source_url already exists
func (r *OpportunityRepository) ExistsBySourceURL(ctx context.Context, sourceURL string) (bool, error) {
	query := `SELECT EXISTS (SELECT 1 FROM opportunities WHERE source_url = $1)`
	var exists bool
	if err := r.db.QueryRowContext(ctx, query, sourceURL).Scan(&exists); err != nil {
		return false, err
	}
	return exists, nil
}
