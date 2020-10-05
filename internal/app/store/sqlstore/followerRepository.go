package sqlstore

import (
	"neckname/internal/app/models"
)

//FollowerRepository ...
type FollowerRepository struct {
	store *Store
}

//CreateFollower creates new instance of user
func (r *FollowerRepository) CreateFollower(u *models.Follower) (*models.Follower, error) {
	return nil, nil
}
