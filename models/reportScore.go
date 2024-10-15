package models

type ReportScore struct {
	Name  string `gorm:"varchar(300)" json:"name"`
	Score int8   `gorm:"integer" json:"score"`
}
