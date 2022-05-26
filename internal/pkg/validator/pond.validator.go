package validator

// Struct that define the validator/binding of Create Pond Request
type CreatePondRequest struct {
	Name   string `json:"name" form:"name" binding:"required,min=1"`
	PondId uint   `json:"pond_id" form:"pond_id" binding:"required"`
}

// Struct that define the validator/binding of Update Farm Request
type UpdatePondRequest struct {
	ID   uint   `json:"id" form:"id"`
	Name string `json:"name" form:"name"`
}
