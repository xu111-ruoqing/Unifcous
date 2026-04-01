package seed

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"

	_ "embed"

	"github.com/unifocus/backend/internal/repository/postgres"
)

//go:embed competitions_full.json
var competitionsMasterJSON []byte

type competitionSeedItem struct {
	Name               string   `json:"name"`
	Level              string   `json:"level"`
	Category           string   `json:"category"`
	URL                string   `json:"official_url"`
	TypicalTimeWindow  string   `json:"typical_time_window"`
	TimelineHint       string   `json:"timeline_hint"`
	Note               string   `json:"note"`
	Tags               []string `json:"tags"`
}

func NormalizeNameKey(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return ""
	}

	replacer := strings.NewReplacer(
		"“", "",
		"”", "",
		"‘", "",
		"’", "",
		"`", "",
		"\"", "",
		"'", "",
		"（", "",
		"）", "",
		"(", "",
		")", "",
	)
	name = replacer.Replace(name)

	name = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return -1
		}
		return r
	}, name)

	return strings.ToLower(name)
}

func SeedCompetitions(ctx context.Context, db *postgres.DB) error {
	var list []competitionSeedItem
	if err := json.Unmarshal(competitionsMasterJSON, &list); err != nil {
		return fmt.Errorf("failed to unmarshal competitions seed JSON: %w", err)
	}

	for _, item := range list {
		name := strings.TrimSpace(item.Name)
		if name == "" {
			continue
		}

		officialURL := normalizeURL(item.URL)

		nameKey := NormalizeNameKey(name)
		if nameKey == "" {
			continue
		}

		level := normalizeLevel(item.Level, item.Category)

		_, err := db.ExecContext(ctx, `
			INSERT INTO competitions (name, name_key, level, official_url, typical_time_window, timeline_hint, note, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
			ON CONFLICT (name_key) DO UPDATE
			SET name = EXCLUDED.name,
			    level = EXCLUDED.level,
			    official_url = EXCLUDED.official_url,
			    typical_time_window = EXCLUDED.typical_time_window,
			    timeline_hint = EXCLUDED.timeline_hint,
			    note = EXCLUDED.note,
			    updated_at = NOW()
		`, name, nameKey, level, officialURL,
			strings.TrimSpace(item.TypicalTimeWindow),
			strings.TrimSpace(item.TimelineHint),
			strings.TrimSpace(item.Note),
		)
		if err != nil {
			return fmt.Errorf("seed upsert failed (%s): %w", name, err)
		}
	}

	return nil
}

func normalizeURL(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}
	parts := strings.Split(url, ";")
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			return p
		}
	}
	return ""
}

func normalizeLevel(level, category string) string {
	level = strings.TrimSpace(level)
	category = strings.TrimSpace(category)
	if level == "" {
		return ""
	}
	switch strings.ToUpper(category) {
	case "A", "A*":
		return level + "A类"
	case "B":
		return level + "B类"
	case "":
		return level
	default:
		return level + category
	}
}

