package app

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/controller"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/middleware"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/router"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/aws"
	"github.com/PRC-36/amikompedia-fiber/shared/mail"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	"gorm.io/gorm"
)

type BootstrapConfig struct {
	DB          *gorm.DB
	App         *fiber.App
	Validate    *validator.Validate
	TokenMaker  token.Maker
	EmailSender mail.EmailSender
	AwsS3       aws.AwsS3Action
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	registerRepository := repository.NewRegisterRepository()
	surveyRepository := repository.NewSurveyRepository()
	userRepository := repository.NewUserRepository()
	imageRepository := repository.NewImageRepository()

	// setup usecases
	registerUsecase := usecase.NewRegisterUsecase(config.DB, config.Validate, config.EmailSender, registerRepository, repository.NewOtpRepository())
	surveyUsecase := usecase.NewSurveyUsecase(config.DB, config.Validate, surveyRepository)
	userUsecase := usecase.NewUserUsecase(config.DB, config.Validate, config.AwsS3, userRepository, imageRepository)
	loginUsecase := usecase.NewLoginUsecase(config.DB, config.Validate, config.EmailSender, config.TokenMaker, userRepository)

	// setup controller
	authController := controller.NewAuthController(registerUsecase, loginUsecase)
	surveyController := controller.NewSurveyController(surveyUsecase)
	userController := controller.NewUserController(userUsecase)

	// setup middleware
	authMiddleware := middleware.AuthMiddleware(config.TokenMaker)

	routeConfig := router.RouteConfig{
		App:              config.App,
		AuthMiddleware:   authMiddleware,
		AuthController:   authController,
		SurveyController: surveyController,
		UserController:   userController,
	}

	routeConfig.Setup()
}
