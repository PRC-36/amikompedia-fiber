package controller

import (
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
)

type AuthController interface {
	Register(ctx *fiber.Ctx) error
	Login(ctx *fiber.Ctx) error
	RenewAccessToken(ctx *fiber.Ctx) error
}

type authControllerImpl struct {
	registerUsecase usecase.RegisterUsecase
	loginUsecase    usecase.LoginUsecase
	sessionUsecase  usecase.SessionUsecase
}

func NewAuthController(registerUsecase usecase.RegisterUsecase, loginUsecase usecase.LoginUsecase, sessionUsecase usecase.SessionUsecase) AuthController {
	return &authControllerImpl{
		registerUsecase: registerUsecase,
		loginUsecase:    loginUsecase,
		sessionUsecase:  sessionUsecase,
	}
}

func (a *authControllerImpl) Register(ctx *fiber.Ctx) error {
	requestBody := new(request.UserRegisterRequest)
	err := ctx.BodyParser(requestBody)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	result, err := a.registerUsecase.Register(ctx.UserContext(), requestBody)

	if err != nil {
		if errors.Is(err, util.EmailAlreadyUsed) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		} else {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusBadRequest,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusCreated,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (a *authControllerImpl) Login(ctx *fiber.Ctx) error {
	requestBody := new(request.LoginRequest)
	userAgent := ctx.Get("User-Agent")
	clientIP := ctx.IP()

	err := ctx.BodyParser(requestBody)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	loginResult, err := a.loginUsecase.Login(ctx.UserContext(), requestBody, userAgent, clientIP)

	if err != nil {
		if errors.Is(err, util.InvalidPassword) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusUnauthorized,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		if errors.Is(err, util.UsernameOrEmailNotFound) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusOK,
			Status: "Success",
			Data:   loginResult,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (a *authControllerImpl) RenewAccessToken(ctx *fiber.Ctx) error {
	requestBody := new(request.TokenRenewRequest)

	err := ctx.BodyParser(requestBody)

	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	result, err := a.sessionUsecase.RenewAccessToken(ctx.UserContext(), requestBody)
	if err != nil {
		if errors.Is(err, util.SessionNotFound) {
			resp, statusCode := util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				},
			)
			return ctx.Status(statusCode).JSON(resp)
		}

		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	resp, statusCode := util.ConstructBaseResponse(
		util.BaseResponse{
			Code:   fiber.StatusOK,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
