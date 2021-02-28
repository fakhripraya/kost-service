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
