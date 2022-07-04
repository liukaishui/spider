package models

import (
	"time"
)

type Urls struct {
	ID                int       `gorm:"column:id"`
	SiteID            int       `gorm:"column:site_id"`
	Url               string    `gorm:"column:url"`
	Info              string    `gorm:"column:info"`
	Content           string    `gorm:"column:content"`
	StatusCode        string    `gorm:"column:status_code"`
	LastExecutionTime time.Time `gorm:"column:last_execution_time"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoCreateTime"`
}
