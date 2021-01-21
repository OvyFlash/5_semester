package apiserver

import (
	"encoding/json"
	"fmt"
	"neckname/internal/config"
	"neckname/internal/jwt"
	"neckname/pkg/sqlitedatabase"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

type server struct {
	router       *mux.Router
	DB           *sqlitedatabase.Sqldatabase
	sessionStore *sessions.CookieStore
	jsonToken    *jwt.JWT
}

//StartAPIServer ...
func StartAPIServer(config *config.Config) error {

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	s := &server{
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
	}
	s.DB = sqlitedatabase.MakeSqldatabase(config)
	s.jsonToken = jwt.NewJWT()
	s.configureRouter()

	headersOk := handlers.AllowedHeaders([]string{"Content-Type", "application/json", "Origin", "Cache-Control", "X-App-Token", "X-Requested-With", "Access-Control-Allow-Origin", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"})
	return http.ListenAndServe(config.BindAddr, handlers.CORS(originsOk, headersOk, methodsOk)(s))
}

func (s *server) configureRouter() {
	s.router.Use(s.logRequest)
	s.router.HandleFunc("/{any}", s.handleSendError())
	s.router.HandleFunc("/play", s.handlePlay()).Methods("POST", "GET").Name("play")
	//start was "/"
	s.router.HandleFunc("/start", s.handleStart()).Methods("GET").Name("start")  //sends posts
	s.router.HandleFunc("/login", s.handleLogIn()).Methods("POST").Name("login") //auth response
	//getMe was "/home"
	s.router.HandleFunc("/getMe", s.handleHome()).Methods("GET").Name("getMe") //profilepage

	
	//post
	postRouter := s.router.PathPrefix("/post").Subrouter()
	postRouter.Use(s.postMiddleware)
	//it was user/post
	postRouter.HandleFunc("/create", s.handleCreatePost()).Methods("POST").Name("createpost")
	postRouter.HandleFunc("/update", s.handleUpdatePost()).Methods("PUT").Name("updatepost")
	postRouter.HandleFunc("/delete", s.handleDeletePost()).Methods("DELETE").Name("deletepost")
	//new method
	postRouter.HandleFunc("/get", s.handleGetPost()).Methods("POST").Name("getpost")
	//add method get concretee post
	//User
	userRouter := s.router.PathPrefix("/user").Subrouter()
	userRouter.Use(s.userMiddleware)
	//get all posts
	userRouter.HandleFunc("/create", s.handleCreateUser()).Methods("POST").Name("createuser")   //auth response
	userRouter.HandleFunc("/update", s.handleUpdateUser()).Methods("PUT").Name("updateuser")    //updates user
	userRouter.HandleFunc("/delete", s.handleDeleteUser()).Methods("DELETE").Name("deleteuser") //auth response (invalid)
	userRouter.HandleFunc("/get", s.handleGetUser()).Methods("GET").Name("getuser")             //gets concretee user

	routeRouter := s.router.PathPrefix("/route").Subrouter()
	routeRouter.Use(s.routeMiddleware) //POMENIAL OUTPUT
	//it was user/route
	routeRouter.HandleFunc("/create", s.handleCreateRoute()).Methods("POST").Name("route") //create route
	//it was user/getroute
	routeRouter.HandleFunc("/get", s.handleGetRoute()).Methods("POST").Name("getroute") //get concretee route
	//it was user/routes
	routeRouter.HandleFunc("/getMy", s.handleGetRoutes()).Methods("GET").Name("routes") //get all user routes

	/*
		ADD PHOTO TO USER AND AND PHOTO TO POST
	*/
	//follow
	//unfollow
	userRouter.HandleFunc("/follow", s.handleFollow()).Methods("POST").Name("follow")
	userRouter.HandleFunc("/unfollow", s.handleUnfollow()).Methods("DELETE").Name("unfollow")

	//post/comment
	userRouter.HandleFunc("/comment", s.handleAddComment()).Methods("POST").Name("createcomment")
	userRouter.HandleFunc("/comment", s.handleUpdateComment()).Methods("UPDATE").Name("updatecomment")
	userRouter.HandleFunc("/comment", s.handleDeleteComment()).Methods("DELETE").Name("deletecomment")
	//post/like
	userRouter.HandleFunc("/like", s.handleLike()).Methods("POST").Name("like")
	userRouter.HandleFunc("/like", s.handleRemoveLike()).Methods("DELETE").Name("rmlike")

}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"Error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) { //For default route
	s.router.ServeHTTP(w, r)
}
