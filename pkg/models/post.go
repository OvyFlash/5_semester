package models

import "database/sql"

//Post ...
type Post struct {
	PostID  int64  `json:"postid,omitempty"`
	RouteID int64  `json:"routeid"`
	Date    int64  `json:"date,omitempty"`
	Text    string `json:"text,omitempty"`
	Route   *Route `json:"route,-"`
}

//PostDecode ...
type PostDecode struct {
	PostID  int64          `json:"postid,omitempty"`
	RouteID int64          `json:"routeid"`
	Date    int64          `json:"date,omitempty"`
	Text    sql.NullString `json:"text,omitempty"`
}

//ToPost ...
func (p *PostDecode) ToPost() *Post {
	post := &Post{}
	if p.Text.Valid {
		post.Text = p.Text.String
	}
	post.PostID = p.PostID
	post.RouteID = p.RouteID
	post.Date = p.Date
	return post
}

// Posts: postid bigint
// Posts: routeid bigint
// Posts: date
// Posts: photopath1
// Posts: photopath2
// Posts: photopath3
// Posts: text
