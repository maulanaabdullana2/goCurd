package commentUsecase

import (
	CommentModels "fiber-crud/internal/domain/comment"
	repository "fiber-crud/internal/repository/comment"

	"github.com/google/uuid"
)

type CommentUsecase interface {
	CreateComment(comment *CommentModels.Comment) error
	Getcommentproductid(ProductID uuid.UUID, UserID uuid.UUID) ([]CommentModels.Comment, error)
}

type commentUsecase struct {
	commentRepository repository.CommentRepository
}

func NewCommentUsecase(repo repository.CommentRepository) CommentUsecase {
	return &commentUsecase{commentRepository: repo}
}

func (r *commentUsecase) CreateComment(comment *CommentModels.Comment) error {
	return r.commentRepository.CreateComment(comment)
}
func (r *commentUsecase) Getcommentproductid(ProductID uuid.UUID, UserID uuid.UUID) ([]CommentModels.Comment, error) {
	return r.commentRepository.Getcommentproductid(ProductID, UserID)
}
