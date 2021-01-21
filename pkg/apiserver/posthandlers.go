package apiserver

import (
	"database/sql"
	"encoding/json"
	"errors"
	"neckname/pkg/models"
	"net/http"
	"strconv"
)

func (s *server) handleCreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Post{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		p.Route, err = s.DB.GetRouteByID(p.RouteID)
		if err != nil {
			if err == sql.ErrNoRows {
				s.error(w, r, http.StatusBadRequest, errors.New("Route with such id does not exists "))
				return
			}
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if int32(userID) != p.Route.UserID {
			s.error(w, r, http.StatusForbidden, errors.New("You cannot add this route"))
			return
		}
		post, err := s.DB.CreatePost(p, int32(userID))
		if err != nil {
			s.error(w, r, http.StatusOK, err)
			return
		}
		s.respond(w, r, http.StatusOK, post)
	}
}

func (s *server) handleUpdatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Post{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if p.PostID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("This post does not exist"))
			return
		}

		post, err := s.DB.UpdatePost(p, int32(userID))
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, post)
	}
}

func (s *server) handleDeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Post{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if p.PostID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("This post does not exist"))
			return
		}

		err = s.DB.DeletePost(p, int32(userID))
		if err != nil {
			s.error(w, r, http.StatusForbidden, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		p := &models.Post{}
		if err := json.NewDecoder(r.Body).Decode(p); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if p.PostID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("This post does not exist"))
			return
		}

		post, err := s.DB.GetPost(p)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, post)
	}
}
