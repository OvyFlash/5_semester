package apiserver

import (
	"encoding/json"
	"errors"
	"neckname/pkg/models"
	"net/http"
	"strconv"
)

func (s *server) handleCreateRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		route := &models.Route{}
		if err := json.NewDecoder(r.Body).Decode(route); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if len(route.Points) < 1 {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("No points"))
			return
		}
		if route.Start == 0 || route.Finish == 0 {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("Start or finish missing"))
			return
		}
		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := s.DB.GetUserByID(int32(userID))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		newRoute, err := s.DB.CreateRoute(u, route)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, newRoute)
	}
}

func (s *server) handleGetRoutes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		routes, err := s.DB.GetUserRoutes(int32(userID))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, routes)
	}
}

func (s *server) handleGetRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		route := &models.Route{}
		if err := json.NewDecoder(r.Body).Decode(route); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if route.RouteID == 0 {
			s.error(w, r, http.StatusUnprocessableEntity, errors.New("No routeid entered"))
		}

		newRoute, err := s.DB.GetRouteByID(route.RouteID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, newRoute)
	}
}
