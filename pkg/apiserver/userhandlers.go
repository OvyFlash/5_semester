package apiserver

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"mime/multipart"
	"neckname/pkg/models"
	"net/http"
	"os"
	"strconv"
)

func (s *server) handleCreateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, formFileError := r.FormFile("Image")
		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			if r.FormValue("JSON") != "" {
				if err := json.Unmarshal([]byte(r.FormValue("JSON")), &u); err != nil {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
			} else {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}

		fmt.Println(u)
		if err := u.CheckEmailAndPassword(); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		//add user to database if not exists
		u, err := s.DB.CreateUser(u)
		if err != nil {
			s.error(w, r, http.StatusConflict, err)
			return
		}
		//check if image present
		if formFileError == nil {
			path := saveImage(file, &w, r, u.UserID)
			if path != "" {
				u.ProfilePic = path
			}
			//w.Header().Set("Content-Length", strconv.Itoa(len([]byte(u))))
		}

		//delete password
		u.ClearPassword()

		s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) handleUpdateUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, _, formFileError := r.FormFile("Image")

		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			if r.FormValue("JSON") != "" && formFileError != nil {
				if err := json.Unmarshal([]byte(r.FormValue("JSON")), &u); err != nil {
					s.error(w, r, http.StatusBadRequest, err)
					return
				}
			} else if formFileError != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
		}

		fmt.Println(u)

		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := s.DB.GetUserByID(int32(userID))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		if u.Email != user.Email && u.Email != "" { //means that user wants to update email
			if err := s.DB.CheckUserEmail(u.Email); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				return
			}
		}

		//compare user with new user and update fields

		user.Difference(u)
		
		if err := user.CheckEmailAndPassword(); err != nil && u.UserPassword != "" && u.Email != ""{
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := s.DB.UpdateUser(user); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		if formFileError == nil {
			path := saveImage(file, &w, r, u.UserID)
			if path != "" {
				user.ProfilePic = path
			}
		}

		user.ClearPassword()
		s.respond(w, r, http.StatusOK, user)
	}
}

func (s *server) handleDeleteUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// session, err := s.sessionStore.Get(r, "session")
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }

		userID, err := strconv.ParseInt(r.Context().Value(userWithID).(string), 10, 32)
		if err != nil {
			http.Error(w, errors.New("User authorized").Error(), http.StatusInternalServerError)
			return
		}

		if err := s.DB.DeleteUser(int32(userID)); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := &models.User{}
		if err := json.NewDecoder(r.Body).Decode(u); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if u.UserID == 0 {
			s.error(w, r, http.StatusInternalServerError, errors.New("Havent received id of user"))
			return
		}
		user, err := s.DB.GetUserByID(u.UserID)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		path := getImage(&w, r, user.UserID, true)
		if path != "" {
			user.ProfilePic = path
		}
		user.ClearPassword()
		s.respond(w, r, http.StatusOK, user)
	}
}
const pathToImages = `internal/photos/`

func saveImage(file multipart.File, w *http.ResponseWriter, r *http.Request, userID int32) string {
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if _, err := os.Stat(pathToImages); os.IsNotExist(err) {
		os.Mkdir(pathToImages, os.ModePerm)
	}
	path := fmt.Sprintf("profile_pic_%d.jpeg", userID)

	f, err := os.Create(pathToImages+path)
	if err != nil {
		fmt.Println(err)
	}
	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	f.Close()

	buffer := new(bytes.Buffer)

	
	f, err = os.Open(path)
	if err != nil {
		return ""
	}
	img, _, err = image.Decode(f)
	if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100}); err != nil {
		fmt.Println("unable to encode image.")
	}
	f.Close()

	//(*w).Header().Set("Content-Type", "image/jpeg")
	//(*w).Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	return path
}

func getImage(w *http.ResponseWriter, r *http.Request, userID int32, profilePic bool) string {
	var path string
	if profilePic {
		path = fmt.Sprintf("profile_pic_%d.jpeg", userID)
	} else {
		path = fmt.Sprintf("post_pic_%d.jpeg", userID)
	}
	if _, err := os.Stat(pathToImages+path); os.IsNotExist(err) {
		return ""
	}

	// buffer := new(bytes.Buffer)
	// f, err := os.Open(pathToImages+path)
	// if err != nil {
	// 	return ""
	// }
	// img, _, err := image.Decode(f)
	// if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100}); err != nil {
	// 	fmt.Println("unable to encode image.")
	// }
	// f.Close()
	// (*w).Header().Set("Content-Type", "image/jpeg")
	// (*w).Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	// (*w).
	// if _, err := (*w).Write(buffer.Bytes()); err != nil {
	// 	log.Println("unable to write image.")
	// }
	return path
}
func (s *server) sendImage(w http.ResponseWriter, r *http.Request) {
	var path string
	path = string(r.RequestURI[1:])
	buffer := new(bytes.Buffer)
	if _, err := os.Stat(pathToImages); os.IsNotExist(err) {
		os.Mkdir(pathToImages, os.ModePerm)
	}
	f, err := os.Open(pathToImages+path)
	if err != nil {
		return
	}
	img, _, err := image.Decode(f)
	if err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100}); err != nil {
		fmt.Println("unable to encode image.")
	}
	f.Close()
	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
	}
}
