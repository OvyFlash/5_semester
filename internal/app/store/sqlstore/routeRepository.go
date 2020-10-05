package sqlstore

import (
	"neckname/internal/app/models"
)

//RouteRepository ...
type RouteRepository struct {
	store *Store
}

//CreateRoute creates new instance of user
func (r *RouteRepository) CreateRoute(u *models.Route) (*models.Route, error) {
	return nil, nil
}
