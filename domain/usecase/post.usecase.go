package usecase

import (
	"context"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"math"
)

type PostUsecase interface {
	CreateNewPost(ctx context.Context, requestData *request.PostRequest) (*response.PostResponse, error)
	FindAllAndSearch(ctx context.Context, requestData *request.SearchPostRequest) (*response.PostResponses, error)
}

type postUsecaseImpl struct {
	DB             *gorm.DB
	Validate       *validator.Validate
	PostRepository repository.PostRepository
}

func NewPostUsecase(DB *gorm.DB, validate *validator.Validate, postRepository repository.PostRepository) PostUsecase {
	return &postUsecaseImpl{
		DB:             DB,
		Validate:       validate,
		PostRepository: postRepository,
	}
}

func (p *postUsecaseImpl) CreateNewPost(ctx context.Context, requestData *request.PostRequest) (*response.PostResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := p.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid username : %+v", err)
		return nil, err
	}

	postEntity := requestData.ToEntity()

	err = p.PostRepository.PostCreate(tx, postEntity)

	if err != nil {
		log.Printf("Failed create post : %+v", err)
		return nil, err
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return postEntity.ToPostResponse(), nil
}

func (p *postUsecaseImpl) FindAllAndSearch(ctx context.Context, requestData *request.SearchPostRequest) (*response.PostResponses, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := p.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid username : %+v", err)
		return nil, err
	}

	posts, total, err := p.PostRepository.PostFindAll(tx, requestData)

	if err != nil {
		log.Printf("Error getting posts : %+v", err)
		return nil, err
	}

	resultPaging := &response.PostPageMetaData{
		Page:      requestData.Page,
		Size:      requestData.Size,
		TotalItem: total,
		TotalPage: int64(math.Ceil(float64(total) / float64(requestData.Size))),
	}

	err = tx.Commit().Error
	if err != nil {
		return nil, err
	}

	return entity.ToPostResponses(posts, resultPaging), nil
}
