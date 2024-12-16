package domain

import "time"

type Role struct {
	ID        uint       `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Role      string     `gorm:"type:varchar(255)" json:"role"`
	Slug      string     `gorm:"type:varchar(255)" json:"slug"`
	CreatedAt time.Time  `gorm:"type:timestamp;autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp;autoUpdateTime;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}
