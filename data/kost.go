package data

import (
	"crypto/rand"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/hashicorp/go-hclog"
	"github.com/srinathgs/mysqlstore"
	"gorm.io/gorm"
)

// Claims determine the current user token holder
type Claims struct {
	Username string
	jwt.StandardClaims
}

// Kost defines a struct for kost flow
type Kost struct {
	logger hclog.Logger
}

// NewKost is a function to create new Kost struct
func NewKost(newLogger hclog.Logger) *Kost {
	return &Kost{newLogger}
}

// GetCurrentUser will get the current user login info
func (kost *Kost) GetCurrentUser(rw http.ResponseWriter, r *http.Request, store *mysqlstore.MySQLStore) (*database.MasterUser, error) {

	// Get a session (existing/new)
	session, err := store.Get(r, "session-name")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)

		return nil, err
	}

	// check the logged in user from the session
	// if user available, get the user info from the session
	if session.Values["userLoggedin"] == nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return nil, fmt.Errorf("Error 401")
	}

	// work with database
	// look for the current user logged in in the db
	var currentUser database.MasterUser
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&currentUser).Error; err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return nil, err
	}

	return &currentUser, nil

}

// GetUOMDesc is a function to get the uom desc by the given uom id
func (kost *Kost) GetUOMDesc(UomID uint) (string, error) {

	var uomDescription string
	if err := config.DB.Raw("SELECT uom_desc FROM master_uoms WHERE id = ?", UomID).Scan(&uomDescription).Error; err != nil {

		return "", err
	}

	return uomDescription, nil
}

// GetLowestPrice is a function to get the lowest price by the given list of prices
func (kost *Kost) GetLowestPrice(KostID uint) (*entities.KostRoomPrice, error) {

	var lowestPrice = &entities.KostRoomPrice{}
	var roomPrices []entities.KostRoomPrice
	if err := config.DB.Raw("SELECT room_price, room_price_uom FROM db_kost_rooms WHERE kost_id = ?", KostID).Scan(&roomPrices).Error; err != nil {

		return nil, err
	}

	for _, price := range roomPrices {

		if lowestPrice.RoomPrice == 0 {
			lowestPrice.RoomPrice = price.RoomPrice
			lowestPrice.RoomPriceUom = price.RoomPriceUom
			lowestPrice.RoomPriceUomDesc, _ = kost.GetUOMDesc(price.RoomPriceUom)
		} else {
			if lowestPrice.RoomPrice > price.RoomPrice {
				lowestPrice.RoomPrice = price.RoomPrice
				lowestPrice.RoomPriceUom = price.RoomPriceUom
				lowestPrice.RoomPriceUomDesc, _ = kost.GetUOMDesc(price.RoomPriceUom)
			}
		}

	}

	return lowestPrice, nil
}

// GetReverseGeocoderResult will get the result of reverse geocoder calculation based on the given latitude and longitude
func (kost *Kost) GetReverseGeocoderResult(latitude string, longitude string) (*entities.Geolocation, error) {

	baseURL, _ := url.Parse("http://api.positionstack.com")

	baseURL.Path += "v1/reverse"

	params := url.Values{}

	// Access Key
	params.Add("access_key", os.Getenv("GEOCODER_API_KEY"))

	// Query = latitude,longitude
	params.Add("query", latitude+","+longitude)

	// trigger the reverse geocoder request to fetch the addresses data
	baseURL.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", baseURL.String(), nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {

		return nil, err
	}

	defer res.Body.Close()

	// create the geo location instance
	geoLocation := &entities.Geolocation{}
	FromJSON(geoLocation, res.Body)

	return geoLocation, nil

}

