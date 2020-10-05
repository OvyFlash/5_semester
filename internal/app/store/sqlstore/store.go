package sqlstore

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql" //mysql driver
)

//Store ...
type Store struct {
	db         *sql.DB
	repository *Repository
}

//NewStore return instance of store
func NewStore(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

//Repository creates instance of repository
func (s *Store) Repository() *Repository {
	if s.repository != nil {
		return s.repository
	}

	s.repository = &Repository{
		Comment:   &CommentRepository{s},
		Follower:  &FollowerRepository{s},
		Point:     &PointRepository{s},
		Post:      &PostRepository{s},
		Route:     &RouteRepository{s},
		RouteStat: &RouteStatRepository{s},
		User:      &UserRepository{s},
	}
	return s.repository
}

//store.Repository().UserRepo
