package apiserver

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
)

type userID int32

const userWithID userID = iota + 1

var findImage = regexp.MustCompile(`^/(profile_pic_[0-9]+)\.jpeg$`)

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		enableCors(&w)

		fmt.Printf("Remote addr: %s\n", r.RemoteAddr)
		fmt.Printf("Started %s %s\n", r.Method, r.RequestURI)

		token := r.Header.Get("Authorization")
		unauthorized := len(token) < 1

		var (
			err    error
			userID int32 = -1
		)
		if len(token) > 1 {
			userID, err = s.jsonToken.IsTokenValid(token)
			if err != nil {
				fmt.Println(err)
			}
		}
		if userID < 0 {
			unauthorized = true
		}

		w.Header().Set("UserID", int32toString(userID))
		r = r.WithContext(context.WithValue(r.Context(), userWithID, int32toString(userID)))

		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// }
		switch uri := r.RequestURI; {
		case string(findImage.Find([]byte(uri))[:]) != "":
			s.sendImage(w, r)
			return
		//start
		//home
		//login
		case uri == "/getMe" && r.Method == "GET":
			// if session.Values["userid"] == nil {
			// 	//	s.error(w, r, http.StatusUnauthorized, errors.New("User unautorizhed"))
			// 	http.Error(w, errors.New("User unauthorized").Error(), http.StatusUnauthorized)
			// 	return
			// }
			// next.ServeHTTP(w, r)
			// return
			if unauthorized {
				//	s.error(w, r, http.StatusUnauthorized, errors.New("User unautorizhed"))
				http.Error(w, errors.New("User unauthorized").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/login" && r.Method == "POST":
			if !unauthorized {
				http.Error(w, "User authorized", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		}
		next.ServeHTTP(w, r)

	})
}

func (s *server) userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// }
		// token := r.Header.Get("Authorization")
		// unauthorized := len(token) < 1

		var (
			userID       int64
			err          error
			unauthorized bool
		)
		userID, err = strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userID < 0 {
			unauthorized = true
		}
		switch uri := r.RequestURI; {
		//redirect to home if user is authentificated
		case uri == "/user/create" && r.Method == "POST":
			// if session.Values["userid"] != nil {
			// 	http.Error(w, errors.New("User authorized").Error(), http.StatusNotAcceptable)
			// 	return
			// }
			// next.ServeHTTP(w, r)
			if !unauthorized {
				http.Error(w, errors.New("User authorized").Error(), http.StatusNotAcceptable)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/update" && r.Method == "PUT":
			if unauthorized {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/delete" && r.Method == "DELETE":
			if unauthorized {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/get" && r.Method == "GET":
			next.ServeHTTP(w, r)
			return

		case uri == "/user/follow" && r.Method == "POST":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/unfollow" && r.Method == "POST":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/comment" && r.Method == "POST":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/comment" && r.Method == "PUT":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/comment" && r.Method == "DELETE":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/like" && r.Method == "POST":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/user/like" && r.Method == "DELETE":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		}
	})
}

func (s *server) routeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			userID       int64
			err          error
			unauthorized bool
		)
		userID, err = strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userID < 0 {
			unauthorized = true
		}
		switch uri := r.RequestURI; {
		//redirect to home if user is authentificated
		case uri == "/route/create" && r.Method == "POST":
			if unauthorized {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/route/getMy" && r.Method == "GET":
			if unauthorized {
				http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/route/get" && r.Method == "POST":
			// if unauthorized {
			// 	http.Error(w, errors.New("User unautorizhed").Error(), http.StatusUnauthorized)
			// 	return
			// }
			next.ServeHTTP(w, r)
			return
		}

	})
}

func (s *server) postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			userID       int64
			err          error
			unauthorized bool
		)
		userID, err = strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if userID < 0 {
			unauthorized = true
		}
		switch uri := r.RequestURI; {
		case uri == "/post/create" && r.Method == "POST":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/post/update" && r.Method == "PUT":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return

		case uri == "/post/delete" && r.Method == "DELETE":
			if unauthorized {
				http.Error(w, "User unautorizhed", http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
			return
		case uri == "/post/get" && r.Method == "POST":
			next.ServeHTTP(w, r)
			return
		}

	})
}
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func int32toString(n int32) string {
	buf := [11]byte{}
	pos := len(buf)
	i := int64(n)
	signed := i < 0
	if signed {
		i = -i
	}
	for {
		pos--
		buf[pos], i = '0'+byte(i%10), i/10
		if i == 0 {
			if signed {
				pos--
				buf[pos] = '-'
			}
			return string(buf[pos:])
		}
	}
}
