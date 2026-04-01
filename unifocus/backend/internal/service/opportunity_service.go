package service

import (
	"context"
	"fmt"

	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
)

// OpportunityService handles opportunity business logic
type OpportunityService struct {
	oppRepo    *postgres.OpportunityRepository
	recService *RecognitionService
}

// NewOpportunityService creates a new opportunity service
func NewOpportunityService(oppRepo *postgres.OpportunityRepository, recService *RecognitionService) *OpportunityService {
	return &OpportunityService{
		oppRepo:    oppRepo,
		recService: recService,
	}
}

// Create creates a new opportunity
func (s *OpportunityService) Create(ctx context.Context, req *domain.CreateOpportunityRequest) (*domain.Opportunity, error) {
	opp := &domain.Opportunity{
		Title:        req.Title,
		Type:         req.Type,
		Description:  req.Description,
		SourceURL:    req.SourceURL,
		SourceType:   "manual", // Manually created
		Organizer:    req.Organizer,
		StartDate:    req.StartDate,
		Deadline:     req.Deadline,
		EventDate:    req.EventDate,
		Location:     req.Location,
		Requirements: req.Requirements,
		TargetMajors: req.TargetMajors,
		Tags:         req.Tags,
		IsActive:     true,
	}

	if err := s.oppRepo.Create(ctx, opp); err != nil {
		return nil, fmt.Errorf("failed to create opportunity: %w", err)
	}

	return opp, nil
}

// GetByID retrieves an opportunity by ID and increments view count
func (s *OpportunityService) GetByID(ctx context.Context, id int64, profileID int64) (*domain.Opportunity, error) {
	opp, err := s.oppRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Enrich with recognition data if applicable
	// Even if profileID is 0, we might want to return default values via ResolveRecognition?
	// The requirement is: "If profileID is not passed, behavior is consistent with before".
	// But before we just returned raw data.
	// If profileID IS passed, we return recognition.
	// If profileID IS NOT passed (0), we can validly return recognition with defaults if we want, or just skip it.
	// Let's call ResolveRecognition anyway if we want consistent struct usage, or only if requested?
	// The user requirement: "If not passed, distinct behavior".
	// But it also says: "return specific level/points and rationale".
	// Let's call it if profileID > 0 OR if we want to populate the new fields.
	// Since `Recognition` is `omitempty`, we only populate it if useful.
	if profileID > 0 {
		rec, err := s.recService.ResolveRecognition(ctx, opp, profileID)
		if err == nil {
			opp.Recognition = rec
		}
		// Log error?
	} else {
		// Fill with defaults? The prompt said "If no profile_id, behavior consistent with now (which is just returning the opp)".
		// But new FE might expect `recognition` field anyway?
		// "不传 profile_id 时行为与现在一致" -> implies old fields used.
		// "传 profile_id 时能返回专业差异...".
		// So if profileID == 0, we leave `Recognition` nil.
	}

	// Increment view count asynchronously (fire and forget)
	go func() {
		_ = s.oppRepo.IncrementViewCount(context.Background(), id)
	}()

	return opp, nil
}

// List retrieves opportunities with filtering and pagination
func (s *OpportunityService) List(ctx context.Context, filter *domain.OpportunityFilter) ([]*domain.Opportunity, int64, error) {
	if filter == nil {
		filter = &domain.OpportunityFilter{
			Limit:  20,
			Offset: 0,
		}
	}

	opportunities, total, err := s.oppRepo.List(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list opportunities: %w", err)
	}

	// Enrich with recognition data if profile_id is provided
	if filter.ProfileID > 0 {
		// Optimization: Could batch this or do it in SQL if possible, but DB rules are complex.
		// For now, iterate.
		for _, opp := range opportunities {
			rec, err := s.recService.ResolveRecognition(ctx, opp, filter.ProfileID)
			if err == nil {
				opp.Recognition = rec
			}
		}
	}

	return opportunities, total, nil
}

// Update updates an opportunity
func (s *OpportunityService) Update(ctx context.Context, id int64, req *domain.CreateOpportunityRequest) (*domain.Opportunity, error) {
	// Get existing opportunity
	opp, err := s.oppRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields
	opp.Title = req.Title
	opp.Type = req.Type
	opp.Description = req.Description
	opp.SourceURL = req.SourceURL
	opp.Organizer = req.Organizer
	opp.StartDate = req.StartDate
	opp.Deadline = req.Deadline
	opp.EventDate = req.EventDate
	opp.Location = req.Location
	opp.Requirements = req.Requirements
	opp.TargetMajors = req.TargetMajors
	opp.Tags = req.Tags

	if err := s.oppRepo.Update(ctx, opp); err != nil {
		return nil, fmt.Errorf("failed to update opportunity: %w", err)
	}

	return opp, nil
}

// Delete soft deletes an opportunity
func (s *OpportunityService) Delete(ctx context.Context, id int64) error {
	return s.oppRepo.Delete(ctx, id)
}

// IncrementSaveCount increments the save count for an opportunity
func (s *OpportunityService) IncrementSaveCount(ctx context.Context, id int64) error {
	return s.oppRepo.IncrementSaveCount(ctx, id)
}
