package sqlstore

import (
	"neckname/internal/app/models"
)

//PostRepository ...
type PostRepository struct {
	store *Store
}

//CreatePost creates new instance of user
func (r *PostRepository) CreatePost(u *models.Post) (*models.Post, error) {
	return nil, nil
}
