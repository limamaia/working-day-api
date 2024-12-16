package swagger

// @Summary Get task by ID
// @Description Retrieve task details by ID. Accessible to the task owner or managers.
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "Task ID"
// @Success 200 {object} TaskResponse "Task details retrieved successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized - Missing or invalid token"
// @Failure 403 {object} ForbiddenResponse "Forbidden - You do not have permission to view this task"
// @Failure 404 {object} NotFoundResponse "Task not found"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error - Error decrypting task summary"
// @Router /v1/task/{id} [get]
func getTaskHandlerDocs() {}

// @Summary Get tasks list
// @Description Retrieve a list of tasks. Accessible to managers or the task owner.
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Success 200 {object} GetTasksSuccessResponse "Tasks retrieved successfully"
// @Success 200 {object} NoTasksResponse "No tasks found"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized - Missing or invalid token"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error - Error fetching tasks"
// @Router /v1/task [get]
func getTasksHandlerDocs() {}

// @Summary Create a new task
// @Description Create a new task with a summary and performed date. Accessible to authenticated users.
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param request body CreateTaskRequest true "Task creation payload"
// @Success 201 {object} CreateTaskSuccessResponse "Task created successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized - Missing or invalid token"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error - Error encrypting summary or creating task"
// @Router /v1/task [post]
func createTaskHandlerDocs() {}

// @Summary Update task by ID
// @Description Update the details of an existing task. Accessible only to the task owner.
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "Task ID"
// @Param request body UpdateTaskRequest true "Updated task data"
// @Success 200 {object} UpdateTaskSuccessResponse "Task updated successfully"
// @Failure 400 {object} ErrorResponse "Invalid request payload"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized - Missing or invalid token"
// @Failure 403 {object} ForbiddenResponse "Forbidden - You do not have permission to update this task"
// @Failure 404 {object} NotFoundResponse "Task not found"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error - Error encrypting summary"
// @Router /v1/task/{id} [put]
func updateTaskHandlerDocs() {}

// @Summary Delete task by ID
// @Description Deletes an existing task. Accessible only to the task owner or managers.
// @Tags Task
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token" default(Bearer <token>)
// @Param id path string true "Task ID"
// @Success 200 {object} DeleteTaskSuccessResponse "Task deleted successfully"
// @Failure 401 {object} UnauthorizedResponse "Unauthorized - Missing or invalid token"
// @Failure 403 {object} ForbiddenResponse "Forbidden - You do not have permission to delete this task"
// @Failure 404 {object} NoTasksResponse "Task not found"
// @Failure 500 {object} InternalServerErrorResponse "Internal server error - Error deleting task"
// @Router /v1/task/{id} [delete]
func deleteTaskHandlerDocs() {}

// TaskResponse represents the structure of a successful task response.
// @Description Task details returned from the API.
type TaskResponse struct {
	ID            uint   `json:"id" example:"1"`
	Task          string `json:"task" example:"Complete project documentation"`
	Summary       string `json:"summary" example:"Finalize all pending tasks"`
	PerformedDate string `json:"performed_date" example:"2024-12-15T19:14:00Z"`
	UserID        uint   `json:"user_id" example:"2"`
	CreatedAt     string `json:"created_at" example:"2024-12-15T20:00:00Z"`
	UpdatedAt     string `json:"updated_at,omitempty" example:"2024-12-15T21:00:00Z"`
}

// CreateTaskRequest defines the input structure for creating a task.
// @Description Structure for creating a new task related to a working day.
type CreateTaskRequest struct {
	Task          string `json:"task" binding:"required" example:"Develop new feature for project X"`
	Summary       string `json:"summary" binding:"required" example:"Implement the user authentication module with JWT and ensure it is tested"`
	PerformedDate string `json:"performed_date" binding:"required" example:"2024-12-18T09:00:00Z"`
}

// UnauthorizedResponse represents the error response for unauthorized access (401).
// @Description Missing or invalid token.
type UnauthorizedResponse struct {
	Error string `json:"error" example:"Unauthorized - Missing or invalid token"`
}

// ForbiddenResponse represents the error response for forbidden access (403).
// @Description The user does not have permission to access the requested resource.
type ForbiddenResponse struct {
	Error string `json:"error" example:"Forbidden - You do not have permission to view this task"`
}

// NotFoundResponse represents the error response for resource not found (404).
// @Description The requested resource was not found.
type NotFoundResponse struct {
	Error string `json:"error" example:"Task not found"`
}

// InternalServerErrorResponse represents the error response for internal server errors (500).
// @Description An internal error occurred during task processing.
type InternalServerErrorResponse struct {
	Error string `json:"error" example:"Internal server error - Error decrypting task summary"`
}

// GetTasksSuccessResponse represents the response structure when tasks are successfully retrieved.
// @Description Success response for retrieving tasks.
type GetTasksSuccessResponse struct {
	Message string         `json:"message" example:"Tasks retrieved successfully"`
	Tasks   []TaskResponse `json:"tasks"`
}

// NoTasksResponse represents the response when no tasks are found.
// @Description Response when no tasks are found.
type NoTasksResponse struct {
	Message string         `json:"message" example:"No tasks found"`
	Tasks   []TaskResponse `json:"tasks"`
}

// CreateTaskSuccessResponse represents the response for a successful task creation.
// @Description Successful creation of a task.
type CreateTaskSuccessResponse struct {
	Message string `json:"message" example:"Task created successfully"`
	Status  int    `json:"status" example:"201"`
}

// UpdateTaskRequest defines the input structure for updating a task.
// @Description Structure for updating an existing task.
type UpdateTaskRequest struct {
	Task          string `json:"task" binding:"required" example:"Update project documentation"`
	Summary       string `json:"summary" binding:"required" example:"Add detailed explanations to the API documentation"`
	PerformedDate string `json:"performed_date" binding:"required" example:"2024-12-18T15:30:00Z"`
}

// UpdateTaskSuccessResponse represents the success response after updating a task.
// @Description Success message when a task is updated successfully.
type UpdateTaskSuccessResponse struct {
	Message string `json:"message" example:"Task updated successfully"`
}

// DeleteTaskSuccessResponse represents the success response after deleting a task.
// @Description Success message when a task is deleted successfully.
type DeleteTaskSuccessResponse struct {
	Message string `json:"message" example:"Task deleted successfully"`
}
