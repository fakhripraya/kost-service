package handlers

import (
	"net/http"
	"time"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/jinzhu/gorm"
)

// AddKost is a method to add the new given kost info
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

	// proceed to create the new user with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		var newKost database.DBKost
		var dbErr error

		newKost.OwnerID = currentUser.ID
		newKost.TypeID = kostReq.TypeID
		newKost.KostCode, dbErr = kostHandler.kost.GenerateCode("k", kostReq.Country[0:1], kostReq.City[0:1])
		if dbErr != nil {
			return dbErr
		}

		newKost.KostName = kostReq.KostName
		newKost.Country = kostReq.Country
		newKost.City = kostReq.City
		newKost.Address = kostReq.Address
		newKost.UpRate = 0
		newKost.UpRateExpired = time.Now().Local()
		newKost.IsVerified = false
		newKost.IsActive = false
		newKost.StatusAktif = 0
		newKost.Created = time.Now().Local()
		newKost.CreatedBy = currentUser.Username
		newKost.Modified = time.Now().Local()
		newKost.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newKost).Error; dbErr != nil {
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
	return
}
