package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type User struct {
    ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    Name           string    `json:"name"`
    Email          string    `gorm:"unique;not null" json:"email"`
    
    // INI YANG DITAMBAHIN:
    // json:"-" itu penting banget! Biar password nggak pernah ikut ke-kirim ke Frontend (Flutter) pas kita narik data user.
    Password       string    `gorm:"not null" json:"-"` 
    
    AvatarURL      string    `json:"avatar_url"`
    IsPremium      bool      `gorm:"default:false" json:"is_premium"`
    DailyScanCount int       `gorm:"default:0" json:"daily_scan_count"`
    LastScanDate   time.Time `json:"last_scan_date"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      time.Time `json:"updated_at"`
}

// BeforeCreate adalah 'hook' GORM yang otomatis berjalan tepat sebelum data user disimpan ke DB.
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
    u.ID = uuid.New()
    u.LastScanDate = time.Now()
    return
}