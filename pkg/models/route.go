package models

//Route ...
type Route struct {
	RouteID int64    `json:"routeid,omitempty"`
	UserID  int32    `json:"userid,omitempty"`
	Date    int64    `json:"date,omitempty"`
	Points  []*Point `json:"points,omitempty"`
	Start   int64    `json:"start,omitempty"`
	Finish  int64    `json:"finish,omitempty"`
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

// //ToRouteStat ...
// func (r *Route) ToRouteStat() *RouteStat {
// 	start, _ := strconv.ParseInt(r.Start, 10, 64)
// 	finish, _ := strconv.ParseInt(r.Finish, 10, 64)
// 	timeOfTraining := time.Unix(start, 0).Sub(time.Unix(finish, 0))
// 	return &RouteStat{

// 	}
// }
