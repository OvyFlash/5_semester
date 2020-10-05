package models

//Follower ...
type Follower struct {
	ID int32
	UserID int32
	FollowedUserID int32
}