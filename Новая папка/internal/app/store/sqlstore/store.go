package sqlstore

import (
	"sync"
)

//Store ...
type Store struct {
	mu *sync.Mutex
	//db         *sql.DB
	repository *Repository
	db	string

}

//NewStore return instance of store
func NewStore(databaseName string) *Store {
	return &Store{
		mu: &sync.Mutex{},
		db: databaseName,
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
