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
	//db.Close()
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
	u.EncryptPassword()
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

//CheckUserEmail checks if user with such email presents
func (mydb *Sqldatabase) CheckUserEmail(email string) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT NOT NULL, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return err
	}

	var id int32
	if err := mydb.db.QueryRow("SELECT userid FROM Users WHERE email = ?",
		email).Scan(&id); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return errors.New("User with such email already exists")
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

//CreateRoute ...
func (mydb *Sqldatabase) CreateRoute(u *models.User, route *models.Route) (*models.Route, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL, start INTEGER NOT NULL, finish INTEGER NOT NULL)"); err != nil {
		return nil, err
	}

	date := time.Now().Unix()
	result, err := mydb.db.Exec("INSERT INTO Route(userid, date, start, finish) VALUES(?, ?, ?, ?)",
		u.UserID,
		date,
		route.Start,
		route.Finish)
	if err != nil {
		return nil, err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	route.RouteID = ID
	route.UserID = u.UserID
	route.Date = date

	if err := mydb.AddPoints(route); err != nil {
		return nil, err
	}

	route.Points = nil
	return route, nil
}

//AddPoints ...
func (mydb *Sqldatabase) AddPoints(route *models.Route) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Points(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, latitude TEXT NOT NULL, longitude TEXT NOT NULL, pointindex INT NOT NULL)"); err != nil {
		return err
	}
	for i, point := range route.Points {
		if _, err := mydb.db.Exec("INSERT INTO Points(routeid, latitude, longitude, pointindex) VALUES(?, ?, ?, ?)",
			route.RouteID, point.Latitude, point.Longitude, i+1); err != nil {
			return err
		}
	}

	return nil
}

//GetRouteByID ...
func (mydb *Sqldatabase) GetRouteByID(routeid int64) (*models.Route, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL, start INTEGER NOT NULL, finish INTEGER NOT NULL)"); err != nil {
		return nil, err
	}

	route := &models.Route{}
	err := mydb.db.QueryRow("SELECT * FROM Route WHERE routeid = ?",
		routeid).Scan(&route.RouteID, &route.UserID, &route.Date, &route.Start, &route.Finish)
	if err != nil {
		return nil, err
	}
	route.Points, err = mydb.GetPointsByRouteID(routeid)
	if err != nil {
		return nil, err
	}

	return route, nil
}

//GetUserRoutes ...
func (mydb *Sqldatabase) GetUserRoutes(userID int32) ([]*models.Route, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Route(routeid INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, date INTEGER NOT NULL, start INTEGER NOT NULL, finish INTEGER NOT NULL)"); err != nil {
		return nil, err
	}

	rows, err := mydb.db.Query("SELECT * FROM Route WHERE userid = ?",
		userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var routes []*models.Route

	for rows.Next() {
		route := &models.Route{}
		rows.Scan(&route.RouteID, &route.UserID, &route.Date, &route.Start, &route.Finish)
		if err != nil {
			return nil, err
		}
		route.Points, err = mydb.GetPointsByRouteID(route.RouteID)
		if err != nil {
			return nil, err
		}
		routes = append(routes, route)
	}

	return routes, nil
}

//GetPointsByRouteID ...
func (mydb *Sqldatabase) GetPointsByRouteID(routeID int64) ([]*models.Point, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Points(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, latitude TEXT NOT NULL, longitude TEXT NOT NULL, pointindex INT NOT NULL)"); err != nil {
		return nil, err
	}
	rows, err := mydb.db.Query("SELECT * FROM Points WHERE routeid = ?", routeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []*models.Point
	for rows.Next() {
		point := &models.Point{}
		err := rows.Scan(
			&point.ID,
			&point.RouteID,
			&point.Latitude,
			&point.Longitude,
			&point.PointIndex)
		if err != nil {
			return nil, err
		}
		points = append(points, point)
	}

	return points, err
}

