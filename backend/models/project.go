package models

import (
	"time"

	"gorm.io/gorm"
)

type Project struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null"`
	Description string         `json:"description"`
	UserID      uint           `json:"user_id" gorm:"not null"`
	Status      string         `json:"status" gorm:"default:pending"` // pending, analyzing, completed, failed
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// リレーション
	User     User       `json:"user" gorm:"foreignKey:UserID"`
	Files    []File     `json:"files" gorm:"foreignKey:ProjectID"`
	Analysis []Analysis `json:"analysis" gorm:"foreignKey:ProjectID"`
}

type File struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProjectID uint           `json:"project_id" gorm:"not null"`
	Name      string         `json:"name" gorm:"not null"`
	Path      string         `json:"path" gorm:"not null"`
	Size      int64          `json:"size"`
	MimeType  string         `json:"mime_type"`
	Content   string         `json:"content,omitempty" gorm:"type:text"`
	Language  string         `json:"language"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// リレーション
	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
}

type Analysis struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	ProjectID uint           `json:"project_id" gorm:"not null"`
	FileID    *uint          `json:"file_id,omitempty"`
	Type      string         `json:"type" gorm:"not null"`          // code_analysis, dependency_map, documentation, pattern_detection
	Status    string         `json:"status" gorm:"default:pending"` // pending, processing, completed, failed
	Result    string         `json:"result,omitempty" gorm:"type:text"`
	Metadata  string         `json:"metadata,omitempty" gorm:"type:json"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// リレーション
	Project Project `json:"project" gorm:"foreignKey:ProjectID"`
	File    *File   `json:"file,omitempty" gorm:"foreignKey:FileID"`
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null"`
	Name      string         `json:"name" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`

	// リレーション
	Projects []Project `json:"projects" gorm:"foreignKey:UserID"`
}
