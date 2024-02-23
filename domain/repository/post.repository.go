package repository

import (
	"fmt"
	"github.com/PRC-36/amikompedia-fiber/delivery/http/dto/request"
	"github.com/PRC-36/amikompedia-fiber/domain/entity"
	"gorm.io/gorm"
	"log"
)

type PostRepository interface {
	PostCreate(tx *gorm.DB, value *entity.Post) error
	PostFindAll(db *gorm.DB, request *request.SearchPostRequest) ([]entity.Post, int64, error)
	FindByID(db *gorm.DB, id string) (*entity.Post, error)
	FindALlCommentsByID(db *gorm.DB, postID []string) ([]entity.Post, error)
	UpdateTotalComments(tx *gorm.DB, postID string) error
	FilterPost(request *request.SearchPostRequest) func(tx *gorm.DB) *gorm.DB
}

type postRepositoryImpl struct {
}

func NewPostRepository() PostRepository {
	return &postRepositoryImpl{}
}

func (p *postRepositoryImpl) PostCreate(tx *gorm.DB, value *entity.Post) error {
	result := tx.Create(value)
	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when creating post : %v", result.Error))
		return result.Error
	}

	find := tx.Preload("User").Preload("User.Images", "image_type NOT LIKE ?", "POST").First(value, value.ID)

	if find.Error != nil {
		log.Println(fmt.Sprintf("Error when creating post : %v", find.Error))
		return find.Error
	}

	return nil

}

func (p *postRepositoryImpl) PostFindAll(db *gorm.DB, request *request.SearchPostRequest) ([]entity.Post, int64, error) {
	var posts []entity.Post

	err := db.Scopes(p.FilterPost(request)).Preload("Images").Preload("User").Preload("User.Images", "image_type NOT LIKE ?", "POST").Offset((request.Page - 1) * request.Size).Limit(request.Size).Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64 = 0
	if err := db.Model(&entity.Post{}).Scopes(p.FilterPost(request)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return posts, total, nil
}

func (p *postRepositoryImpl) FindByID(db *gorm.DB, id string) (*entity.Post, error) {
	var post entity.Post
	result := db.Preload("Images").Preload("User").Preload("User.Images", "image_type NOT LIKE ?", "POST").First(&post, "id = ?", id)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find post by id : %v", result.Error))
		return nil, result.Error
	}

	return &post, nil
}

func (p *postRepositoryImpl) FindALlCommentsByID(db *gorm.DB, postID []string) ([]entity.Post, error) {
	var posts []entity.Post
	result := db.Preload("Images").Preload("User").Preload("User.Images", "image_type NOT LIKE ?", "POST").Find(&posts, "ref_post_id IN ?", postID)

	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when find post by id : %v", result.Error))
		return nil, result.Error
	}

	return posts, nil
}

func (p *postRepositoryImpl) UpdateTotalComments(tx *gorm.DB, postID string) error {
	result := tx.Model(&entity.Post{}).Where("id = ?", postID).Update("total_comments", gorm.Expr("total_comments + ?", 1))
	if result.Error != nil {
		log.Println(fmt.Sprintf("Error when updating total comments : %v", result.Error))
		return result.Error
	}

	return nil
}

func (p *postRepositoryImpl) FilterPost(request *request.SearchPostRequest) func(tx *gorm.DB) *gorm.DB {
	return func(tx *gorm.DB) *gorm.DB {

		if keyword := request.Keyword; keyword != "" {
			keyword = "%" + keyword + "%"
			tx = tx.Where("content LIKE ?", keyword)
		}

		if postId := request.PostID; postId != "" {
			tx = tx.Where("ref_post_id = ?", postId)
		}
		return tx
	}
}