//GetRouteStatByRouteID ...
func (mydb *Sqldatabase) GetRouteStatByRouteID(id int64) (*models.RouteStat, error) {
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

//FollowUser ...
func (mydb *Sqldatabase) FollowUser(followed *models.User, follower *models.User) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Followers(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, follower_userid INTEGER NOT NULL)"); err != nil {
		return err
	}
	count := 0
	if err := mydb.db.QueryRow("SELECT COUNT(*) FROM Followers WHERE userid = ? AND follower_userid = ?").Scan(
		&count,
	); err != nil {
		return err
	}

	if count > 0 {
		return errors.New("User already following")
	}

	_, err := mydb.db.Exec("INSERT INTO Followers(userid, follower_userid) VALUES(?, ?)", followed.UserID, follower.UserID)
	if err != nil {
		return err
	}
	return nil
}

//UnfollowUser ...
func (mydb *Sqldatabase) UnfollowUser(followed *models.User, follower *models.User) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Followers(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, follower_userid INTEGER NOT NULL)"); err != nil {
		return err
	}
	count := 0
	if err := mydb.db.QueryRow("SELECT COUNT(*) FROM Followers WHERE userid = ? AND follower_userid = ?").Scan(
		&count,
	); err != nil {
		return err
	}

	if count == 0 {
		return errors.New("User is not following")
	}

	_, err := mydb.db.Exec("DELETE FROM Followers where userid = ? AND follower_userid = ?", followed.UserID, follower.UserID)
	if err != nil {
		return err
	}
	return nil
}

//CreatePost ...
func (mydb *Sqldatabase) CreatePost(post *models.Post, userID int32) (*models.Post, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Posts(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, date INTEGER NOT NULL, post_text TEXT)"); err != nil {
		return nil, err
	}

	date := time.Now().Unix()
	post.Date = date
	result, err := mydb.db.Exec("INSERT INTO Posts(routeid, date, post_text) VALUES(?, ?, ?)",
		post.RouteID, post.Date, post.Text)
	if err != nil {
		return nil, err
	}
	post.Route, err = mydb.GetRouteByID(post.RouteID)
	if err != nil {
		return nil, err
	}
	post.PostID, _ = result.LastInsertId()
	return post, nil
}

//UpdatePost ...
func (mydb *Sqldatabase) UpdatePost(post *models.Post, userID int32) (*models.Post, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Posts(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, date INTEGER NOT NULL, post_text TEXT)"); err != nil {
		return nil, err
	}

	var userid int32 = 0
	postDecode := &models.PostDecode{}
	err := mydb.db.QueryRow("SELECT Posts.id, Posts.routeid, Posts.date, Posts.post_text, Route.userid FROM Posts JOIN Route ON Route.routeid = Posts.routeid WHERE Posts.id = ?", post.PostID).Scan(
		&postDecode.PostID,
		&postDecode.RouteID,
		&postDecode.Date,
		&postDecode.Text,
		&userid)
	if err != nil {
		return nil, err
	}
	if userid != userID {
		return nil, errors.New("You cant edit this post")
	}

	post.RouteID = postDecode.RouteID
	post.Date = postDecode.Date
	post.PostID = postDecode.PostID
	if postDecode.Text.Valid {
		if post.Text == "" && post.Text != postDecode.Text.String {
			post.Text = postDecode.Text.String
		}
	}
	if _, err := mydb.db.Exec("UPDATE Posts SET post_text = ? WHERE id = ?", post.Text, post.PostID); err != nil {
		return nil, err
	}

	post.Route, err = mydb.GetRouteByID(post.RouteID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

//DeletePost ...
func (mydb *Sqldatabase) DeletePost(post *models.Post, userID int32) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Posts(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, date INTEGER NOT NULL, post_text TEXT)"); err != nil {
		return err
	}

	var userid int32 = 0
	err := mydb.db.QueryRow("SELECT Route.userid FROM Posts JOIN Route ON Route.routeid = Posts.routeid WHERE Posts.id = ?", post.PostID).Scan(
		&userid)
	if err != nil {
		return err
	}
	if userid != userID {
		return errors.New("You cant delete this post")
	}
	if _, err := mydb.db.Exec("DELETE FROM Posts WHERE id = ?", post.PostID); err != nil {
		return err
	}

	return nil
}

