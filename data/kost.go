package data

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/hashicorp/go-hclog"
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

// UpdateKost is a function to update the given kost model
func (kost *Kost) UpdateKost(targetUser *entities.Kost) error {

	return nil
}
