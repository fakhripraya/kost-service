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

	// Get a session (existing/new)
	session, err := kostHandler.store.Get(r, "session-name")
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// check the logged in user from the session
	// if user available, get the user info from the session
	if session.Values["userLoggedin"] == nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	// work with database
	// look for the current user logged in in the db
	var currentUser database.MasterUser
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&currentUser).Error; err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	// proceed to create the new user with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// do some database operations in the transaction (use 'tx' from this point, not 'db')
		var newKost database.DBKost
		var dbErr error

		newKost.OwnerID = currentUser.ID
		newKost.TypeID = kostReq.TypeID
		newKost.KostCode = kostReq.KostCode
		newKost.KostName = kostReq.KostName
		newKost.Address = kostReq.Address
		newKost.UpRate = 0
		newKost.UpRateExpired = time.Now().Local()
		newKost.IsVerified = false
		newKost.IsActive = true
		newKost.Created = time.Now().Local()
		newKost.CreatedBy = "SYSTEM"
		newKost.Modified = time.Now().Local()
		newKost.ModifiedBy = "SYSTEM"

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