//GetPost ...
func (mydb *Sqldatabase) GetPost(post *models.Post) (*models.Post, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Posts(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, date INTEGER NOT NULL, post_text TEXT)"); err != nil {
		return nil, err
	}

	postDecode := &models.PostDecode{}
	err := mydb.db.QueryRow("SELECT * FROM Posts WHERE id = ?", post.PostID).Scan(
		&postDecode.PostID,
		&postDecode.RouteID,
		&postDecode.Date,
		&postDecode.Text)
	if err != nil {
		return nil, err
	}

	post.RouteID = postDecode.RouteID
	post.Date = postDecode.Date
	post.PostID = postDecode.PostID
	if postDecode.Text.Valid {
		post.Text = postDecode.Text.String

	}

	post.Route, err = mydb.GetRouteByID(post.RouteID)
	if err != nil {
		return nil, err
	}

	return post, nil
}

//AddComment ...
func (mydb *Sqldatabase) AddComment(c *models.Comment) (*models.Comment, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Comment(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, postid INTEGER NOT NULL, commentary TEXT NOT NULL)"); err != nil {
		return nil, err
	}

	result, err := mydb.db.Exec("INSERT INTO Comment(userid, postid, commentary) VALUES(?, ?, ?)",
		c.UserID, c.PostID, c.Commentary)
	if err != nil {
		return nil, err
	}
	c.CommentID, _ = result.LastInsertId()
	return c, nil
}

//UpdateComment ...
func (mydb *Sqldatabase) UpdateComment(c *models.Comment) (*models.Comment, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Comment(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, postid INTEGER NOT NULL, commentary TEXT NOT NULL)"); err != nil {
		return nil, err
	}

	_, err := mydb.db.Exec("UPDATE Comment SET commentary = ? WHERE postid = ? AND userid = ? AND id = ?", c.Commentary, c.PostID, c.UserID, c.CommentID)
	if err != nil {
		return nil, err
	}

	return c, nil
}

//DeleteComment ...
func (mydb *Sqldatabase) DeleteComment(c *models.Comment) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Comment(id INTEGER PRIMARY KEY AUTOINCREMENT, userid INTEGER NOT NULL, postid INTEGER NOT NULL, commentary TEXT NOT NULL)"); err != nil {
		return err
	}

	_, err := mydb.db.Exec("DELETE FROM Comment WHERE postid = ? AND userid = ? AND id = ?", c.PostID, c.UserID, c.CommentID)
	if err != nil {
		return err
	}

	return nil
}

//Like ...
func (mydb *Sqldatabase) Like(l *models.Like) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Like(id INTEGER PRIMARY KEY AUTOINCREMENT, postid INTEGER NOT NULL, userid INTEGER NOT NULL)"); err != nil {
		return err
	}

	count := 0
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM Like WHERE postid = ? AND userid = ?", l.PostID, l.UserID).Scan(&count)
	if err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Already liked")
	}
	_, err = mydb.db.Exec("INSERT INTO Like(postid, userid) VALUES(?, ?)", l.PostID, l.UserID)
	if err != nil {
		return err
	}
	return nil
}

//RemoveLike ...
func (mydb *Sqldatabase) RemoveLike(l *models.Like) error {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Like(id INTEGER PRIMARY KEY AUTOINCREMENT, postid INTEGER NOT NULL, userid INTEGER NOT NULL)"); err != nil {
		return err
	}

	count := 0
	err := mydb.db.QueryRow("SELECT COUNT(*) FROM Like WHERE postid = ? AND userid = ?", l.PostID, l.UserID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New("Not liked")
	}
	_, err = mydb.db.Exec("DELETE FROM Like WHERE postid = ? AND userid = ?", l.PostID, l.UserID)
	if err != nil {
		return err
	}
	return nil
}

//GetAllPosts ...
func (mydb *Sqldatabase) GetAllPosts() ([]*models.Post, error) {
	if _, err := mydb.db.Exec("CREATE TABLE IF NOT EXISTS Posts(id INTEGER PRIMARY KEY AUTOINCREMENT, routeid INTEGER NOT NULL, date INTEGER NOT NULL, post_text TEXT)"); err != nil {
		return nil, err
	}

	rows, err := mydb.db.Query("SELECT * FROM Posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post
	for rows.Next() {
		postD := &models.PostDecode{}
		err := rows.Scan(
			&postD.PostID,
			&postD.RouteID,
			&postD.Date,
			&postD.Text,
		)
		if err != nil {
			return nil, err
		}
		p := postD.ToPost()
		p.Route, err = mydb.GetRouteByID(p.RouteID)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, nil
}
