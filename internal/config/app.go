package config

import (
	"todo-app/internal/delivery/http"
	"todo-app/internal/delivery/http/middleware"
	"todo-app/internal/delivery/http/route"
	"todo-app/internal/repository"
	"todo-app/internal/usecase"

	"github.com/IBM/sarama"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB       *gorm.DB
	App      *fiber.App
	Log      *logrus.Logger
	Validate *validator.Validate
	Config   *viper.Viper
	Producer sarama.SyncProducer
}

func Bootstrap(config *BootstrapConfig) {
	// setup repositories
	userRepository := repository.NewUserRepository(config.Log)
	roleRepository := repository.NewRoleRepository(config.Log)

	// setup use cases
	userUseCase := usecase.NewUserUseCase(config.DB, config.Log, config.Validate, userRepository)
	roleUseCase := usecase.NewRoleUseCase(config.DB, config.Log, config.Validate, roleRepository)

	// setup controller
	userController := http.NewUserController(userUseCase, config.Log)
	roleController := http.NewRoleController(roleUseCase, config.Log)

	// setup middleware
	authMiddleware := middleware.NewAuth(userUseCase)

	routeConfig := route.RouteConfig{
		App:            config.App,
		UserController: userController,
		RoleController: roleController,
		AuthMiddleware: authMiddleware,
	}
	routeConfig.Setup()
}
