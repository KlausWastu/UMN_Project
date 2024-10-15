package models

type User struct {
	UserID   int64  `gorm:"primaryKey" json:"userID"`
	Fullname string `gorm:"varchar(300)" json:"fullname"`
	Username string `gorm:"varchar(300)" json:"username"`
	Email    string `gorm:"varchar(255)" json:"email"`
	Password string `gorm:"varchar(300)" json:"password"`
}
