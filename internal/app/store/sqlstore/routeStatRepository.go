package sqlstore

import (
	"neckname/internal/app/models"
)

//RouteStatRepository ...
type RouteStatRepository struct {
	store *Store
}

//CreateRouteStat creates new instance of user
func (r *RouteStatRepository) CreateRouteStat(u *models.RouteStat) (*models.RouteStat, error) {
	return nil, nil
}
