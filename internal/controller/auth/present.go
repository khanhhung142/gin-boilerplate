package auth_controller

// Define the input and output of the controller layer here
// Gin already support validation, so we don't need to init a validator

// SignIn
type SignInRequest struct {
	Username string `json:"username" binding:"required,min=8,max=50"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type SignInResponse struct {
	Token string `json:"token"`
	Exp   int64  `json:"exp_time"`
}

// SignUp
type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=8,max=50"`
	Password string `json:"password" binding:"required,min=8,max=50"`
}

type SignUpResponse struct {
	Success bool `json:"success"`
}
