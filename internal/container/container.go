package container

import (
	"working-day-api/config"
	"working-day-api/database"
	"working-day-api/internal/controllers"
	"working-day-api/internal/helpers"
	"working-day-api/internal/repositories"
	"working-day-api/internal/services"
	"working-day-api/messaging"
)

type Container struct {
	Config           *config.AppVars
	HealthController *controllers.HealthController
	AuthController   *controllers.AuthController
	TaskController   *controllers.TaskController
	UserController   *controllers.UserController
}

func NewContainer(config *config.AppVars) *Container {
	database.Connection(config)
	messaging.Connection(config)

	userRepo := &repositories.UserRepositoryImpl{}
	taskRepo := &repositories.TaskRepositoryImpl{}

	messenger := &messaging.RabbitMessenger{}
	jwtService := helpers.NewJWTService()
	hasher := helpers.NewPasswordHasher()

	healthService := &services.HealthService{}
	userService := &services.UserService{UserRepo: userRepo, Hasher: hasher}

	loginService := &services.LoginService{
		UserRepo: userRepo,
		JWT:      jwtService,
		Hasher:   hasher,
	}
	taskService := &services.TaskService{
		TaskRepo:  taskRepo,
		Messenger: messenger,
		Decrypt:   helpers.Decrypt,
		Encrypt:   helpers.Encrypt,
	}

	healthController := controllers.NewHealthController(healthService)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(loginService)
	taskController := controllers.NewTaskController(taskService)

	return &Container{
		Config:           config,
		HealthController: healthController,
		UserController:   userController,
		AuthController:   authController,
		TaskController:   taskController,
	}
}
