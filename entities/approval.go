package entities

// AdminApprovalKost is an entity to communicate with the AdminApprovalKost client side
type AdminApprovalKost struct {
	KostID       uint `json:"kost_id"`
	FlagApproval bool `json:"flag_approval"`
}
