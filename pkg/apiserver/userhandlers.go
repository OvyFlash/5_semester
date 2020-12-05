package apiserver

import (
	"encoding/json"
	"fmt"
	"neckname/pkg/models"
	"net/http"
)

func (s *server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		fmt.Println(u)
		if err := u.CheckEmailAndPassword(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		//add user to database if not exists
		u, err := s.DB.CreateUser(u)
		if err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		//delete password
		u.ClearPassword()
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
		}
		session.Values["userid"] = u.UserID
		session.Save(r, w)
		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		fmt.Println(u)

		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		user, err := s.DB.GetUserByID(session.Values["userid"].(int32))
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}
		//compare user with new user and update fields
		u.Difference(user)

		if err := s.DB.UpdateUser(u); err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
		}
		
		u.ClearPassword()
		s.respond(w, r, http.StatusOK, u)
	}
}

func (s *server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if err := s.DB.DeleteUser(session.Values["userid"].(int32)); err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
		}

		session.Options.MaxAge = -1
		session.Save(r, w)
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleCreateRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		u, err := s.DB.GetUserByID(session.Values["userid"].(int32))
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}

		route, err := s.DB.CreateRoute(u)
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}
		session.Values["routeid"] = route.RouteID
		session.Save(r, w)
		s.respond(w, r, http.StatusOK, route)
	}
}


// func (s *server) handleDeleteRoute() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		session, err := s.sessionStore.Get(r, "session")
// 		if err != nil {
// 			s.error(w, r, http.StatusInternalServerError, err)
// 			return
// 		}

// 		route, err := s.DB.GetRouteByID(session.Values["routeid"].(int64))
// 		if err != nil {
// 			s.error(w,r, http.StatusInternalServerError, err)
// 			return
// 		}

// 		if err := s.DB.DeleteRoute(u); err != nil {
// 			s.error(w,r, http.StatusInternalServerError, err)
// 			return
// 		}
// 		session.Values["routeid"] = route.RouteID
// 		session.Save(r, w)
// 		s.respond(w, r, http.StatusOK, route)
// 	}
// }

func (s *server) handleAddPoint() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		point := &models.Point{}
		if err := json.NewDecoder(r.Body).Decode(point); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		route, err := s.DB.GetRouteByID(session.Values["routeid"].(int64))
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}
		point.RouteID = route.RouteID

		//maybe add some logic to decrease server
		if err := s.DB.AddPoint(point); err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleFinishRoute() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		point := &models.Point{}
		if err := json.NewDecoder(r.Body).Decode(point); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		session, err := s.sessionStore.Get(r, "session")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		route, err := s.DB.GetRouteByID(session.Values["routeid"].(int64))
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}
		session.Values["routeid"] = nil
		session.Save(r, w)

		point.RouteID = route.RouteID
		//maybe add some logic to decrease server
		routestat, err := s.DB.FinishRoute(point, session.Values["userid"].(int32))
		if err != nil {
			s.error(w,r, http.StatusInternalServerError, err)
			return
		}
		
		s.respond(w, r, http.StatusOK, routestat)
	}
}