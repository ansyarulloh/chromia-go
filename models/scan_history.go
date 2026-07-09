package models

import (
    "time"

    "github.com/google/uuid"
    "gorm.io/gorm"
)

type ScanHistory struct {
    ID        uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
    UserID    string    `gorm:"index;not null" json:"user_id"`
    ColorName string    `json:"color_name"`
    HexCode   string    `json:"hex_code"`
    CreatedAt time.Time `json:"created_at"`
}

func (h *ScanHistory) BeforeCreate(tx *gorm.DB) (err error) {
    h.ID = uuid.New()
    h.CreatedAt = time.Now()
    return
}
