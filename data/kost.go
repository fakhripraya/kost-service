package data

import (
	"crypto/rand"
	"fmt"
	"io"
	"net/http"
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
		newKostRoom.RoomArea = targetKostRoom.RoomArea

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

// UpdateKost is a function to update the given kost model
func (kost *Kost) UpdateKost(targetKost *entities.Kost) error {

	return nil
}
