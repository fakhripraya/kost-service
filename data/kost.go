package data

import (
	"crypto/rand"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/hashicorp/go-hclog"
	"github.com/jinzhu/gorm"
	"github.com/srinathgs/mysqlstore"
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

		return nil, err
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
func (kost *Kost) GenerateCode(kostType, country, city string) (string, error) {

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

	var finalCode string = kostType +
		"/" + country +
		"-" + city +
		"/" + strconv.Itoa(time.Now().UTC().Year()) + "-" + time.Now().UTC().Month().String()[0:1] +
		"/" + crypted

	return finalCode, nil

}

// UpdateKost is a function to update the given kost model
func (kost *Kost) UpdateKost(targetKost *entities.Kost) error {

	return nil
}

// AddRoom is a function to add kost room based on the given kost id
func (kost *Kost) AddRoom(currentUser *database.MasterUser, kostID uint, targetKostRoom *entities.KostRoom) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		var newKostRoom database.DBKostRoom
		var dbErr error

		newKostRoom.KostID = kostID
		newKostRoom.RoomDesc = targetKostRoom.RoomDesc
		newKostRoom.RoomPrice = targetKostRoom.RoomPrice
		newKostRoom.RoomArea = targetKostRoom.RoomArea
		newKostRoom.IsActive = true
		newKostRoom.Created = time.Now().Local()
		newKostRoom.CreatedBy = currentUser.Username
		newKostRoom.Modified = time.Now().Local()
		newKostRoom.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newKostRoom).Error; dbErr != nil {
			return dbErr
		}

		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var newKostRoomPict database.DBKostRoomPict
			var dbErr2 error

			// loop the room pict slices
			for _, roomPic := range targetKostRoom.RoomPicts {

				newKostRoomPict.RoomID = newKostRoom.ID
				newKostRoomPict.PictDesc = roomPic.PictDesc
				newKostRoomPict.URL = roomPic.URL
				newKostRoomPict.IsActive = true
				newKostRoomPict.Created = time.Now().Local()
				newKostRoomPict.CreatedBy = currentUser.Username
				newKostRoomPict.Modified = time.Now().Local()
				newKostRoomPict.ModifiedBy = currentUser.Username

				if dbErr2 = tx.Create(&newKostRoomPict).Error; dbErr2 != nil {
					return dbErr2
				}

			}

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
func (kost *Kost) AddFacilities(currentUser *database.MasterUser, kostID uint, targetFacilities *entities.MasterFacilities) error {

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		var newKostFacilities database.DBKostFacilities
		var dbErr error

		newKostFacilities.KostID = kostID
		newKostFacilities.FacID = targetFacilities.ID
		newKostFacilities.Created = time.Now().Local()
		newKostFacilities.CreatedBy = currentUser.Username
		newKostFacilities.Modified = time.Now().Local()
		newKostFacilities.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newKostFacilities).Error; dbErr != nil {
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
