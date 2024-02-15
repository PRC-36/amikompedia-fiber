package controller

import (
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
)

type OtpController interface {
	OtpValidation(ctx *fiber.Ctx) error
	ResendOtp(ctx *fiber.Ctx) error
}

type otpControllerImpl struct {
	otpUsecase usecase.OtpUsecase
}

func NewOtpController(otpUsecase usecase.OtpUsecase) OtpController {
	return &otpControllerImpl{otpUsecase: otpUsecase}
}

func (o otpControllerImpl) OtpValidation(ctx *fiber.Ctx) error {
	requestBody := new(request.OtpValidateRequest)
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

	result, err := o.otpUsecase.OtpValidation(ctx.UserContext(), requestBody, userAgent, clientIP)

	if err != nil {
		if errors.Is(err, util.RefCodeNotFound) || errors.Is(err, util.UserRegisterNotFound) {
			util.ConstructBaseResponse(
				util.BaseResponse{
					Code:   fiber.StatusNotFound,
					Status: err.Error(),
				})
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
			Code:   fiber.StatusCreated,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (o otpControllerImpl) ResendOtp(ctx *fiber.Ctx) error {
	requestBody := new(request.OtpSendRequest)
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

	err = o.otpUsecase.ResendOtp(ctx.UserContext(), requestBody)

	if err != nil {
		if errors.Is(err, util.RefCodeNotFound) {
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
			Code:   fiber.StatusCreated,
			Status: "Success",
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
