package swagger

// @Summary User login
// @Description Generates a JWT token for the authenticated user
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.Login true "Login credentials"
// @Success 200 {object} LoginSuccessResponse "Token successfully generated"
// @Failure 400 {object} ErrorResponseInvalidJSON "Invalid JSON payload"
// @Failure 401 {object} ErrorResponseInvalidCredentials "Invalid credentials"
// @Router /v1/login [post]
func loginHandlerDocs() {}

// Login defines the credentials for authentication.
// @Description Structure for login credentials.
type Login struct {
	Email    string `json:"email" binding:"required,email" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"supersecret"`
}

// LoginSuccessResponse represents the successful login response.
type LoginSuccessResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}

// ErrorResponseInvalidJSON represents the error response for invalid JSON payload.
// @Description Invalid JSON payload.
type ErrorResponseInvalidJSON struct {
	Error string `json:"error" example:"Invalid JSON payload"`
}

// ErrorResponseInvalidCredentials represents the error response for invalid credentials.
// @Description Invalid credentials provided.
type ErrorResponseInvalidCredentials struct {
	Error string `json:"error" example:"Invalid credentials"`
}

// ErrorResponse represents a generic error response with a variable message.
// @Description Generic error response.
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid credentials"`
}
