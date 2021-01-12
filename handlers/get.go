package handlers

import (
	"net/http"

	"github.com/fakhripraya/authentication-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/entities"
)

// GetKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKost(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetMyKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetMyKost(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetKostList is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKostList(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetMyKostList is a method to fetch the given kost info
func (kostHandler *KostHandler) GetMyKostList(rw http.ResponseWriter, r *http.Request) {

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
	// look for the target user in the db
	var targetKost entities.Kost
	if err := config.DB.Where("username = ?", session.Values["userLoggedin"].(string)).First(&targetKost).Error; err != nil {
		rw.WriteHeader(http.StatusUnauthorized)

		return
	}

	rw.WriteHeader(http.StatusOK)
	return
}
