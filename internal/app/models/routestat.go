package models

//RouteStat ...
type RouteStat struct {
	ID int32
	RouteID int32
	UserID int32
	WorkoutTime int16
	AverageSpeed int16
	BurnedFats int16
}