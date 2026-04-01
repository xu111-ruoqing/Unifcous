package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/unifocus/backend/internal/domain"
)

// Mocking required types since we can't import internal/domain if we run this as standalone easily without module setup
// But we are in the project, so `go run` should work if we put this in `cmd/demo/main.go` and use proper imports.

func main() {
	deadline := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

	// Construct a mock Opportunity with all new fields
	opp := domain.Opportunity{
		ID:                 101,
		Title:              "2024 National AI Challenge",
		Type:               "competition",
		Description:        "A top-tier AI competition for students.",
		SourceURL:          "https://example.com/ai-challenge",
		SourceType:         "crawler",
		CompetitionLevel:   "National", // Original Level
		AwardLevel:         "First Prize",
		DefaultPointsValue: 50, // Renamed Field
		IsOfficial:         true,
		StartDate:          &deadline,
		Deadline:           &deadline,
		Requirements: domain.Requirements{
			Major:  []string{"Computer Science", "Software Engineering"},
			Skills: []string{"Python", "PyTorch"},
			Grade:  []int{2021, 2022},
		},
		TargetMajors: []string{"CS", "SE"},
		Tags:         []string{"AI", "Official", "HighValue"},
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),

		// New Dynamic Recognition Field
		Recognition: &domain.RecognitionResult{
			Level:      "National A-Class",
			Grade:      "A",
			Points:     100, // Enriched value based on Profile
			Confidence: 0.95,
			Rationale: map[string]interface{}{
				"rule_id": 5,
				"reason":  "Matches 'AI Challenge' in CS whitelist",
			},
		},
	}

	// Marshal to JSON to show structure
	bytes, err := json.MarshalIndent(opp, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("=== 数据库模型 (Domain Object) 样例输出 ===")
	fmt.Println(string(bytes))
}
