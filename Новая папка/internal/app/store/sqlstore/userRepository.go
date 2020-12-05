package sqlstore

import (
	"database/sql"
	"neckname/internal/app/models"
	"neckname/internal/app/store"
	_ "github.com/mattn/go-sqlite3"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//CreateUser creates new instance of user
func (r *UserRepository) CreateUser(u *models.User) error { //for POST Method
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	db, err := sql.Open("sqlite3", r.store.db)
	if err != nil {
		return err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return err
	}

	result, err := db.Exec("INSERT INTO Users(email, encrypted_password) VALUES($1, $2)", u.Email, u.EncryptedPassword)
	if err != nil {
		return err
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return err
	}
	u.UserID = int32(ID)

	return nil
}

//GetUser ...
func (r *UserRepository) GetUser(email string) (*models.UserDecode, error) {
	u := &models.UserDecode{}
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	db, err := sql.Open("sqlite3", r.store.db)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return nil, err
	}

	if err := db.QueryRow(
		"SELECT * FROM Users WHERE email = $1", email,
	).Scan(
		&u.UserID,
		&u.UserName,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.EncryptedPassword,
		&u.PhoneNumber,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
//krossovki kurtki shtani kofti
//UpdateUser updates user in db
func (r *UserRepository) UpdateUser(user *models.User) (error) {
		checkPassword, err := r.getPassword(user.Email)
		if err != nil {
			return err
		}
		if !checkPassword.ComparePassword(user.UserPassword) {
			return err
		}

		r.store.mu.Lock()
		defer r.store.mu.Unlock()
	
		db, err := sql.Open("sqlite3", r.store.db)
		if err != nil {
			return err
		}
		defer db.Close()


		if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
			return err
		}
		
		if _, err = db.Exec("UPDATE Users SET username = $1, firstname = $2, lastname = $3, phone_number = $4 WHERE email = $5", 
			user.UserName, user.FirstName, user.LastName, user.PhoneNumber, user.Email); err != nil {
				return err
			}
		
		return nil
}

//UpdateUser updates user in db
func (r *UserRepository) DeleteUser(user *models.User) (error) {
	checkPassword, err := r.getPassword(user.Email)
	if err != nil {
		return err
	}
	if !checkPassword.ComparePassword(user.UserPassword) {
		return err
	}

	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	db, err := sql.Open("sqlite3", r.store.db)
	if err != nil {
		return err
	}
	defer db.Close()


	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return err
	}
	
	if _, err = db.Exec("DELETE FROM Users WHERE email = $1", user.Email); err != nil {
			return err
		}
	
	return nil
}

//GetUser ...
func (r *UserRepository) getPassword(email string) (*models.User, error) {
	u := &models.User{}
	r.store.mu.Lock()
	defer r.store.mu.Unlock()

	db, err := sql.Open("sqlite3", r.store.db)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	if _, err := db.Exec("CREATE TABLE IF NOT EXISTS Users(userid INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT, firstname TEXT, lastname TEXT, email TEXT NOT NULL UNIQUE, encrypted_password TEXT NOT NULL, phone_number INTEGER)"); err != nil {
		return nil, err
	}

	if err := db.QueryRow(
		"SELECT encrypted_password FROM Users WHERE email = $1", email,
	).Scan(	
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}