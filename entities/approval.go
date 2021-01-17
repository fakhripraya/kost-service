package entities

// ApprovalKost is an entity to communicate with the ApprovalKost client side
type ApprovalKost struct {
	KostID       uint `json:"kost_id"`
	FlagApproval bool `json:"flag_approval"`
}
