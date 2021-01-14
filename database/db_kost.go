package database

import "time"

// DBKost is an entity that directly communicate with the Kost table in the database
type DBKost struct {
	ID            uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	OwnerID       uint      `gorm:"not null" json:"owner_id"`
	TypeID        uint      `gorm:"not null" json:"type_id"`
	KostCode      string    `gorm:"not null" json:"kost_code"`
	KostName      string    `gorm:"not null" json:"kost_name"`
	Country       string    `gorm:"not null" json:"country"`
	City          string    `gorm:"not null" json:"city"`
	Address       string    `gorm:"not null" json:"address"`
	UpRate        uint64    `json:"up_rate"`
	UpRateExpired time.Time `json:"up_rate_expired"`
	Rate          uint64    `json:"rate"`
	IsVerified    bool      `gorm:"not null;default:false" json:"is_verified"`
	IsActive      bool      `gorm:"not null;default:true" json:"is_active"`
	StatusAktif   uint64    `gorm:"not null;default:0" json:"status_aktif"`
	ThumbnailURL  string    `json:"thumbnail_url"`
	Created       time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy     string    `json:"created_by"`
	Modified      time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy    string    `json:"modified_by"`
}

// DBKostRoom is an entity that directly communicate with the KostRoom table in the database
type DBKostRoom struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	KostID     uint      `gorm:"not null" json:"kost_id"`
	RoomDesc   string    `gorm:"not null" json:"room_desc"`
	RoomPrice  uint64    `gorm:"not null" json:"room_price"`
	RoomArea   uint64    `gorm:"not null" json:"room_area"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// DBKostRoomPict is an entity that directly communicate with the KostRoomPict table in the database
type DBKostRoomPict struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	RoomID     uint      `gorm:"not null" json:"room_id"`
	PictDesc   string    `gorm:"not null" json:"pict_desc"`
	URL        string    `gorm:"not null" json:"url"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// DBKostFacilities is an entity that directly communicate with the KostFacilities table in the database
type DBKostFacilities struct {
	FacID      uint      `gorm:"not null" json:"fac_id"`
	KostID     uint      `gorm:"not null" json:"kost_id"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// KostTable set the migrated struct table name
func (dbKost *DBKost) KostTable() string {
	return "dbKost"
}

// KostRoomTable set the migrated struct table name
func (dbKostRoom *DBKostRoom) KostRoomTable() string {
	return "dbKostRoom"
}

// KostRoomPictTable set the migrated struct table name
func (dbKostRoomPict *DBKostRoomPict) KostRoomPictTable() string {
	return "dbKostRoomPict"
}

// KostFacilitiesTable set the migrated struct table name
func (dbKostFacilities *DBKostFacilities) KostFacilitiesTable() string {
	return "dbKostFacilities"
}
