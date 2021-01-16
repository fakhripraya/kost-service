package database

import "time"

// MasterPeriod is an entity that directly communicate with the MasterPeriod table in the database
type MasterPeriod struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	PeriodDesc string    `gorm:"not null" json:"period_desc"` // annual , monthly , weekly , daily dll
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterPeriodTable set the migrated struct table name
func (masterPeriod *MasterPeriod) MasterPeriodTable() string {
	return "dbMasterPeriod"
}
