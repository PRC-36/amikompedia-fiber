package controller

import (
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/middleware"
	"github.com/PRC-36/amikompedia-fiber/domain/usecase"
	"github.com/PRC-36/amikompedia-fiber/shared/token"
	"github.com/PRC-36/amikompedia-fiber/shared/util"
	"github.com/gofiber/fiber/v2"
	"log"
)

type PostController interface {
	Create(ctx *fiber.Ctx) error
	CreateComment(ctx *fiber.Ctx) error
	FindAll(ctx *fiber.Ctx) error
}

type postControllerImpl struct {
	postUsecase usecase.PostUsecase
}

func NewPostController(postUsecase usecase.PostUsecase) PostController {
	return &postControllerImpl{
		postUsecase: postUsecase,
	}
}

func (p *postControllerImpl) Create(ctx *fiber.Ctx) error {
	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	imagePost, _ := ctx.FormFile("image_post")
	content := ctx.FormValue("content")

	requestBody := &request.PostRequest{Content: content}

	result, err := p.postUsecase.CreateNewPost(ctx.UserContext(), requestBody, authPayload.UserID, imagePost)

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

func (p *postControllerImpl) CreateComment(ctx *fiber.Ctx) error {
	requestBody := new(request.PostCommentRequest)

	authPayload := ctx.Locals(middleware.AuthorizationPayloadKey).(*token.Payload)

	err := ctx.BodyParser(requestBody)

	log.Printf("error parser", err)
	if err != nil {
		resp, statusCode := util.ConstructBaseResponse(
			util.BaseResponse{
				Code:   fiber.StatusBadRequest,
				Status: err.Error(),
			},
		)
		return ctx.Status(statusCode).JSON(resp)
	}

	result, err := p.postUsecase.CommentPost(ctx.UserContext(), requestBody, authPayload.UserID)

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

func (p *postControllerImpl) FindAll(ctx *fiber.Ctx) error {

	request := &request.SearchPostRequest{
		Keyword: ctx.Query("keyword", ""),
		Page:    ctx.QueryInt("page", 1),
		Size:    ctx.QueryInt("size", 10),
	}

	result, err := p.postUsecase.FindAllAndSearch(ctx.UserContext(), request)

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
			Code:   fiber.StatusOK,
			Status: "Success",
			Data:   result,
		},
	)

	return ctx.Status(statusCode).JSON(resp)
}
