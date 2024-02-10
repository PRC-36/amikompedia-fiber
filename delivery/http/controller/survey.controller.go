package controller

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type SurveyController interface {
	Create(ctx *fiber.Ctx) error
}

type surveyControllerImpl struct {
	surveyUsecase usecase.SurveyUsecase
}

func NewSurveyController(surveyUsecase usecase.SurveyUsecase) SurveyController {
	return &surveyControllerImpl{
		surveyUsecase: surveyUsecase,
	}
}

func (s *surveyControllerImpl) Create(ctx *fiber.Ctx) error {
	requestBody := new(request.SurveyRequest)

	//authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	//log.Println("authPayload", authPayload)

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

	result, err := s.surveyUsecase.Create(ctx.UserContext(), uuid.NewString(), requestBody)

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
