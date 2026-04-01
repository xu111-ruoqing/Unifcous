package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/unifocus/backend/internal/config"
	"github.com/unifocus/backend/internal/repository/postgres"
	"github.com/unifocus/backend/pkg/logger"
)

type SeedCompetition struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Level    string `json:"level"`
	Category string `json:"category"`
	URL      string `json:"url"`
}

func main() {
	var filePath string
	var configPath string
	flag.StringVar(&filePath, "file", "data/competitions_master.json", "path to competitions JSON file (array)")
	flag.StringVar(&configPath, "config", "configs/config.dev.yaml", "backend config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}
	defer logger.Sync()

	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}
	defer db.Close()

	payload, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read %s: %v", filePath, err)
	}

	var list []SeedCompetition
	if err := json.Unmarshal(payload, &list); err != nil {
		log.Fatalf("Invalid JSON array in %s: %v", filePath, err)
	}

	ctx := context.Background()

	var upserted int
	var skipped int

	for _, c := range list {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			skipped++
			continue
		}

		officialURL := normalizeURL(c.URL)
		level := normalizeLevel(c.Level, c.Category)

		officialURLArg := nullIfEmpty(officialURL) // prefer NULL when missing
		_, err := db.DB.ExecContext(ctx, `
			INSERT INTO competitions (name, official_url, level, created_at, updated_at)
			VALUES ($1, $2, $3, NOW(), NOW())
			ON CONFLICT (name) DO UPDATE SET
				official_url = EXCLUDED.official_url,
				level = EXCLUDED.level,
				updated_at = NOW()
		`, name, officialURLArg, nullIfEmpty(level))
		if err != nil && strings.Contains(err.Error(), "official_url") && strings.Contains(err.Error(), "not-null") {
			_, err = db.DB.ExecContext(ctx, `
				INSERT INTO competitions (name, official_url, level, created_at, updated_at)
				VALUES ($1, $2, $3, NOW(), NOW())
				ON CONFLICT (name) DO UPDATE SET
					official_url = EXCLUDED.official_url,
					level = EXCLUDED.level,
					updated_at = NOW()
			`, name, emptyStringIfMissing(officialURL), nullIfEmpty(level))
		}
		if err != nil {
			log.Fatalf("Upsert failed (%s): %v", name, err)
		}
		upserted++
	}

	fmt.Printf("✓ competitions 导入完成：upsert=%d, skipped=%d\n", upserted, skipped)
	fmt.Println(`验证：psql ... -c "SELECT id,name,level,official_url FROM competitions ORDER BY id;"`)
}

func normalizeLevel(level, category string) string {
	level = strings.TrimSpace(level)
	category = strings.TrimSpace(category)
	if level == "" {
		return ""
	}
	if category == "" {
		return level
	}
	switch strings.ToUpper(category) {
	case "A", "A*":
		return level + "A类"
	case "B":
		return level + "B类"
	default:
		return level + category
	}
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

func nullIfEmpty(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	return s
}

func emptyStringIfMissing(s string) string {
	if strings.TrimSpace(s) == "" {
		return ""
	}
	return s
}
