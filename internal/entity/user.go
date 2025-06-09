package entity

type User struct {
	ID        string `json:"id" gorm:"primaryKey;type:uuid;default:uuid_generate_v4()"`
	Email     string `json:"email" gorm:"uniqueIndex;not null"`
	Password  string `json:"password" gorm:"not null"`
	Role      string `json:"role" gorm:"not null"`
	Status    int    `json:"status" gorm:"not null;default:'active'"`
	CreatedAt string `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt string `json:"updated_at" gorm:"autoUpdateTime;not null"`
}
