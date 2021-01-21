package models

//Follower ...
type Follower struct {
	ID int64 `json:"followid,omitempty"`
	UserID int32 `json:"userid,omitempty"`
	FollowerUserID int32 `json:"follower_userid"`
}