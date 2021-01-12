package database

import "time"

// DBKost is an entity that directly communicate with the Kost table in the database
type DBKost struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	OwnerID    uint      `gorm:"not null" json:"owner_id"`
	TypeID     uint      `gorm:"not null" json:"type_id"`
	KostCode   string    `gorm:"not null" json:"kost_code"`
	KostName   string    `gorm:"not null" json:"kost_name"`
	Address    string    `gorm:"not null" json:"address"`
	Rate       uint64    `json:"rate"`
	IsVerified bool      `gorm:"not null;default:false" json:"is_verified"`
	IsActive   bool      `gorm:"not null;default:true" json:"is_active"`
	Created    time.Time `gorm:"type:datetime" json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `gorm:"type:datetime" json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// DBKostRoom is an entity that directly communicate with the KostRoom table in the database
type DBKostRoom struct {
	ID         uint      `gorm:"primary_key;autoIncrement;not null" json:"id"`
	KostID     string    `gorm:"not null" json:"type_code"`
	RoomDesc   string    `gorm:"not null" json:"room_desc"`
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

// KostFacilitiesTable set the migrated struct table name
func (dbKostFacilities *DBKostFacilities) KostFacilitiesTable() string {
	return "dbKostFacilities"
}
