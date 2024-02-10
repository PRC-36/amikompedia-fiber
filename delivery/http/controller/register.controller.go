package controller

import (
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
)

type RegisterController interface {
	Register(ctx *fiber.Ctx) error
}

type registerControllerImpl struct {
	registerUsecase usecase.RegisterUsecase
}

func NewRegisterController(registerUsecase usecase.RegisterUsecase) RegisterController {
	return &registerControllerImpl{registerUsecase: registerUsecase}
}

func (r *registerControllerImpl) Register(ctx *fiber.Ctx) error {
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

	result, err := r.registerUsecase.Register(ctx.UserContext(), requestBody)

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
