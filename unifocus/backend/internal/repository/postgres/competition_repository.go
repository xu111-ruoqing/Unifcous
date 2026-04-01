package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/unifocus/backend/internal/domain"
)

type CompetitionRepository struct {
	db *DB
}

func NewCompetitionRepository(db *DB) *CompetitionRepository {
	return &CompetitionRepository{db: db}
}

func (r *CompetitionRepository) ListCompetitions(ctx context.Context, q string, level string, limit int, offset int) ([]domain.Competition, error) {
	if limit <= 0 {
		limit = 200
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	var where []string
	var args []any
	argN := 1

	where = append(where, "name_key IS NOT NULL")

	if strings.TrimSpace(q) != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", argN))
		args = append(args, "%"+strings.TrimSpace(q)+"%")
		argN++
	}
	if strings.TrimSpace(level) != "" {
		where = append(where, fmt.Sprintf("level = $%d", argN))
		args = append(args, strings.TrimSpace(level))
		argN++
	}

	args = append(args, limit, offset)
	limitArg := argN
	offsetArg := argN + 1

	query := fmt.Sprintf(`
		SELECT id, name, COALESCE(level, ''), COALESCE(official_url, ''), name_key,
		       COALESCE(typical_time_window, ''), COALESCE(timeline_hint, ''), COALESCE(note, '')
		FROM competitions
		WHERE %s
		ORDER BY id ASC
		LIMIT $%d OFFSET $%d
	`, strings.Join(where, " AND "), limitArg, offsetArg)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []domain.Competition
	for rows.Next() {
		var c domain.Competition
		var nameKey sql.NullString
		if err := rows.Scan(&c.ID, &c.Name, &c.Level, &c.OfficialURL, &nameKey,
			&c.TypicalTimeWindow, &c.TimelineHint, &c.Note); err != nil {
			return nil, err
		}
		if nameKey.Valid {
			c.NameKey = nameKey.String
		}
		out = append(out, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *CompetitionRepository) GetCompetition(ctx context.Context, id int64) (*domain.Competition, error) {
	var c domain.Competition
	var nameKey sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, COALESCE(level, ''), COALESCE(official_url, ''), name_key,
		       COALESCE(typical_time_window, ''), COALESCE(timeline_hint, ''), COALESCE(note, '')
		FROM competitions WHERE id = $1
	`, id).Scan(&c.ID, &c.Name, &c.Level, &c.OfficialURL, &nameKey,
		&c.TypicalTimeWindow, &c.TimelineHint, &c.Note)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if nameKey.Valid {
		c.NameKey = nameKey.String
	}
	return &c, nil
}

func (r *CompetitionRepository) CreateCompetition(ctx context.Context, c *domain.Competition) error {
	c.NameKey = normalizeNameKey(c.Name)
	return r.db.QueryRowContext(ctx, `
		INSERT INTO competitions (name, name_key, level, official_url, typical_time_window, timeline_hint, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`, c.Name, c.NameKey, c.Level, c.OfficialURL, c.TypicalTimeWindow, c.TimelineHint, c.Note).Scan(&c.ID)
}

func (r *CompetitionRepository) UpdateCompetition(ctx context.Context, c *domain.Competition) error {
	c.NameKey = normalizeNameKey(c.Name)
	result, err := r.db.ExecContext(ctx, `
		UPDATE competitions
		SET name=$1, name_key=$2, level=$3, official_url=$4,
		    typical_time_window=$5, timeline_hint=$6, note=$7, updated_at=NOW()
		WHERE id=$8
	`, c.Name, c.NameKey, c.Level, c.OfficialURL, c.TypicalTimeWindow, c.TimelineHint, c.Note, c.ID)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *CompetitionRepository) DeleteCompetition(ctx context.Context, id int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM competitions WHERE id=$1`, id)
	if err != nil {
		return err
	}
	n, _ := result.RowsAffected()
	if n == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func normalizeNameKey(name string) string {
	name = strings.TrimSpace(name)
	replacer := strings.NewReplacer(
		"\u201c", "", "\u201d", "", "\u2018", "", "\u2019", "",
		"`", "", "\"", "", "'", "", "\uff08", "", "\uff09", "", "(", "", ")", "",
	)
	name = replacer.Replace(name)
	var b strings.Builder
	for _, r := range name {
		if r != ' ' && r != '\t' {
			b.WriteRune(r)
		}
	}
	return strings.ToLower(b.String())
}
