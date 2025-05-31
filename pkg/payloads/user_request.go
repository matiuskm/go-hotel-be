package payloads

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name" validate:"required"`
	Role string `json:"role" validate:"required,oneof=admin manager receptionist housekeeping"`
}