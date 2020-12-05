package sqlitedatabase

import (
	"database/sql"
	"errors"
	"neckname/internal/config"
	"neckname/pkg/models"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

//Sqldatabase ...
type Sqldatabase struct {
	db   *sql.DB
	name string
}

//MakeSqldatabase ...
func MakeSqldatabase(config *config.Config) *Sqldatabase {

	db, err := sql.Open("sqlite3", config.DatabaseURL)
	if err != nil {
		panic(err)
	}
	if err := db.Ping(); err != nil {
		panic(err)
	}
	mydb := &Sqldatabase{
		db:   db,
		name: config.DatabaseURL,
	}
	return mydb
}

//CreateUser ...
func (mydb *Sqldatabase) CreateUser(u *models.User) (*models.User, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return nil, err
	}

	exists := 0
	if err := mydb.db.QueryRow("SELECT COUNT(*) FROM Users WHERE email = ?", u.Email).Scan(&exists); err != nil {
		return nil, err
	}

	if exists == 0 {
		u.EncryptPassword()
		result, err := mydb.db.Exec("INSERT INTO Users(username, firstname, lastname, phone_number, email, encrypted_password) VALUES(?, ?, ?, ?, ?, ?)",
			u.UserName,
			u.FirstName,
			u.LastName,
			u.PhoneNumber,
			u.Email,
			u.EncryptedPassword)
		if err != nil {
			return nil, err
		}

		ID, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}
		u.UserID = int32(ID)

		return u, nil
	}
	return nil, errors.New("User already exists")
}

//GetUserByID ...
func (mydb *Sqldatabase) GetUserByID(userID int32) (*models.User, error) {
	u := &models.UserDecode{}

	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return nil, err
	}

	if err := mydb.db.QueryRow(
		"SELECT * FROM Users WHERE userid = ?", userID,
	).Scan(
		&u.UserID,    //
		&u.UserName,  //
		&u.FirstName, //
		&u.LastName,  //
		&u.Email,     //
		&u.EncryptedPassword,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("User not found")
		}
		return nil, err
	}

	user := &models.User{
		UserID:    u.UserID,
		Email:     u.Email,
		FirstName: u.FirstName,
	}
	if u.UserName.Valid {
		user.UserName = u.UserName.String
	}
	if u.LastName.Valid {
		user.LastName = u.LastName.String
	}
	if u.PhoneNumber.Valid {
		user.PhoneNumber = int(u.PhoneNumber.Int32)
	}

	return user, nil
}

//GetUserByEmail ...
func (mydb *Sqldatabase) GetUserByEmail(email string) (*models.User, error) {
	u := &models.UserDecode{}

	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return nil, err
	}

	if err := mydb.db.QueryRow(
		"SELECT * FROM Users WHERE email = ?", email,
	).Scan(
		&u.UserID,    //
		&u.UserName,  //
		&u.FirstName, //
		&u.LastName,  //
		&u.Email,     //
		&u.EncryptedPassword,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Email or password does not match")
		}
		return nil, err
	}

	user := &models.User{
		UserID:            u.UserID,
		Email:             u.Email,
		FirstName:         u.FirstName,
		EncryptedPassword: u.EncryptedPassword,
	}
	if u.UserName.Valid {
		user.UserName = u.UserName.String
	}
	if u.LastName.Valid {
		user.LastName = u.LastName.String
	}
	if u.PhoneNumber.Valid {
		user.PhoneNumber = int(u.PhoneNumber.Int32)
	}

	return user, nil
}

//UpdateUser ...
func (mydb *Sqldatabase) UpdateUser(u *models.User) error {

	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return err
	}

	if _, err := mydb.db.Exec("UPDATE Users SET username = ?, firstname = ?, lastname = ?, email = ?, encrypted_password = ?, phone_number = ? WHERE userid = ?",
		u.UserName,
		u.FirstName,
		u.LastName,
		u.Email,
		u.EncryptedPassword,
		u.PhoneNumber,
		u.UserID,
	); err != nil {
		return err
	}
	return nil
}

//DeleteUser ...
func (mydb *Sqldatabase) DeleteUser(userID int32) error {

	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return err
	}

	if _, err := mydb.db.Exec("DELETE FROM Users WHERE userid = ?",
		userID,
	); err != nil {
		return err
	}
	return nil
}

// //FollowUser ...
// func (mydb *Sqldatabase) FollowUser(u *models.User) (error) {

// 	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Followers(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, followed_userid INTEGER NOT NULL)"); err != nil {
// 		return err
// 	}

// 	count := 0
// 	if _, err := mydb.db.Exec("SELECT * FROM Followers WHERE userid = ? AND followed_userid = ?",
// 			userID,
// 		); err != nil {
// 			return err
// 		}

