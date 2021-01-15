package database

import "time"

// MasterStatus is an entity that directly communicate with the MasterStatus table in the database
type MasterStatus struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	Category   uint      `gorm:"not null" json:"category"`
	StatusDesc string    `gorm:"not null" json:"status_desc"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterStatusTable set the migrated struct table name
func (masterStatus *MasterStatus) MasterStatusTable() string {
	return "dbMasterStatusKost"
}
