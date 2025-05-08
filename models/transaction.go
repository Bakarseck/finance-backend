package models

import "time"

type Transaction struct {
    ID         uint      `gorm:"primaryKey" json:"id"`
    Icon       string    `json:"icon"`
    Title      string    `json:"title"`
    Subtitle   string    `json:"subtitle"`
    Date       time.Time `json:"date"`
    Amount     float64   `json:"amount"`
    Type       string    `json:"type"`        
    Category   string    `json:"category"`   	
    CategoryID uint      `json:"category_id"` 
    UserID     uint      `json:"user_id"`
}