// 	if _, err := mydb.db.Exec("INSERT INTO Followers(userid, followed_userid) VALUES(?, ?)",
// 			userID,
// 		); err != nil {
// 			return err
// 		}
// 	return nil
// }

//CreateRoute ...
func (mydb *Sqldatabase) CreateRoute(u *models.User) (*models.Route, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL"); err != nil {
		return nil, err
	}

	date := time.Now().Unix()
	result, err := mydb.db.Exec("INSERT INTO Route(userid, date) VALUES(?, ?)",
		u.UserID,
		date)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	route := &models.Route{
		RouteID: ID,
		UserID:  u.UserID,
		Date:    date,
	}

	return route, nil
}

//GetRouteByID ...
func (mydb *Sqldatabase) GetRouteByID(routeid int64) (*models.Route, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL"); err != nil {
		return nil, err
	}

	route := &models.Route{}
	err := mydb.db.QueryRow("SELECT * FROM Route WHERE routeid = ?",
		routeid).Scan(&route.RouteID, &route.UserID, &route.Date)
	if err != nil {
		return nil, err
	}

	return route, nil
}

// //GetRouteByID ...
// func (mydb *Sqldatabase) GetRouteByID(routeid int64) (*models.Route, error) {
// 	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL"); err != nil {
// 		return nil, err
// 	}

// 	route := &models.Route{}
// 	err := mydb.db.QueryRow("SELECT * FROM Route WHERE routeid = ?",
// 		routeid).Scan(&route.RouteID, &route.UserID, &route.Date)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return route, nil
// }

/*//Point ...
type Point struct {
	ID            	  int64  `json:"pointid,omitempty"`
	RouteID           int64  `json:"routeid,omitempty"`
	Latitude          string    `json:"latitude"`
	Longitude         string `json:"longitude"`
	PointIndex		  int32  `json:"pointindex,omitempty"`
}*/

//AddPoint ...
func (mydb *Sqldatabase) AddPoint(point *models.Point) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Points(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, latitude TEXT NOT NULL, longitude TEXT NOT NULL, pointindex INT NOT NULL)"); err != nil {
		return err
	}

	count := 0
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM Points WHERE routeid = ?",
		point.RouteID).Scan(&count)
	if err != nil {
		return err
	}
	if _, err = mydb.db.Exec("INSERT INTO Points(routeid, latitude, longitude, pointindex) VALUES(?, ?, ?, ?)",
		point.RouteID, point.Latitude, point.Longitude, count); err != nil {
		return err
	}
	return nil
}

//FinishRoute ...
func (mydb *Sqldatabase) FinishRoute(point *models.Point, userID int32) (*models.RouteStat, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Points(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, latitude TEXT NOT NULL, longitude TEXT NOT NULL, pointindex INT NOT NULL)"); err != nil {
		return nil, err
	}
	count := 0
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM Points WHERE routeid = ?",
		point.RouteID).Scan(&count)
	if err != nil {
		return nil, err
	}
	if _, err = mydb.db.Exec("INSERT INTO Points(routeid, latitude, longitude, pointindex) VALUES(?, ?, ?, ?)",
		point.RouteID, point.Latitude, point.Longitude, count); err != nil {
		return nil, err
	}

	/*
	type RouteStat struct {
	ID          int64 `json:"routeStat,omitempty"`
	RouteID     int64 `json:"routeid,omitempty"`
	UserID      int32 `json:"userid,omitempty"`
	Workouttime
	Date        int64 `json:"date,omitempty"`
	}
	*/
	//creating routestats
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS RouteStats(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, userid INTEGER NOT NULL, workouttime INTEGER NOT NULL, date INTEGER NOT NULL)"); err != nil {
		return nil, err
	}

	route, err := mydb.GetRouteByID(point.RouteID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	workouttime := now.Sub(time.Unix(route.Date, 0))

	result, err := mydb.db.Exec("INSERT INTO RouteStats(routeid, userid, workouttime, date) VALUES(?, ?, ?, ?)",
		point.RouteID, userID, workouttime, now)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	} 
	routestat, err := mydb.GetRouteStatByRouteID(ID)
	if err != nil {
		return nil, err
	}

	return routestat, nil
}

//GetRouteStatByRouteID ...
func(mydb *Sqldatabase)GetRouteStatByRouteID(id int64) (*models.RouteStat, error){
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS RouteStats(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, userid INTEGER NOT NULL, workouttime INTEGER NOT NULL, date INTEGER NOT NULL)"); err != nil {
		return nil, err
	}

	routestat := &models.RouteStat{}
	if err := mydb.db.QueryRow("SELECT * FROM RouteStats WHERE id = ?", id).Scan(
		&routestat.ID,
		&routestat.RouteID,
		&routestat.UserID,
		&routestat.Workouttime,
		&routestat.Date,
	); err != nil {
		return nil, err
	}
	return routestat, nil
}