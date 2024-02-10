package app

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/controller"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/middleware"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/router"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
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
}

func Bootstrap(config *BootstrapConfig) {

	// setup repositories
	registerRepository := repository.NewRegisterRepository()
	surveyRepository := repository.NewSurveyRepository()

	// setup usecases
	registerUsecase := usecase.NewRegisterUsecase(config.DB, config.Validate, config.EmailSender, registerRepository, repository.NewOtpRepository())
	surveyUsecase := usecase.NewSurveyUsecase(config.DB, config.Validate, surveyRepository)

	// setup controller
	registerController := controller.NewRegisterController(registerUsecase)
	surveyController := controller.NewSurveyController(surveyUsecase)

	// setup middleware
	authMiddleware := middleware.AuthMiddleware(config.TokenMaker)

	routeConfig := router.RouteConfig{
		App:                config.App,
		AuthMiddleware:     authMiddleware,
		RegisterController: registerController,
		SurveyController:   surveyController,
	}

	routeConfig.Setup()
}
