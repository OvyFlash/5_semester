package sqlstore

import (
	"neckname/internal/app/models"
)

//CommentRepository ...
type CommentRepository struct {
	store *Store
}

//CreateComment creates new instance of commenta
func (r *CommentRepository) CreateComment(u *models.Comment) (*models.Comment, error) {
	return nil, nil
}
