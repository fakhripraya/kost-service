package entities

import "time"

// Kost is an entity to communicate with Kost table in database
type Kost struct {
	ID            uint               `json:"id"`
	OwnerID       uint               `json:"owner_id"`
	TypeID        uint               `json:"type_id"`
	KostCode      string             `json:"kost_code"`
	KostName      string             `json:"kost_name"`
	Country       string             `json:"country"`
	City          string             `json:"city"`
	Address       string             `json:"address"`
	UpRate        uint64             `json:"up_rate"`
	UpRateExpired time.Time          `json:"up_rate_expired"`
	Rate          uint64             `json:"rate"`
	Rooms         []KostRoom         `json:"rooms"`
	Facilities    []MasterFacilities `json:"facilities"`
	IsVerified    bool               `json:"is_verified"`
	IsActive      bool               `json:"is_active"`
	StatusAktif   uint64             `json:"status_aktif"`
	ThumbnailURL  string             `json:"thumbnail_url"`
	Created       time.Time          `json:"created"`
	CreatedBy     string             `json:"created_by"`
	Modified      time.Time          `json:"modified"`
	ModifiedBy    string             `json:"modified_by"`
}

// KostRoom is an entity to communicate with KostRoom table in database
type KostRoom struct {
	ID         uint           `json:"id"`
	KostID     uint           `json:"kost_id"`
	RoomDesc   string         `json:"room_desc"`
	RoomPrice  uint64         `json:"room_price"`
	RoomArea   uint64         `json:"room_area"`
	RoomPicts  []KostRoomPict `json:"room_picts"`
	IsActive   bool           `json:"is_active"`
	Created    time.Time      `json:"created"`
	CreatedBy  string         `json:"created_by"`
	Modified   time.Time      `json:"modified"`
	ModifiedBy string         `json:"modified_by"`
}

// KostRoomPict is an entity to communicate with KostRoomPict table in database
type KostRoomPict struct {
	ID         uint      `json:"id"`
	RoomID     uint      `json:"room_id"`
	PictDesc   string    `json:"pict_desc"`
	URL        string    `json:"url"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// KostFacilities is an entity to communicate with KostFacilities table in database
type KostFacilities struct {
	FacID      uint      `json:"fac_id"`
	KostID     uint      `json:"kost_id"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterKostType is an entity to communicate with MasterKostType table in database
type MasterKostType struct {
	ID         uint      `json:"id"`
	TypeDesc   string    `json:"type_desc"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterFacilities is an entity to communicate with MasterFacilities table in database
type MasterFacilities struct {
	ID         uint      `json:"id"`
	FacName    string    `json:"fac_name"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterStatusKost is an entity to communicate with MasterStatusKost table in database
type MasterStatusKost struct {
	ID         uint      `json:"id"`
	StatusDesc string    `json:"status_desc"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
