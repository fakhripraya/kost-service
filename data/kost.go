package data

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/hashicorp/go-hclog"
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

// UpdateKost is a function to update the given kost model
func (user *Kost) UpdateKost(targetUser *entities.Kost) error {

	return nil
}
