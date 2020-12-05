package apiserver

import (
	"neckname/internal/app/store/sqlstore"
	"net/http"

	"github.com/gorilla/sessions"
)

//StartAPIServer ...
func StartAPIServer(config *Config) error {
	// db, err := newDB(config.DatabaseULR)
	// if err != nil {
	// 	return err
	// }
	// defer db.Close()

	store := sqlstore.NewStore(config.DatabaseURL)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	
	srv := newServer(store, sessionStore)

	return http.ListenAndServe(config.BindAddr, srv)
}

// func newDB(databaseURL string) (*sql.DB, error) {
// 	db, err := sql.Open("mysql", databaseURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	if err := db.Ping(); err != nil {
// 		return nil, err
// 	}

// 	return db, nil
// }
