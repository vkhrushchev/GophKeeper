package dto

type APIResponse struct {
	APICode    int    `json:"api_code"`
	APIMessage string `json:"api_message"`
}

type APIRegisterUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type APIRegisterUserResponse struct {
	APIResponse
}

type APILoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type APILoginResponse struct {
	APIResponse
	Token string `json:"token"`
}

type APILogoutRequest struct {
}

type APILogoutResponse struct {
	APIResponse
}
