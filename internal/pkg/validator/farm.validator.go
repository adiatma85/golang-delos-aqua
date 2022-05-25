package validator

// Struct that define the binding of Create Farm Request
type CreateFarmRequest struct {
	Name string `json:"name" form:"name" binding:"required,min=1"`
}

// Struct that define the validator/binding of Update Farm Request
type UpdateFarmRequest struct {
	Name string `json:"name" form:"name" binding:"required,min=1"`
}
