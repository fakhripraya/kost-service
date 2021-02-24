package entities

// User is an entity to communicate with the current logged in user additional info on the client side
type User struct {
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}
