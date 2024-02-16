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
	"github.com/PRC-36/amikompedia-fiber/shared/util"
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
	ViperConfig util.Config
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	registerRepository := repository.NewRegisterRepository()
	surveyRepository := repository.NewSurveyRepository()
	userRepository := repository.NewUserRepository()
	imageRepository := repository.NewImageRepository()
	sessionRepository := repository.NewSessionRepository()
	postRepository := repository.NewPostRepository()
	otpRepository := repository.NewOtpRepository()

	// setup usecases
	registerUsecase := usecase.NewRegisterUsecase(config.DB, config.Validate, config.EmailSender, registerRepository, otpRepository)
	surveyUsecase := usecase.NewSurveyUsecase(config.DB, config.Validate, surveyRepository)
	userUsecase := usecase.NewUserUsecase(config.DB, config.Validate, config.AwsS3, config.EmailSender, userRepository, imageRepository, otpRepository)
	loginUsecase := usecase.NewLoginUsecase(config.DB, config.Validate, config.EmailSender, config.TokenMaker, config.ViperConfig, userRepository, sessionRepository)
	sessionUsecase := usecase.NewSessionUsecase(config.DB, config.Validate, config.TokenMaker, config.ViperConfig, sessionRepository)
	postUsecase := usecase.NewPostUsecase(config.DB, config.Validate, postRepository)
	otpUsecase := usecase.NewOtpUsecase(config.DB, config.Validate, config.EmailSender, config.TokenMaker, config.ViperConfig, otpRepository, registerRepository, userRepository, sessionRepository)

	// setup controller
	authController := controller.NewAuthController(registerUsecase, loginUsecase, sessionUsecase)
	surveyController := controller.NewSurveyController(surveyUsecase)
	userController := controller.NewUserController(userUsecase)
	postController := controller.NewPostController(postUsecase)
	otpController := controller.NewOtpController(otpUsecase)

	// setup middleware
	authMiddleware := middleware.AuthMiddleware(config.TokenMaker, config.ViperConfig)

	routeConfig := router.RouteConfig{
		App:              config.App,
		AuthMiddleware:   authMiddleware,
		AuthController:   authController,
		SurveyController: surveyController,
		UserController:   userController,
		PostController:   postController,
		OtpController:    otpController,
	}

	routeConfig.Setup()
}
