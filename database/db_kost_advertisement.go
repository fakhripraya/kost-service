package database

import "time"

// DBKostAds will migrate a kost ads table with the given specification into the database
type DBKostAds struct {
	ID                     uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	OwnerID                uint      `gorm:"not null" json:"owner_id"`
	AdsTypeID              uint      `gorm:"not null" json:"ads_type_id"`
	KostTypeID             uint      `gorm:"not null" json:"kost_type_id"`
	GenderTypeID           uint      `gorm:"not null" json:"gender_type_id"`
	PetAllowedTypeID       uint      `gorm:"not null" json:"pet_allowed_type_id"`
	Status                 uint      `gorm:"not null" json:"status"`
	AdsCode                string    `gorm:"not null" json:"ads_code"`
	AdsDesc                string    `gorm:"not null" json:"ads_desc"`
	AdsPostScheduleRequest string    `gorm:"not null" json:"ads_post_schedule_request"`
	IsActive               bool      `gorm:"not null;default:true" json:"is_active"`
	Created                time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy              string    `json:"created_by"`
	Modified               time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy             string    `json:"modified_by"`
}

// DBKostAdsFiles will migrate a kost ads files table with the given specification into the database
type DBKostAdsFiles struct {
	ID           uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	AdsID        uint      `gorm:"not null" json:"kost_id"`
	BASE64STRING string    `gorm:"not null" json:"base64_string"`
	IsActive     bool      `gorm:"not null;default:true" json:"is_active"`
	Created      time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy    string    `json:"created_by"`
	Modified     time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy   string    `json:"modified_by"`
}

// KostAdvertisementTable set the migrated struct table name
func (dbKostAds *DBKostAds) KostAdsTable() string {
	return "dbKostAds"
}

// KostAdsFilesTable set the migrated struct table name
func (dbKostAdsFiles *DBKostAdsFiles) KostAdsFilesTable() string {
	return "dbKostAdsFiles"
}
