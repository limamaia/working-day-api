package swagger

// @Summary Create a new user
// @Description Creates a new user with the provided details.
// @Tags User
// @Accept json
// @Produce json
// @Param request body domain.CreateUserRequest true "User creation payload"
// @Success 201 {object} CreateUserSuccessResponse "User created successfully"
// @Failure 400 {object} EmailErrorResponse "Invalid input data"
// @Failure 500 {object} EmailErrorResponse "Internal server error or email already exists"
// @Router /v1/user [post]
func createUserHandlerDocs() {}

// @Summary Get user by ID
// @Description Retrieve user details by ID. Only accessible to managers or the user themselves.
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "User ID"
// @Success 200 {object} UserResponse "User details retrieved successfully"
// @Failure 401 {object} ErrorResponse "Unauthorized - Missing or invalid token"
// @Failure 403 {object} UserErrorResponse "Forbidden - You do not have permission to view this user"
// @Failure 404 {object} UserNotFoundResponse "User not found"
// @Router /v1/user/{id} [get]
func getUserHandlerDocs() {}

// CreateUserRequest represents the payload for creating a user.
// @Description Structure for creating a new user.
type CreateUserRequest struct {
	Name     string `json:"name" binding:"required,min=3,max=255" example:"John Doe"`
	Email    string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Password string `json:"password" binding:"required,min=8" example:"P@ssw0rd"`
	RoleID   uint   `json:"role_id" binding:"required" example:"2"`
}

// CreateUserSuccessResponse represents the response for successful user creation.
// @Description Response for successfully creating a user.
type CreateUserSuccessResponse struct {
	Message string `json:"message" example:"User created successfully"`
	Status  int    `json:"status" example:"201"`
}

// ErrorResponse represents a generic error response.
// @Description Generic error message format.
type EmailErrorResponse struct {
	Error   string `json:"error" example:"A user with this email already exists"`
	Details string `json:"details,omitempty" example:"Invalid data format"`
}

// UserResponse represents the response for a user.
// @Description Response structure for retrieving a user.
type UserResponse struct {
	ID        uint   `json:"id" example:"1"`
	Name      string `json:"name" example:"John Doe"`
	Email     string `json:"email" example:"john.doe@example.com"`
	RoleID    uint   `json:"role_id" example:"2"`
	CreatedAt string `json:"created_at" example:"2024-12-15T19:14:00Z"`
	UpdatedAt string `json:"updated_at,omitempty" example:"2024-12-16T10:00:00Z"`
}

// ErrorResponse represents a generic error response.
// @Description Generic error message format.
type UserErrorResponse struct {
	Error string `json:"error" example:"you do not have permission to view this user"`
}

// ErrorResponse represents a generic error response.
// @Description Generic error message format.
type UserNotFoundResponse struct {
	Error string `json:"error" example:"User not found"`
}
