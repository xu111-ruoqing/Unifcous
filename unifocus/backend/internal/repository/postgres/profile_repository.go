package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/unifocus/backend/internal/domain"
)

// ProfileRepository handles user profile data access operations
type ProfileRepository struct {
	db *DB
}

// NewProfileRepository creates a new profile repository
func NewProfileRepository(db *DB) *ProfileRepository {
	return &ProfileRepository{db: db}
}

// CreateOrUpdate creates or updates a user profile
func (r *ProfileRepository) CreateOrUpdate(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (user_id, resume_text, skills, certificates, interests, resume_vector)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id) 
		DO UPDATE SET 
			resume_text = EXCLUDED.resume_text,
			skills = EXCLUDED.skills,
			certificates = EXCLUDED.certificates,
			interests = EXCLUDED.interests,
			resume_vector = EXCLUDED.resume_vector,
			updated_at = CURRENT_TIMESTAMP
		RETURNING id, updated_at
	`

	err := r.db.QueryRowContext(ctx, query,
		profile.UserID,
		profile.ResumeText,
		pq.Array(profile.Skills),
		pq.Array(profile.Certificates),
		pq.Array(profile.Interests),
		pq.Array(profile.ResumeVector),
	).Scan(&profile.ID, &profile.UpdatedAt)

	if err != nil {
		return err
	}

	return nil
}

// GetByUserID retrieves a user profile by user ID
func (r *ProfileRepository) GetByUserID(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	query := `
		SELECT id, user_id, resume_text, skills, certificates, interests, resume_vector, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	profile := &domain.UserProfile{}
	var skills, interests []string
	var certificates []domain.Cert
	var resumeVector []float32

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.ResumeText,
		pq.Array(&skills),
		pq.Array(&certificates),
		pq.Array(&interests),
		pq.Array(&resumeVector),
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("profile not found")
		}
		return nil, err
	}

	profile.Skills = skills
	profile.Certificates = certificates
	profile.Interests = interests
	profile.ResumeVector = resumeVector

	return profile, nil
}

// UpdateSkills updates user skills
func (r *ProfileRepository) UpdateSkills(ctx context.Context, userID int64, skills []string) error {
	query := `
		UPDATE user_profiles
		SET skills = $1, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $2
	`

	result, err := r.db.ExecContext(ctx, query, pq.Array(skills), userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("profile not found")
	}

	return nil
}