// CalculateDistanceBetween will calculate the distance between two given point
func (kost *Kost) CalculateDistanceBetween(lat1 float64, lng1 float64, lat2 float64, lng2 float64, unit ...string) float64 {

	const PI float64 = 3.141592653589793

	radlat1 := float64(PI * lat1 / 180)
	radlat2 := float64(PI * lat2 / 180)

	theta := float64(lng1 - lng2)
	radtheta := float64(PI * theta / 180)

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)

	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / PI
	dist = dist * 60 * 1.1515

	// Unit is a unit of measure for length
	// Kilometer(K)
	// Miles(M)
	// Nautican miles(N)

	if len(unit) > 0 {
		if unit[0] == "K" {
			dist = dist * 1.609344
		} else if unit[0] == "N" {
			dist = dist * 0.8684
		}
	}

	return dist

}

// GenerateCode will generate the new given type code
func (kost *Kost) GenerateCode(codeType, country, city string) (string, error) {

	// generate 8 random crypted number
	var max int = 8
	var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

	b := make([]byte, max)
	n, err := io.ReadAtLeast(rand.Reader, b, max)
	if n != max {
		return "", err
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}

	// returns the crypted random 8 number
	var crypted string = string(b)

	var finalCode string = codeType +
		"/" + country +
		"-" + city +
		"/" + strconv.Itoa(time.Now().UTC().Year()) + "-" + time.Now().UTC().Month().String()[0:1] +
		"/" + crypted

	return finalCode, nil

}

