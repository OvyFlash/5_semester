package models

//Comment ...
type Comment struct {
	ID int32
	PostID int32
	UserID int32
	Commentary string
}