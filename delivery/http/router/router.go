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
	PostController   controller.PostController
	OtpController    controller.OtpController
}

func (c *RouteConfig) Setup() {
	c.SetupGuestRoute()
	c.SetupAuthRoute()
}

func (c *RouteConfig) SetupGuestRoute() {

	// Register
	c.App.Post("/api/v1/auth/_register", c.AuthController.Register)

	// Otp
	c.App.Post("/api/v1/otp/_validation", c.OtpController.OtpValidation)
	c.App.Post("/api/v1/otp/_resend", c.OtpController.ResendOtp)

	// Login
	c.App.Post("/api/v1/auth/_login", c.AuthController.Login)

	// Forgot Password
	c.App.Post("/api/v1/users/_forgot-password", c.UserController.ForgotPassword)
	c.App.Patch("/api/v1/users/_reset-password", c.UserController.ResetPassword)

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
	c.App.Patch("/api/v1/users/update-password", c.UserController.UpdatePassword)
	c.App.Post("/api/v1/users/follow", c.UserController.Follow)
	c.App.Delete("/api/v1/users/follow/:follow_id", c.UserController.Unfollow)

	// Post
	c.App.Post("/api/v1/posts", c.PostController.Create)
	c.App.Get("/api/v1/posts", c.PostController.FindAll)
	c.App.Post("/api/v1/posts/comments", c.PostController.CreateComment)
	c.App.Get("/api/v1/posts/:id", c.PostController.DetailPost)
}
