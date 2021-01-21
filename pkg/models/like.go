package models

type Like struct {
	LikeID int64 `json:"likeid,omitempty"`
	PostID int64 `json:"postid"`
	UserID int32 `json:"userid,omitempty"`
}


// Likes: likeid
// Likes: postid
// Likes: userid