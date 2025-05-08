package models

type Category struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    Name     string `json:"name"`       
    Type     string `json:"type"`       
    Icon     string `json:"icon"`       
    UserID   uint   `json:"user_id"`    
}
