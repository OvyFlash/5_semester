package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"neckname/pkg/models"
	"net/http"
	"os"
	"strconv"
)

func (s *server) handleFollow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followed := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(followed); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		follower := &models.User{UserID: userID}

		if err := s.DB.FollowUser(followed, follower); err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleUnfollow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		followed := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(followed); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		follower := &models.User{UserID: userID}

		if err := s.DB.UnfollowUser(followed, follower); err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleAddComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &models.Comment{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		c.UserID = userID
		comment, err := s.DB.AddComment(c)
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}
		s.respond(w, r, http.StatusOK, comment)
	}
}

func (s *server) handleUpdateComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &models.Comment{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if c.CommentID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("This post does not exist"))
			return
		}
		c.UserID = userID
		commentary, err := s.DB.UpdateComment(c)
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, commentary)
	}
}

func (s *server) handleDeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &models.Comment{}
		if err := json.NewDecoder(r.Body).Decode(c); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if c.CommentID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("This post does not exist"))
			return
		}

		c.UserID = userID
		err = s.DB.DeleteComment(c)
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleLike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &models.Like{}
		if err := json.NewDecoder(r.Body).Decode(l); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		l.UserID = userID

		err = s.DB.Like(l)
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, l)
	}
}

func (s *server) handleRemoveLike() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := &models.Like{}
		if err := json.NewDecoder(r.Body).Decode(l); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		l.UserID = userID
		err = s.DB.RemoveLike(l)
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

//handleHome sends user to homepage if cookie is present
func (s *server) handleHome() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		userID, err := s.jsonToken.IsTokenValid(token)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		u, err := s.DB.GetUserByID(userID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		path := getImage(&w, r, u.UserID, true)
		if path != "" {
			u.ProfilePic = path
		}
		s.respond(w, r, http.StatusCreated, u)
	}
}

//handleLogOut sends user to homepage if cookie is present
func (s *server) handleLogOut() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }
		// session.Options.MaxAge = -1
		// session.Save(r, w)
		// token := r.Header.Get("Authorization")
		// userID, err := s.jsonToken.IsTokenValid(token)
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }

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
		if !user.ComparePassword(u) {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("Email or password does not match"))
			return
		}
		//delete password
		user.ClearPassword()
		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// }
		// session.Values["userid"] = user.UserID
		// session.Save(r, w)
		newToken := s.jsonToken.GetToken(user.UserID)
		user.JSONToken = newToken
		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) handleStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		posts, err := s.DB.GetAllPosts()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, posts)

	}
}

func (s *server) handleSendError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.error(w, r, http.StatusNotFound, errors.New("404 page not found haha"))
	}
}
func (s *server) handlePlay() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//saveImage(r)
		buffer := new(bytes.Buffer)
		f, err := os.Open("/home/oleg/Рабочий стол/KPI/5 Семестр/Системи обробки сигналів/NECKNAME SocialNetwork/5_semester-master/2.jpeg")
		if err != nil {
			fmt.Println("Unable to open file")
			return
		}
		defer f.Close()
		img, _, err := image.Decode(f)
		if err := jpeg.Encode(buffer, img, &jpeg.Options{100}); err != nil {
			fmt.Println("unable to encode image.")
		}

		w.Header().Set("Content-Type", "image/jpeg")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

		if _, err := w.Write(buffer.Bytes()); err != nil {
			fmt.Println("unable to write image.")
		}

	}
}

/*

	buffer := new(bytes.Buffer)
	f, err := os.Open("/home/oleg/Рабочий стол/KPI/5 Семестр/Системи обробки сигналів/NECKNAME SocialNetwork/5_semester-master/Снимок экрана от 2020-08-15 00-15-43.png")
	if err != nil {
		return
	}
	img, _, err := image.Decode(f)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

*/
