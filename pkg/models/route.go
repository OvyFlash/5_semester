package models

//Route ...
type Route struct {
	RouteID int64 `json:"routeid,omitempty"`
	UserID  int32 `json:"userid,omitempty"`
	Date    int64 `json:"date,omitempty"`
}

// //RouteDecode route for decoding
// type RouteDecode struct {
// 	UserID            int32          `json:"userid"`
// 	UserName          sql.NullString `json:"username,omitempty"`
// 	FirstName         string         `json:"firstname"`
// 	LastName          sql.NullString `json:"lastname,omitempty"`
// 	Email             string         `json:"email"`
// 	PhoneNumber       sql.NullInt32  `json:"phonenumber,omitempty"`
// 	EncryptedPassword string         `json:"-"`
// }
