package sqlstore

//Repository ...
type Repository struct {
	Comment   *CommentRepository
	Follower  *FollowerRepository
	Point     *PointRepository
	Post      *PostRepository
	Route     *RouteRepository
	RouteStat *RouteStatRepository
	User      *UserRepository
}
