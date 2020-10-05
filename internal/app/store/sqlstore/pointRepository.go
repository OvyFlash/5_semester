package sqlstore

import (
	"neckname/internal/app/models"
)

//PointRepository ...
type PointRepository struct {
	store *Store
}

//CreatePoint creates new instance of user
func (r *PointRepository) CreatePoint(u *models.Point) (*models.Point, error) {
	return nil, nil
}