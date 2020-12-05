package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"neckname/internal/config"
	"neckname/pkg/models"
	"neckname/pkg/sqlitedatabase"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

//StartAPIServer ...
func StartAPIServer(config *config.Config) error {

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := &server{
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
	}
	s.DB = sqlitedatabase.MakeSqldatabase(config)
	s.configureRouter()
	return http.ListenAndServe(config.BindAddr, s)
}

type server struct {
	router       *mux.Router
	DB           *sqlitedatabase.Sqldatabase
	sessionStore *sessions.CookieStore
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}))) //s liubih domenov
	s.router.Use(s.logRequest)

	//logout deletes cookies
	s.router.HandleFunc("/logout", s.handleLogOut()).Methods("GET").Name("logout")
	//login creates cookie and allows to enter if passwords match
	s.router.HandleFunc("/login", s.handleLogIn()).Methods("POST").Name("login")
	//send information about authentificated user
	s.router.HandleFunc("/home", s.handleHome()).Methods("GET", "POST").Name("home") //??
	// s.router.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello WORLD"))
	// })
	// s.router.HandleFunc("/follow", s.handleFollow()).Methods("POST").Name("follow")
	// s.router.HandleFunc("/unfollow", s.handleUnfollow()).Methods("DELETE").Name("unfollow")


	//User
	userRoute := s.router.PathPrefix("/user").Subrouter()
	userRoute.Use(s.userMiddleware)
	userRoute.HandleFunc("/create", s.handleCreateUser()).Methods("POST").Name("createuser")
	userRoute.HandleFunc("/update", s.handleUpdateUser()).Methods("PUT").Name("updateuser")
	userRoute.HandleFunc("/delete", s.handleDeleteUser()).Methods("DELETE").Name("deleteuser")
	
	userRoute.HandleFunc("/createroute", s.handleCreateRoute()).Methods("POST").Name("createroute")
	//userRoute.HandleFunc("/deleteroute", s.handleDeleteRoute()).Methods("DELETE").Name("deleteroute")
	
	userRoute.HandleFunc("/addpoint", s.handleAddPoint()).Methods("POST").Name("addpoint")
	userRoute.HandleFunc("/finishroute", s.handleFinishRoute()).Methods("GET").Name("finishroute")

	/*
	//follow sbd
	//unfollow sbd
	
	//createroute(should return route with new routeid) (mb cookie)
	//deleteroute
	//addpoint(with cookie add to routeid)
	//finish route(create route stat and delete cookie)

	//add comment
	//delete comment
	//update comment
	
	//create post
	//delete post
	//update post
	
	//give like
	//remove like
	*/
	// s.router.HandleFunc("/getuser", s.handleUsersGet()).Methods("GET")

	//private := s.router.PathPrefix("/private").Subrouter()
	//private.Use(s.authenticateUser)
	//private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")

}

//handleHome sends user to homepage if cookie is present
func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u, err := s.DB.GetUserByID(session.Values["userid"].(int32))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusCreated, u)
	}
}

//handleLogOut sends user to homepage if cookie is present
func (s *server) handleLogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Options.MaxAge = -1
		session.Save(r, w)

		s.respond(w, r, http.StatusOK, nil)
	}
}

//handleLogOut sends user to homepage if cookie is present
func (s *server) handleLogIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		fmt.Println(u)

		//get user by email
		user, err := s.DB.GetUserByEmail(u.Email)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		//compare passwords
		if user.ComparePassword(u) {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("Email or password does not match"))
			return
		}
		//delete password
		user.ClearPassword()
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		session.Values["userid"] = user.UserID
		session.Save(r, w)
		s.respond(w, r, http.StatusOK, user)
	}
}

// //handleFollow sends user to homepage if cookie is present
// func (s *server) handleFollow() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, err := s.sessionStore.Get(r, "session")
// 		if err != nil {
// 			s.error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}
// 		u, err := s.DB.GetUserByID(session.Values["userid"].(int32))
// 		if err != nil {
// 			s.error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}


// 		s.respond(w, r, http.StatusOK, nil)
// 	}
// }

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"Error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) { //For default route
	s.router.ServeHTTP(w, r)
}
