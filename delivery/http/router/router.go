package router

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/controller"
	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	App              *fiber.App
	AuthMiddleware   fiber.Handler
	AuthController   controller.AuthController
	SurveyController controller.SurveyController
	UserController   controller.UserController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	// Register
	c.App.Post("/api/v1/auth/_register", c.AuthController.Register)

	// Login
	c.App.Post("/api/v1/auth/_login", c.AuthController.Login)

	// Renew Access Token
	c.App.Post("/api/v1/auth/_renewtoken", c.AuthController.RenewAccessToken)

	// User TEMP for testing
	c.App.Post("/api/v1/users", c.UserController.Create)
}

func (c *RouteConfig) SetupAuthRoute() {
	c.App.Use(c.AuthMiddleware)

	// Survey
	c.App.Post("/api/v1/surveys", c.SurveyController.Create)

	// User
	c.App.Get("/api/v1/users/profile", c.UserController.Profile)
	c.App.Patch("/api/v1/users", c.UserController.Update)
}
