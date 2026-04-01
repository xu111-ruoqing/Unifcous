package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"

	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
)

// ProfileService handles user profile business logic
type ProfileService struct {
	profileRepo *postgres.ProfileRepository
	nlpClient   NLPClient // NLP服务客户端（待实现）
}

// NLPClient NLP服务客户端接口
type NLPClient interface {
	ExtractTextFromPDF(ctx context.Context, fileData []byte) (string, error)
	VectorizeText(ctx context.Context, text string) ([]float32, error)
	ExtractSkills(ctx context.Context, text string) ([]string, error)
}

// NewProfileService creates a new profile service
func NewProfileService(profileRepo *postgres.ProfileRepository, nlpClient NLPClient) *ProfileService {
	return &ProfileService{
		profileRepo: profileRepo,
		nlpClient:   nlpClient,
	}
}

// GetProfile retrieves user profile
func (s *ProfileService) GetProfile(ctx context.Context, userID int64) (*domain.UserProfile, error) {
	profile, err := s.profileRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}

	return profile, nil
}

// UpdateProfile updates user profile
func (s *ProfileService) UpdateProfile(ctx context.Context, userID int64, req *domain.UpdateProfileRequest) (*domain.UserProfile, error) {
	profile := &domain.UserProfile{
		UserID:       userID,
		ResumeText:   req.ResumeText,
		Skills:       req.Skills,
		Certificates: req.Certificates,
		Interests:    req.Interests,
	}

	// 如果有简历文本，进行向量化
	if req.ResumeText != "" && s.nlpClient != nil {
		vector, err := s.nlpClient.VectorizeText(ctx, req.ResumeText)
		if err == nil {
			profile.ResumeVector = vector
		}
	}

	if err := s.profileRepo.CreateOrUpdate(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to update profile: %w", err)
	}

	return profile, nil
}

// UploadResume uploads and parses a resume file
func (s *ProfileService) UploadResume(ctx context.Context, userID int64, file multipart.File, filename string) (*domain.UserProfile, error) {
	// 读取文件内容
	fileData, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 提取文本（根据文件类型）
	var text string
	if s.nlpClient != nil {
		// 假设是PDF文件
		text, err = s.nlpClient.ExtractTextFromPDF(ctx, fileData)
		if err != nil {
			return nil, fmt.Errorf("failed to extract text: %w", err)
		}
	} else {
		// 如果没有NLP客户端，返回错误
		return nil, fmt.Errorf("NLP service not available")
	}

	// 提取技能
	var skills []string
	if s.nlpClient != nil {
		skills, err = s.nlpClient.ExtractSkills(ctx, text)
		if err != nil {
			// 技能提取失败不影响整体流程
			skills = []string{}
		}
	}

	// 更新或创建profile
	profile := &domain.UserProfile{
		UserID:     userID,
		ResumeText: text,
		Skills:     skills,
	}

	// 向量化简历文本
	if s.nlpClient != nil {
		vector, err := s.nlpClient.VectorizeText(ctx, text)
		if err == nil {
			profile.ResumeVector = vector
		}
	}

	if err := s.profileRepo.CreateOrUpdate(ctx, profile); err != nil {
		return nil, fmt.Errorf("failed to save profile: %w", err)
	}

	return profile, nil
}
