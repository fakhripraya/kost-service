package entities

import (
	"time"

	"github.com/fakhripraya/kost-service/database"
)

// Kost is an entity to communicate with the kost client side
type Kost struct {
	ID            uint                        `json:"id"`
	OwnerID       uint                        `json:"owner_id"`
	TypeID        uint                        `json:"type_id"`
	Status        uint                        `json:"status"`
	KostCode      string                      `json:"kost_code"`
	KostName      string                      `json:"kost_name"`
	KostDesc      string                      `json:"kost_desc"`
	Country       string                      `json:"country"`
	City          string                      `json:"city"`
	Address       string                      `json:"address"`
	Latitude      string                      `json:"latitude"`
	Longitude     string                      `json:"longitude"`
	UpRate        uint64                      `json:"up_rate"`
	UpRateExpired time.Time                   `json:"up_rate_expired"`
	Rooms         []KostRoom                  `json:"rooms"`
	Facilities    []database.DBKostFacilities `json:"facilities"`
	KostPeriods   []database.DBKostPeriod     `json:"kost_periods"`
	KostPicts     []database.DBKostPict       `json:"kost_picts"`
	IsVerified    bool                        `json:"is_verified"`
	ThumbnailURL  string                      `json:"thumbnail_url"`
	IsActive      bool                        `json:"is_active"`
	Created       time.Time                   `json:"created"`
	CreatedBy     string                      `json:"created_by"`
	Modified      time.Time                   `json:"modified"`
	ModifiedBy    string                      `json:"modified_by"`
}

// KostPeriod is an entity to communicate with the kost period client side
type KostPeriod struct {
	ID         uint      `json:"id"`
	KostID     uint      `json:"kost_id"`
	PeriodID   uint      `json:"period_id"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// KostPict is an entity to communicate with the kost pict client side
type KostPict struct {
	ID         uint      `json:"id"`
	KostID     uint      `json:"kost_id"`
	PictDesc   string    `json:"pict_desc"`
	URL        string    `json:"url"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// KostRoom is an entity to communicate with the kost room client side
type KostRoom struct {
	ID           uint                        `json:"id"`
	KostID       uint                        `json:"kost_id"`
	RoomDesc     string                      `json:"room_desc"`
	RoomPrice    float64                     `json:"room_price"`
	RoomPriceUOM uint                        `json:"room_price_uom"`
	RoomArea     float64                     `json:"room_area"`
	RoomAreaUOM  uint                        `json:"room_area_uom"`
	MaxPerson    uint                        `json:"max_person"`
	FloorLevel   uint                        `json:"floor_level"`
	RoomPicts    []database.DBKostRoomPict   `json:"room_picts"`
	RoomDetails  []database.DBKostRoomDetail `json:"room_details"`
	IsActive     bool                        `json:"is_active"`
	Created      time.Time                   `json:"created"`
	CreatedBy    string                      `json:"created_by"`
	Modified     time.Time                   `json:"modified"`
	ModifiedBy   string                      `json:"modified_by"`
}

// KostRoomDetail is an entity to communicate with the kost room detail client side
type KostRoomDetail struct {
	ID         uint      `json:"id"`
	RoomID     uint      `json:"room_id"`
	RoomNumber string    `json:"room_number"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// KostRoomPict is an entity to communicate with the kost room pict client side
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

// KostFacilities is an entity to communicate with the kost facilities client side
type KostFacilities struct {
	ID         uint      `json:"id"`
	FacID      uint      `json:"fac_id"`
	KostID     uint      `json:"kost_id"`
	FacDesc    string    `json:"fac_desc"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterKostType is an entity to communicate with the master kost type client side
type MasterKostType struct {
	ID         uint      `json:"id"`
	TypeDesc   string    `json:"type_desc"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterFacilities is an entity to communicate with the master facilities client side
type MasterFacilities struct {
	ID         uint      `json:"id"`
	FacName    string    `json:"fac_name"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}

// MasterUOM is an entity to communicate with the master uom client side
type MasterUOM struct {
	ID         uint      `json:"id"`
	UOMType    string    `json:"uom_type"`
	UOMDesc    string    `json:"uom_desc"`
	UOMRate    string    `json:"uom_rate"`
	IsActive   bool      `json:"is_active"`
	Created    time.Time `json:"created"`
	CreatedBy  string    `json:"created_by"`
	Modified   time.Time `json:"modified"`
	ModifiedBy string    `json:"modified_by"`
}
