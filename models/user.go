package models

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
    FullName string `json:"full_name"`
    Email    string `gorm:"unique" json:"email"`
    Currency string `json:"currency"`
    Password string `json:"password"`
}
