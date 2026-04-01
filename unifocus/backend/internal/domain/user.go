package domain

import (
	"time"
)

// User 用户实体
type User struct {
	ID        int64     `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // 不在JSON中返回密码
	School    string    `json:"school" db:"school"`
	Major     string    `json:"major" db:"major"`
	Grade     int       `json:"grade" db:"grade"` // 年级: 1-4
	AvatarURL string    `json:"avatar_url" db:"avatar_url"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfile 用户画像
type UserProfile struct {
	ID           int64     `json:"id" db:"id"`
	UserID       int64     `json:"user_id" db:"user_id"`
	ResumeText   string    `json:"resume_text" db:"resume_text"`
	Skills       []string  `json:"skills" db:"skills"`         // JSONB
	Certificates []Cert    `json:"certificates" db:"certificates"` // JSONB
	Interests    []string  `json:"interests" db:"interests"`   // JSONB
	ResumeVector []float32 `json:"-" db:"resume_vector"`       // 向量表示
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Cert 证书信息
type Cert struct {
	Name  string `json:"name"`
	Score int    `json:"score,omitempty"`
	Date  string `json:"date,omitempty"`
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	School   string `json:"school" binding:"required"`
	Major    string `json:"major" binding:"required"`
	Grade    int    `json:"grade" binding:"required,min=1,max=4"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// UpdateProfileRequest 更新画像请求
type UpdateProfileRequest struct {
	ResumeText   string   `json:"resume_text"`
	Skills       []string `json:"skills"`
	Certificates []Cert   `json:"certificates"`
	Interests    []string `json:"interests"`
}
