package domain

import "time"

type Task struct {
	ID            uint      `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Task          string    `gorm:"type:varchar(255)" json:"task"`
	Summary       string    `gorm:"type:text" json:"summary"`
	PerformedDate time.Time `gorm:"type:datetime" json:"performed_date"`

	UserID uint  `gorm:"not null" json:"user_id,omitempty"`
	User   *User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE;" json:"user"`

	CreatedAt time.Time  `gorm:"type:timestamp;autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp;autoUpdateTime;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}
