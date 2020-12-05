package models

//Point ...
type Point struct {
	ID            	  int64  `json:"pointid,omitempty"`
	RouteID           int64  `json:"routeid,omitempty"`
	Latitude          string    `json:"latitude"`
	Longitude         string `json:"longitude"`
	PointIndex		  int32  `json:"pointindex,omitempty"`
}

// //PointDecode user for decoding
// type PointDecode struct {
// 	UserID            int32          `json:"userid"`
// 	UserName          sql.NullString `json:"username,omitempty"`
// 	FirstName         string         `json:"firstname"`
// 	LastName          sql.NullString `json:"lastname,omitempty"`
// 	Email             string         `json:"email"`
// 	PhoneNumber       sql.NullInt32  `json:"phonenumber,omitempty"`
// 	EncryptedPassword string         `json:"-"`
// }