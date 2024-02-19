package usecase

import (
	"context"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/response"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"github.com/PRC-36/amikompedia-fiber/domain/repository"
	"github.com/PRC-36/amikompedia-fiber/shared/aws"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"log"
	"math"
	"mime/multipart"
)

type PostUsecase interface {
	CreateNewPost(ctx context.Context, requestData *request.PostRequest, userId string, imgPost *multipart.FileHeader) (*response.PostResponse, error)
	FindAllAndSearch(ctx context.Context, requestData *request.SearchPostRequest) (*response.PostResponses, error)
	CommentPost(ctx context.Context, requestData *request.PostCommentRequest, userId string) (*response.PostResponse, error)
}

type postUsecaseImpl struct {
	DB              *gorm.DB
	Validate        *validator.Validate
	AwsS3           aws.AwsS3Action
	PostRepository  repository.PostRepository
	ImageRepository repository.ImageRepository
}

func NewPostUsecase(DB *gorm.DB, validate *validator.Validate, awsS3 aws.AwsS3Action,
	postRepository repository.PostRepository,
	imageRepository repository.ImageRepository) PostUsecase {
	return &postUsecaseImpl{
		DB:              DB,
		Validate:        validate,
		AwsS3:           awsS3,
		PostRepository:  postRepository,
		ImageRepository: imageRepository,
	}
}

func (p *postUsecaseImpl) CreateNewPost(ctx context.Context, requestData *request.PostRequest, userId string, imgPost *multipart.FileHeader) (*response.PostResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := p.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid username : %+v", err)
		return nil, err
	}

	postEntity := requestData.ToEntity(userId)

	err = p.PostRepository.PostCreate(tx, postEntity)

	if err != nil {
		log.Printf("Failed create post : %+v", err)
		return nil, err
	}

	if imgPost != nil {
		uploaded, err := p.AwsS3.UploadFile(imgPost, aws.ImgPost)
		if err != nil {
			log.Printf("Failed upload avatar image : %+v", err)
			return nil, err
		}

		imgPostEntity := &entity.Image{UserID: postEntity.UserID,
			PostID:    postEntity.ID,
			ImageType: uploaded.ImageType,
			ImageUrl:  uploaded.ImageUrl,
			FilePath:  uploaded.FilePath,
		}

		err = p.ImageRepository.ImageSave(tx, imgPostEntity)
		if err != nil {
			log.Printf("Failed upload save image avatar : %+v", err)
			return nil, err
		}

		postEntity.Images = append(postEntity.Images, *imgPostEntity)
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

func (p *postUsecaseImpl) CommentPost(ctx context.Context, requestData *request.PostCommentRequest, userId string) (*response.PostResponse, error) {
	tx := p.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	err := p.Validate.Struct(requestData)
	if err != nil {
		log.Printf("Invalid username : %+v", err)
		return nil, err
	}

	postEntity := requestData.ToEntity(userId)

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
