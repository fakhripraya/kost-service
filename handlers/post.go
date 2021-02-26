package handlers

import (
	"net/http"
	"time"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"gorm.io/gorm"
)

// AddKost is a method to add the new given kost info to the database
func (kostHandler *KostHandler) AddKost(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// proceed to create the new kost with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newKost database.DBKost
		var dbErr error

		newKost.OwnerID = currentUser.ID
		newKost.TypeID = kostReq.TypeID // kategori kos kosan atau kontrakan atau dll
		newKost.Status = 0              // TODO: create a documented status later // status 0 = baru
		newKost.KostCode, dbErr = kostHandler.kost.GenerateCode("K", kostReq.Country[0:1], kostReq.City[0:1])

		if dbErr != nil {
			return dbErr
		}

		newKost.KostName = kostReq.KostName
		newKost.KostDesc = kostReq.KostDesc
		newKost.Country = kostReq.Country
		newKost.City = kostReq.City
		newKost.Address = kostReq.Address
		newKost.Latitude = kostReq.Latitude
		newKost.Longitude = kostReq.Longitude
		newKost.ThumbnailURL = kostReq.Rooms[0].RoomPicts[0].URL
		newKost.UpRate = 0
		newKost.UpRateExpired = time.Now().Local()
		newKost.IsVerified = false
		newKost.IsActive = true
		newKost.Created = time.Now().Local()
		newKost.CreatedBy = currentUser.Username
		newKost.Modified = time.Now().Local()
		newKost.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newKost).Error; dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost periods with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostPeriod = kostReq.KostPeriods

			// add the kost id to the slices
			for i := range kostPeriod {
				(&kostPeriod[i]).KostID = newKost.ID
				(&kostPeriod[i]).IsActive = true
				(&kostPeriod[i]).Created = time.Now().Local()
				(&kostPeriod[i]).CreatedBy = currentUser.Username
				(&kostPeriod[i]).Modified = time.Now().Local()
				(&kostPeriod[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost periods to database
			if dbErr2 = tx2.Create(&kostPeriod).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost picts with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostPicts = kostReq.KostPicts

			// add the kost id to the slices
			for i := range kostPicts {
				(&kostPicts[i]).KostID = newKost.ID
				(&kostPicts[i]).IsActive = true
				(&kostPicts[i]).Created = time.Now().Local()
				(&kostPicts[i]).CreatedBy = currentUser.Username
				(&kostPicts[i]).Modified = time.Now().Local()
				(&kostPicts[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost picts to database
			if dbErr2 = tx2.Create(&kostPicts).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// loop the room slices
		for _, room := range kostReq.Rooms {

			// add the kostReq room slices into the database
			dbErr = kostHandler.kost.AddRoom(currentUser, newKost.ID, &room)

			// if transaction error, return the error
			if dbErr != nil {
				return dbErr
			}

		}

		// add the kost facilities to the database
		dbErr = kostHandler.kost.AddFacilities(currentUser, newKost.ID, kostReq.Facilities)

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericError{Message: "Sukses menambah kost baru"}, rw)
	return
}
