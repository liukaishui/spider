package models

import (
	"time"
)

type Sites struct {
	ID                int       `gorm:"column:id"`
	Name              string    `gorm:"column:name"`
	Url               string    `gorm:"column:url"`
	Status            string    `gorm:"column:status"`
	LastExecutionTime time.Time `gorm:"column:last_execution_time"`
	CreatedAt         time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt         time.Time `gorm:"column:updated_at;autoCreateTime"`
}
