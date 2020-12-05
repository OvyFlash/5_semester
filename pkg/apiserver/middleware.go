package apiserver

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Remote addr: %s\n", r.RemoteAddr)
		fmt.Printf("Started %s %s\n", r.Method, r.RequestURI)

		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		switch uri := r.RequestURI; {
		case uri == "/home" && (r.Method == "GET" || r.Method == "POST"):
			if session.Values["userid"] == nil {
				//	s.error(w, r, http.StatusUnauthorized, errors.New("User unautorizhed"))
				http.Error(w, errors.New("User unauthorized").Error(), http.StatusUnauthorized)
				return
			}

		case uri == "/login" && r.Method == "POST":
			if session.Values["userid"] != nil {
				http.Error(w, "User authorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/logout" && r.Method == "GET":
			if session.Values["userid"] == nil {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		case uri == "/follow" && r.Method == "POST":
			if session.Values["userid"] == nil {
				http.Error(w, "User unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w,r)
			return
		case uri == "/unfollow" && r.Method == "DELETE":
			if session.Values["userid"] == nil {
				http.Error(w, "User unauthorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w,r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *server) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}

		switch uri := r.RequestURI; {
		//redirect to home if user is authentificated
		case uri == "/user/create" && r.Method == "POST":
			if session.Values["userid"] != nil {
				fmt.Println("Redirecting to home")
				url, err := mux.CurrentRoute(r).Subrouter().Get("home").URL()
				if err != nil {
					panic(err)
				}
				http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
				return
			}
			next.ServeHTTP(w, r)
		case uri == "/user/update" && r.Method == "PUT":
			if session.Values["userid"] == nil {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		case uri == "/user/delete" && r.Method == "DELETE":
			if session.Values["userid"] == nil {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)


		case uri == "/user/createroute" && r.Method == "POST":
			//add check for route cookie
			if session.Values["userid"] == nil {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			if session.Values["routeid"] != nil {
				url, err := mux.CurrentRoute(r).Subrouter().Get("addpoint").URL()
				if err != nil {
					panic(err)
				}
				http.Redirect(w, r, url.String(), http.StatusPermanentRedirect)
				return
			}
			next.ServeHTTP(w, r)
		// case uri == "/user/deleteroute" && r.Method == "DELETE":
		// 	if session.Values["userid"] == nil {
		// 		http.Error(w, "User unautorizhed", http.StatusUnauthorized)
		// 		return
		// 	}
		// 	if session.Values["routeid"] == nil {
		// 		http.Error(w, "User has no route started", http.StatusUnauthorized)
		// 		return
		// 	}
		// 	next.ServeHTTP(w, r)
		// default:
		// 	next.ServeHTTP(w, r)
		case uri == "/user/addpoing" && r.Method == "POST":
			if session.Values["userid"] == nil {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			if session.Values["routeid"] == nil {
				http.Error(w, "User has no route started", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		
		case uri == "/user/finishroute" && r.Method == "GET":
			if session.Values["userid"] == nil {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			if session.Values["routeid"] == nil {
				http.Error(w, "User has no route started", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		}
	})
}
