package model

type Unit struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"not null" json:"name"`                         // e.g., "Minutes"
	Symbol      string      `gorm:"not null" json:"symbol"`                       // e.g., "min"
	Measurement Measurement `gorm:"type:varchar(20);not null" json:"measurement"` // e.g., "time"
}
