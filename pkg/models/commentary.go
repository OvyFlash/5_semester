package models

//Comment ...
type Comment struct {
	CommentID  int64  `json:"commentid,omitempty"`
	UserID     int32  `json:"userid,omitempty"`
	PostID     int64  `json:"postid"`
	Commentary string `json:"text"`
}

//Comments: id bigint
//Comments: userid int
//Comments: postid bigint
//Comments: commentary text
