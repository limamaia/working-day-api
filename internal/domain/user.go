package domain

import "time"

type User struct {
	ID       uint   `gorm:"primaryKey;autoIncrement:true" json:"id"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Email    string `gorm:"type:varchar(320);uniqueIndex" json:"email"`
	Password string `gorm:"type:varchar(65)" json:"password"`

	RoleID *uint `gorm:"default:NULL" json:"role_id,omitempty"`
	Role   *Role `gorm:"foreignKey:RoleID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"role"`

	CreatedAt time.Time  `gorm:"type:timestamp;autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt *time.Time `gorm:"type:timestamp;autoUpdateTime;default:CURRENT_TIMESTAMP;ON UPDATE CURRENT_TIMESTAMP;" json:"updated_at"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=255"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	RoleID   uint   `json:"role_id" binding:"required"`
}
