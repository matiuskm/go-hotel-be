package payloads

type UserResponse struct {
	ID 				uint		`json:"id"`
	Username 		string		`json:"username"`
	FullName 		string		`json:"full_name"`
	Role 			string		`json:"role"`
}