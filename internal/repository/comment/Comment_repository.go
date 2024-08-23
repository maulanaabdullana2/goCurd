package repository

import (
	CommentModels "fiber-crud/internal/domain/comment"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommentRepository interface {
	CreateComment(comment *CommentModels.Comment) error
	Getcommentproductid(ProductID uuid.UUID, UserID uuid.UUID) ([]CommentModels.Comment, error)
}

type Commentrepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) CommentRepository {
	return &Commentrepository{db: db}
}

func (r *Commentrepository) CreateComment(comment *CommentModels.Comment) error {

	return r.db.Create(comment).Error
}

func (r *Commentrepository) Getcommentproductid(ProductID uuid.UUID, UserID uuid.UUID) ([]CommentModels.Comment, error) {

	var comments []CommentModels.Comment
	if err := r.db.Where("product_id = ? AND user_id = ?", ProductID, UserID).Find(&comments).Error; err != nil {
		return nil, err
	}
	return comments, nil
}
