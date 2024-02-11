package controller

import (
	"errors"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/middleware"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
)

type UserController interface {
	Create(ctx *fiber.Ctx) error
	Profile(ctx *fiber.Ctx) error
	Update(ctx *fiber.Ctx) error
}

type userControllerImpl struct {
	userUsecase usecase.UserUsecase
}

func NewUserController(userUsecase usecase.UserUsecase) UserController {
	return &userControllerImpl{
		userUsecase: userUsecase,
	}
}

func (u *userControllerImpl) Create(ctx *fiber.Ctx) error {
	requestBody := new(request.UserRequest)
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

	result, err := u.userUsecase.CreateNewUser(ctx.UserContext(), requestBody)

	if err != nil {
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

func (u *userControllerImpl) Profile(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	result, err := u.userUsecase.ProfileUser(ctx.UserContext(), authPayload.UserID)

	if err != nil {
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
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}

func (u *userControllerImpl) Update(ctx *fiber.Ctx) error {

	//requestBody := new(request.UserUpdateRequest)
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	//err := ctx.BodyParser(requestBody)

	username := ctx.FormValue("username")
	name := ctx.FormValue("name")
	bio := ctx.FormValue("bio")
	requestBody := &request.UserUpdateRequest{Username: username, Name: name, Bio: bio}

	imgAvtr, _ := ctx.FormFile("image_avatar")
	imgHeader, _ := ctx.FormFile("image_header")
	//images := []*multipart.FileHeader{imgAvtr, imgHeader}

	//if err != nil {
	//	resp, statusCode := util.ConstructBaseResponse(
	//		util.BaseResponse{
	//			Code:   fiber.StatusBadRequest,
	//			Status: err.Error(),
	//		},
	//	)
	//	return ctx.Status(statusCode).JSON(resp)
	//}

	result, err := u.userUsecase.UpdateUser(ctx.UserContext(), authPayload.UserID, requestBody, imgAvtr, imgHeader)

	if err != nil {
		if errors.Is(err, util.UserNotFound) {
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
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
