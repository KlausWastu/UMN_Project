package models

type TemporaryData struct {
	No    int64  `gorm:"primaryKey" json:"no"`
	Name  string `gorm:"varchar(255)" json:"name"`
	Score int8   `gorm:"integer" json:"score"`
}
