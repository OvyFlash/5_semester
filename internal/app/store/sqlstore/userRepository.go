package sqlstore

import (
	"database/sql"
	"neckname/internal/app/models"
	"neckname/internal/app/store"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//CreateUser creates new instance of user
func (r *UserRepository) CreateUser(u *models.User) error {
	if err := u.Validate(); err != nil {
		return err
	}

	if err := u.BeforeCreate(); err != nil {
		return err
	}

	result, err := r.store.db.Exec("INSERT INTO Users(email, encrypted_password) VALUES(?, ?)",
		u.Email, u.EncryptedPassword) //add basic values (email password)
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

//FindByEmail ...
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT userid, email, encrypted_password FROM Users WHERE email = ?",
		email,
	).Scan(
		&u.UserID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

//FindByID ...
func (r *UserRepository) FindByID(id int) (*models.User, error) {
	u := &models.User{}
	if err := r.store.db.QueryRow(
		"SELECT userid, email, encrypted_password FROM Users WHERE id = ?",
		id,
	).Scan(
		&u.UserID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}