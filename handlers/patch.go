package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"gorm.io/gorm"
)

// AdminApprovalKost is a method to approve the kost info by the admin
func (kostHandler *KostHandler) AdminApprovalKost(rw http.ResponseWriter, r *http.Request) {

	// get the approval via context
	approvalReq := r.Context().Value(KeyApproval{}).(*entities.ApprovalKost)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// proceed to create the new approval with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var targetKost database.DBKost
		var dbErr error

		// look for the existing kost by the given kost id
		if dbErr := config.DB.Where("id = ?", approvalReq.KostID).First(&targetKost).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)

			return dbErr
		}

		// occurs when transaction already been approved by the tenant
		if targetKost.Status != 0 {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Status kost tidak valid untuk di approve")
		}

		// only tenant can approve the book transaction in this method
		if currentUser.RoleID != 0 {
			rw.WriteHeader(http.StatusForbidden)

			return fmt.Errorf("Hanya admin yang bisa approve kost")
		}

		// Status 1 = approved by owner
		// Status 2 = reject
		if approvalReq.FlagApproval == true {
			targetKost.Status = 1
			targetKost.Modified = time.Now().Local()
			targetKost.ModifiedBy = currentUser.Username
		} else {
			targetKost.Status = 2
			targetKost.Modified = time.Now().Local()
			targetKost.ModifiedBy = currentUser.Username
		}

		// update the kost
		dbErr = config.DB.Save(targetKost).Error

		if dbErr != nil {
			return dbErr
		}

		return nil

	})

	// if transaction error
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// TODO: send notif

	rw.WriteHeader(http.StatusOK)
	if approvalReq.FlagApproval == true {
		data.ToJSON(&GenericError{Message: "Sukses Approve kost"}, rw)
	} else {
		data.ToJSON(&GenericError{Message: "Sukses Reject kost"}, rw)
	}

	return

}
