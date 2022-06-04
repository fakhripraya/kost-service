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
	KostBenchmark []database.DBKostBenchmark  `json:"kost_benchmark"`
	KostAccess    []database.DBKostAccess     `json:"kost_access"`
	KostAround    []database.DBKostAround     `json:"kost_around"`
	IsVerified    bool                        `json:"is_verified"`
	ThumbnailURL  string                      `json:"thumbnail_url"`
	Distance      float64                     `json:"distance"`
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
	PeriodDesc string    `json:"period_desc"`
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
	ID               uint                        `json:"id"`
	KostID           uint                        `json:"kost_id"`
	RoomDesc         string                      `json:"room_desc"`
	RoomPrice        float64                     `json:"room_price"`
	RoomPriceUOM     uint                        `json:"room_price_uom"`
	RoomPriceUOMDesc string                      `json:"room_price_uom_desc"`
	RoomLength       float64                     `json:"room_length"`
	RoomWidth        float64                     `json:"room_width"`
	RoomArea         float64                     `json:"room_area"`
	RoomAreaUOM      uint                        `json:"room_area_uom"`
	RoomAreaUOMDesc  string                      `json:"room_area_uom_desc"`
	MaxPerson        uint                        `json:"max_person"`
	AllowedGender    string                      `json:"allowed_gender"`
	Comments         string                      `json:"comments"`
	RoomPicts        []database.DBKostRoomPict   `gorm:"-" json:"room_picts"`
	RoomDetails      []database.DBKostRoomDetail `gorm:"-" json:"room_details"`
	IsActive         bool                        `json:"is_active"`
	Created          time.Time                   `json:"created"`
	CreatedBy        string                      `json:"created_by"`
	Modified         time.Time                   `json:"modified"`
	ModifiedBy       string                      `json:"modified_by"`
}

// KostRoomDetail is an entity to communicate with the kost room detail client side
type KostRoomDetail struct {
	ID          uint                 `json:"id"`
	KostID      uint                 `json:"kost_id"`
	RoomID      uint                 `json:"room_id"`
	RoomDesc    string               `json:"room_desc"`
	RoomNumber  string               `json:"room_number"`
	FloorLevel  uint                 `json:"floor_level"`
	Price       float64              `json:"price"`
	Currency    string               `json:"currency"`
	Status      uint                 `json:"status"`
	Booker      *database.MasterUser `json:"booker"`
	PrevPayment time.Time            `json:"prev_payment"`
	NextPayment time.Time            `json:"next_payment"`
	IsActive    bool                 `json:"is_active"`
	Created     time.Time            `json:"created"`
	CreatedBy   string               `json:"created_by"`
	Modified    time.Time            `json:"modified"`
	ModifiedBy  string               `json:"modified_by"`
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
	ID          uint      `json:"id"`
	FacID       uint      `json:"fac_id"`
	KostID      uint      `json:"kost_id"`
	FacCategory uint      `json:"fac_category"`
	FacDesc     string    `json:"fac_desc"`
	IsActive    bool      `json:"is_active"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Modified    time.Time `json:"modified"`
	ModifiedBy  string    `json:"modified_by"`
}

// KostRoomFacilities is an entity to communicate with the kost room client side
type KostRoomFacilities struct {
	ID          uint      `json:"id"`
	FacID       uint      `json:"fac_id"`
	RoomID      uint      `json:"room_id"`
	FacCategory uint      `json:"fac_category"`
	FacDesc     string    `json:"fac_desc"`
	IsActive    bool      `json:"is_active"`
	Created     time.Time `json:"created"`
	CreatedBy   string    `json:"created_by"`
	Modified    time.Time `json:"modified"`
	ModifiedBy  string    `json:"modified_by"`
}

// KostReview is an entity to communicate with the kost review client side
type KostReview struct {
	ID             uint      `json:"id"`
	KostID         uint      `json:"owner_id"`
	UserID         uint      `json:"user_id"`
	DisplayName    string    `json:"display_name"`
	ProfilePicture string    `json:"profile_picture"`
	Cleanliness    float64   `json:"cleanliness"`
	Convenience    float64   `json:"convenience"`
	Security       float64   `json:"security"`
	Facilities     float64   `json:"facilities"`
	Comments       string    `json:"comments"`
	IsActive       bool      `json:"is_active"`
	Created        time.Time `json:"created"`
	CreatedBy      string    `json:"created_by"`
	Modified       time.Time `json:"modified"`
	ModifiedBy     string    `json:"modified_by"`
}

// KostRoomPrice is an entity to communicate with the kost room price client side
type KostRoomPrice struct {
	RoomPrice        float64 `json:"room_price"`
	RoomPriceUom     uint    `json:"room_price_uom"`
	RoomPriceUomDesc string  `json:"room_price_uom_desc"`
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

// KostAds is an entity to communicate with the kost ads client side
type KostAds struct {
	ID                     uint                      `json:"id"`
	Status                 uint                      `json:"status"`
	AdsCode                string                    `json:"ads_code"`
	AdsType                string                    `json:"ads_type"`
	AdsKostType            string                    `json:"ads_kost_type"`
	AdsOwner               string                    `json:"ads_owner"`
	AdsOwnerIG             string                    `json:"ads_owner_ig"`
	AdsPhoneNumber         string                    `json:"ads_phone_number"`
	AdsPICWhatsapp         string                    `json:"ads_pic_whatsapp"`
	AdsPropertyAddress     string                    `json:"ads_property_address"`
	AdsPropertyCity        string                    `json:"ads_property_city"`
	AdsPropertyPrice       string                    `json:"ads_property_price"`
	AdsDesc                string                    `json:"ads_desc"`
	AdsGender              string                    `json:"ads_gender"`
	AdsPetAllowed          string                    `json:"ads_pet_allowed"`
	AdsPostScheduleRequest string                    `json:"ads_post_schedule_request"`
	AdsHashtag             string                    `json:"ads_hashtag"`
	AdsLinkSwipeUp         string                    `json:"ads_link_swipe_up"`
	AdsIgBioLink           string                    `json:"ads_ig_bio_link"`
	AdsFiles               []database.DBKostAdsFiles `gorm:"-" json:"ads_files"`
	IsActive               bool                      `json:"is_active"`
	Created                time.Time                 `json:"created"`
	CreatedBy              string                    `json:"created_by"`
	Modified               time.Time                 `json:"modified"`
	ModifiedBy             string                    `json:"modified_by"`
}

// KostAdsFiles is an entity to communicate with the kost ads file client side
type KostAdsFiles struct {
	ID           uint      `json:"id"`
	AdsID        uint      `json:"ads_id"`
	AdsFileType  string    `json:"ads_file_type"`
	BASE64STRING string    `json:"base64_string"`
	IsActive     bool      `json:"is_active"`
	Created      time.Time `json:"created"`
	CreatedBy    string    `json:"created_by"`
	Modified     time.Time `json:"modified"`
	ModifiedBy   string    `json:"modified_by"`
}