// AddRoom is a function to add kost room based on the given kost id
func (kost *Kost) AddRoom(currentUser *database.MasterUser, kostID uint, targetKostRoom *entities.KostRoom) error {

	// add the kostReq room into the database with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newKostRoom database.DBKostRoom
		var targetPriceUOM database.MasterUOM
		var targetAreaUOM database.MasterUOM
		var dbErr error

		newKostRoom.KostID = kostID
		newKostRoom.RoomDesc = targetKostRoom.RoomDesc
		newKostRoom.RoomPrice = targetKostRoom.RoomPrice

		// look for the requested uom from the database
		if dbErr = config.DB.Where("id = ?", targetKostRoom.RoomPriceUOM).First(&targetPriceUOM).Error; dbErr != nil {
			return dbErr
		}

		// check the uom type, if not currency return error
		if targetPriceUOM.UOMType != "currency" {
			return fmt.Errorf("Invalid UOM Type")
		}

		newKostRoom.RoomPriceUOM = targetKostRoom.RoomPriceUOM
		newKostRoom.RoomLength = targetKostRoom.RoomLength
		newKostRoom.RoomWidth = targetKostRoom.RoomWidth
		newKostRoom.RoomArea = targetKostRoom.RoomLength * targetKostRoom.RoomWidth

		// look for the requested uom from the database
		if dbErr = config.DB.Where("id = ?", targetKostRoom.RoomAreaUOM).First(&targetAreaUOM).Error; dbErr != nil {
			return dbErr
		}

		// check the uom type, if not length return error
		if targetAreaUOM.UOMType != "length" {
			return fmt.Errorf("Invalid UOM Type")
		}

		newKostRoom.RoomAreaUOM = targetKostRoom.RoomAreaUOM
		newKostRoom.MaxPerson = targetKostRoom.MaxPerson
		newKostRoom.FloorLevel = targetKostRoom.FloorLevel
		newKostRoom.AllowedGender = targetKostRoom.AllowedGender
		newKostRoom.IsActive = true
		newKostRoom.Created = time.Now().Local()
		newKostRoom.CreatedBy = currentUser.Username
		newKostRoom.Modified = time.Now().Local()
		newKostRoom.ModifiedBy = currentUser.Username

		// insert the new room to the database
		if dbErr = tx.Create(&newKostRoom).Error; dbErr != nil {
			return dbErr
		}

		// add the room details into the database with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var roomDetails = targetKostRoom.RoomDetails

			// add the room id to the slices
			for i := range roomDetails {
				(&roomDetails[i]).RoomID = newKostRoom.ID
				(&roomDetails[i]).IsActive = true
				(&roomDetails[i]).Created = time.Now().Local()
				(&roomDetails[i]).CreatedBy = currentUser.Username
				(&roomDetails[i]).Modified = time.Now().Local()
				(&roomDetails[i]).ModifiedBy = currentUser.Username
			}

			// insert the new room details to database
			if dbErr2 = tx2.Create(&roomDetails).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		if dbErr != nil {
			return dbErr
		}

		// add the room picts into the database with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var roomPicts = targetKostRoom.RoomPicts

			// add the room id to the slices
			for i := range roomPicts {
				(&roomPicts[i]).RoomID = newKostRoom.ID
				(&roomPicts[i]).IsActive = true
				(&roomPicts[i]).Created = time.Now().Local()
				(&roomPicts[i]).CreatedBy = currentUser.Username
				(&roomPicts[i]).Modified = time.Now().Local()
				(&roomPicts[i]).ModifiedBy = currentUser.Username
			}

			// insert the new room picts to database
			if dbErr2 = tx2.Create(&roomPicts).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error
		if dbErr != nil {

			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil
}

// AddFacilities is a function to add kost facilities based on the given kost id
func (kost *Kost) AddFacilities(currentUser *database.MasterUser, kostID uint, targetFacilities []database.DBKostFacilities) error {

	// add the room facilities into the database with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var dbErr error
		var facilities = targetFacilities

		// add the kost id to the slices
		for i := range facilities {
			(&facilities[i]).KostID = kostID
			(&facilities[i]).IsActive = true
			(&facilities[i]).Created = time.Now().Local()
			(&facilities[i]).CreatedBy = currentUser.Username
			(&facilities[i]).Modified = time.Now().Local()
			(&facilities[i]).ModifiedBy = currentUser.Username
		}

		// insert the facilities to the database
		if dbErr = tx.Create(targetFacilities).Error; dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {

		return err
	}

	return nil
}

// GetKostFacilities is a function to get kost facilities by the kost
func (kost *Kost) GetKostFacilities(KostID uint, RoomID string) ([]entities.KostFacilities, []entities.KostRoomFacilities, error) {

	var facTableName string
	var facTableKey string

	var id uint
	var model *gorm.DB
	var kostFacilities []entities.KostFacilities
	var kostRoomFacilities []entities.KostRoomFacilities

	// filter whether id empty or not
	if RoomID != "" {

		facTableKey = "room_id"
		facTableName = "db_kost_room_facilities"
		model = config.DB.Model(&database.DBKostRoomFacilities{})

		roomID, err := strconv.ParseUint(RoomID, 10, 32)
		if err != nil {

			return nil, nil, err
		}

		id = uint(roomID)

	} else {

		id = KostID
		facTableKey = "kost_id"
		facTableName = "db_kost_facilities"
		model = config.DB.Model(&database.DBKostFacilities{})

	}

	// if id is not empty
	// query will execute a select sql statement towards db_kost_room_facilities
	// if not it will go towards db_kost_facilities instead
	finalQuery := model.
		Select(facTableName+".id,"+facTableName+".fac_id,"+facTableName+"."+facTableKey+",master_facilities.fac_category as fac_category, master_facilities.fac_name as fac_desc").
		Joins("inner join master_facilities on master_facilities.id = "+facTableName+".fac_id").
		Where(facTableName+"."+facTableKey+" = ?", id)

	if RoomID != "" {
		finalQuery = finalQuery.Scan(&kostRoomFacilities)
	} else {
		finalQuery = finalQuery.Scan(&kostFacilities)
	}

	if err := finalQuery.Error; err != nil {

		return nil, nil, err
	}

	return kostFacilities, kostRoomFacilities, nil
}

// GetKostByOwner is a function to get kost by owner id
func (kost *Kost) GetKostByOwner(ownerID uint) (*database.DBKost, error) {

	// look for the current kost list in the db
	myKost := &database.DBKost{}
	if err := config.DB.Where("owner_id = ?", ownerID).First(myKost).Error; err != nil {

		return nil, err
	}

	return myKost, nil
}

// GetKostListByOwner is a function to get kost list by owner id
func (kost *Kost) GetKostListByOwner(ownerID uint) ([]database.DBKost, error) {

	// look for the current kost list in the db
	var kostList []database.DBKost
	if err := config.DB.Where("owner_id = ?", ownerID).Find(&kostList).Error; err != nil {

		return nil, err
	}

	return kostList, nil
}

// GetKostList is a function to get kost list
func (kost *Kost) GetKostList(page int) ([]database.DBKost, error) {

	// look for the current kost list in the db
	// 10 is the default limit
	var kostList []database.DBKost
	if err := config.DB.Limit((10 * page)).Find(&kostList).Error; err != nil {

		return nil, err
	}

	return kostList, nil
}

// GetNearbyKostList is a function to get nearby kost list
func (kost *Kost) GetNearbyKostList(latitude, longitude string, page int) ([]entities.KostRanges, error) {

	// Get the geolocation by reversing it
	geoLocation, err := kost.GetReverseGeocoderResult(latitude, longitude)
	if err != nil {

		return nil, err
	}

	limit := 10 * page

	// look for the current kost list in the db
	var nearbyKostList []database.DBKost
	if err := config.DB.Limit(limit).Where("is_active = ? AND city >= ?", true, geoLocation.GeoData[0].County).Find(&nearbyKostList).Error; err != nil {

		return nil, err
	}

	var listRanges []entities.KostRanges
	var tempRanges []entities.KostRanges

	for _, nearby := range nearbyKostList {

		lat1, _ := strconv.ParseFloat(latitude, 64)
		lng1, _ := strconv.ParseFloat(longitude, 64)
		lat2, _ := strconv.ParseFloat(nearby.Latitude, 64)
		lng2, _ := strconv.ParseFloat(nearby.Longitude, 64)

		distance := kost.CalculateDistanceBetween(lat1, lng1, lat2, lng2, "K")
		tempRanges = append(tempRanges, entities.KostRanges{
			KostID:   nearby.ID,
			Distance: distance,
		})

	}

	for i := 1; i < limit; i++ {

		if len(tempRanges) > 0 {

			var smallest float64 = -1
			for _, obj := range tempRanges {

				if smallest == -1 {
					smallest = obj.Distance
				} else {
					if smallest > obj.Distance {
						smallest = obj.Distance
					}
				}
			}

			var keepRanges []entities.KostRanges

			for _, num := range tempRanges {

				if num.Distance == smallest {
					listRanges = append(listRanges, num)
				} else {
					keepRanges = append(keepRanges, num)
				}
			}

			tempRanges = keepRanges

		}

	}

	return listRanges, nil
}

// GetKostRoom is a function to get kost room based on the given room id
func (kost *Kost) GetKostRoom(roomID uint) (*database.DBKostRoom, error) {

	kostRoom := &database.DBKostRoom{}
	if err := config.DB.Where("room_id = ?", roomID).Find(&kostRoom).Error; err != nil {

		return nil, err
	}

	return kostRoom, nil
}

// GetKostRoomDetails is a function to get kost room details based on the given room id
func (kost *Kost) GetKostRoomDetails(roomID uint) ([]database.DBKostRoomDetail, error) {

	var kostRoomDetails []database.DBKostRoomDetail
	if err := config.DB.Where("room_id = ?", roomID).Find(&kostRoomDetails).Error; err != nil {

		return nil, err
	}

	return kostRoomDetails, nil
}

// GetKostRoomPicts is a function to get kost room picts based on the given room id
func (kost *Kost) GetKostRoomPicts(roomID uint) ([]database.DBKostRoomPict, error) {

	var kostRoomPicts []database.DBKostRoomPict
	if err := config.DB.Where("room_id = ?", roomID).Find(&kostRoomPicts).Error; err != nil {

		return nil, err
	}

	return kostRoomPicts, nil
}

// GetKostRoomBookedList is a function to get kost booked room list based on the given room id
func (kost *Kost) GetKostRoomBookedList(roomID uint) ([]database.DBTransactionRoomBook, error) {

	var kostRoomBookedList []database.DBTransactionRoomBook
	if err := config.DB.Where("room_id = ?", roomID).Find(&kostRoomBookedList).Error; err != nil {

		return nil, err
	}

	return kostRoomBookedList, nil
}
