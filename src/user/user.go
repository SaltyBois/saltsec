package user

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saltsec/database"
	"saltsec/middleware"

	"github.com/gorilla/mux"
)

// TODO: MILE
type User struct {
	Username string   `gorm:"primaryKey" json:"username"`
	Password string   `json:"password"`
	Salt     string   `json:"-"`
	Certs    []string `json:"certs"`
}

type UserDTO struct {
	Username         string `json:"username"`
	Password         string `json:"password"`
	CertType         string `json:"certType"`
	ParentCommonName string `json:"parentCommonName"`
}

func AddUserDB(uos *User, db *database.DBConn) error {
	return db.DB.Create(uos).Error
}

func UpdateUserDB(uos *User, db *database.DBConn) error {
	newUos := User{}
	if result := db.DB.First(&newUos); result.Error != nil {
		return result.Error
	}

	newUos.Username = uos.Username
	newUos.Password = uos.Password
	return db.DB.Save(&newUos).Error
}

func RemoveUserDB(id uint64, db *database.DBConn) error {
	return db.DB.Delete(id).Error
}

func GetUserDB(username string, uos *User, db *database.DBConn) error {
	return db.DB.First(uos, username).Error
}

func GetAllUsersDB(uos *[]User, db *database.DBConn) error {
	return db.DB.Find(uos).Error
}

func GetUserByUsernameDB(username string, uos *User, db *database.DBConn) error {
	return db.DB.Where("username = ?", username).First(uos).Error
}

func GetAllUsers(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uoss := []User{}
		if err := GetAllUsersDB(&uoss, db); err != nil {
			log.Print(err)
		}
		_ = json.NewEncoder(w).Encode(uoss)
	}
}

func GetUser(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto UserDTO
		params := mux.Vars(r)
		username := params["username"]
		dto.Username = username
		uos := User{}
		if err := GetUserByUsernameDB(dto.Username, &uos, db); err != nil {
			log.Print(err)
		}
		json.NewEncoder(w).Encode(uos)
	}
}

// TODO(Jovan): Split into separate?
func AddUserAndCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto UserDTO
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := dto.loadCertDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Lemara Request: "+err.Error(), http.StatusBadRequest)
			return
		}
		uos := User{Username: dto.Username, Password: dto.Password}

		_ = AddUserDB(&uos, db)
	}
}

func (dto *UserDTO) loadCertDTO(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}

func (a User) ToString() string {
	return fmt.Sprintf("Admin {Username: %s, Password: %s}", a.Username, a.Password)
}
