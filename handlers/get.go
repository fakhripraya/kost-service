package handlers

import (
	"fmt"
	"net/http"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
)

// GetKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKost(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetKostList is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKostList(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetMyKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetMyKost(rw http.ResponseWriter, r *http.Request) {

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// look for the current kost in the db
	var myKost database.DBKost
	if err := config.DB.Where("owner_id = ?", currentUser.ID).First(&myKost).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err = data.ToJSON(myKost, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

// GetMyKostList is a method to fetch the list of the given kost info
func (kostHandler *KostHandler) GetMyKostList(rw http.ResponseWriter, r *http.Request) {

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// look for the current kost list in the db
	var kostList database.DBKost
	if err := config.DB.Where("owner_id = ?", currentUser.ID).Find(&kostList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err = data.ToJSON(kostList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}

// GetEventList is a method to fetch the list of application event
func (kostHandler *KostHandler) GetEventList(rw http.ResponseWriter, r *http.Request) {

	// look for the current kost list in the db
	var eventList []database.MasterEvent
	if err := config.DB.Where("is_active = ?", true).Find(&eventList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(eventList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetNearYouList is a method to fetch the list of nearby kost
func (kostHandler *KostHandler) GetNearYouList(rw http.ResponseWriter, r *http.Request) {

	// get the current logged in user additional info via context
	userReq := r.Context().Value(KeyUser{}).(*entities.User)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	kostHandler.logger.Info(currentUser.DisplayName)

	// look for the current kost list in the db
	var nearbyKostList []database.DBKost
	if err := config.DB.Where("is_active = ?", true).Find(&nearbyKostList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	fmt.Printf("%v", nearbyKostList)

	// parse the given instance to the response writer
	err = data.ToJSON(nearbyKostList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}
