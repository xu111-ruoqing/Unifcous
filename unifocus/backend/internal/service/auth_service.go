package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/unifocus/backend/internal/domain"
	"github.com/unifocus/backend/internal/repository/postgres"
	"github.com/unifocus/backend/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication business logic
type AuthService struct {
	userRepo *postgres.UserRepository
	jwtMgr   *jwt.Manager
}

// NewAuthService creates a new authentication service
func NewAuthService(userRepo *postgres.UserRepository, jwtMgr *jwt.Manager) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		jwtMgr:   jwtMgr,
	}
}

// Register registers a new user
func (s *AuthService) Register(ctx context.Context, req *domain.CreateUserRequest) (*domain.User, string, error) {
	// Check if email already exists
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check email existence: %w", err)
	}
	if exists {
		return nil, "", errors.New("email already exists")
	}

	// Check if username already exists
	exists, err = s.userRepo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, "", fmt.Errorf("failed to check username existence: %w", err)
	}
	if exists {
		return nil, "", errors.New("username already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	// Create user
	user := &domain.User{
		Username:  req.Username,
		Email:     req.Email,
		Password:  string(hashedPassword),
		School:    req.School,
		Major:     req.Major,
		Grade:     req.Grade,
		AvatarURL: "",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := s.jwtMgr.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// Login authenticates a user and returns a JWT token
func (s *AuthService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.User, string, error) {
	// Get user by email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := s.jwtMgr.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	return user, token, nil
}

// ValidateToken validates a JWT token and returns the user
func (s *AuthService) ValidateToken(ctx context.Context, tokenString string) (*domain.User, error) {
	claims, err := s.jwtMgr.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}

// RefreshToken generates a new token from an existing valid token
func (s *AuthService) RefreshToken(ctx context.Context, tokenString string) (string, error) {
	claims, err := s.jwtMgr.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// Verify user still exists
	_, err = s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Generate new token
	return s.jwtMgr.RefreshToken(tokenString)
}
