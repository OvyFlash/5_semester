package models

//RouteStat ...
type RouteStat struct {
	ID          int64 `json:"id,omitempty"`
	RouteID     int64 `json:"routeid,omitempty"`
	UserID      int32 `json:"userid,omitempty"`
	Workouttime int64 `json:"workouttime,omitempty"`
	Date        int64 `json:"date,omitempty"`
}
