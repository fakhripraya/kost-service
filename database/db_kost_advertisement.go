package database

import "time"

// DBKostAds will migrate a kost ads table with the given specification into the database
type DBKostAds struct {
	ID                     uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	Status                 uint      `gorm:"not null" json:"status"`
	AdsCode                string    `gorm:"not null" json:"ads_code"`
	AdsType                string    `gorm:"not null" json:"ads_type"`
	AdsKostType            string    `gorm:"not null" json:"ads_kost_type"`
	AdsOwner               string    `gorm:"not null" json:"ads_owner"`
	AdsOwnerIG             string    `gorm:"not null" json:"ads_owner_ig"`
	AdsPhoneNumber         string    `gorm:"not null" json:"ads_phone_number"`
	AdsPICWhatsapp         string    `gorm:"not null" json:"ads_pic_whatsapp"`
	AdsPropertyAddress     string    `gorm:"not null" json:"ads_property_address"`
	AdsPropertyCity        string    `gorm:"not null" json:"ads_property_city"`
	AdsPropertyPrice       string    `gorm:"not null" json:"ads_property_price"`
	AdsDesc                string    `gorm:"not null" json:"ads_desc"`
	AdsGender              string    `gorm:"not null" json:"ads_gender"`
	AdsPetAllowed          string    `gorm:"not null" json:"ads_pet_allowed"`
	AdsPostScheduleRequest string    `gorm:"not null" json:"ads_post_schedule_request"`
	AdsHashtag             string    `gorm:"not null" json:"ads_hashtag"`
	AdsLinkSwipeUp         string    `json:"ads_link_swipe_up"`
	AdsIgBioLink           string    `json:"ads_ig_bio_link"`
	IsActive               bool      `gorm:"not null;default:true" json:"is_active"`
	Created                time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy              string    `json:"created_by"`
	Modified               time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy             string    `json:"modified_by"`
}

// DBKostAdsFiles will migrate a kost ads files table with the given specification into the database
type DBKostAdsFiles struct {
	ID           uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	AdsID        uint      `gorm:"not null" json:"ads_id"`
	AdsFileType  string    `gorm:"not null" json:"ads_file_type"`
	AdsDirPath   string    `gorm:"not null" json:"ads_dir_path"`
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
