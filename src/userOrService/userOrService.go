package userOrService

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"saltsec/database"
	"saltsec/middleware"
)

// TODO: MILE
type UserOrService struct {
	ID       uint64 `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserOrServiceDTO struct {
	ID               uint64 `json:"id"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	CertType         string `json:"certType"`
	ParentCommonName string `json:"parentCommonName"`
}

func AddUserOrServiceToDB(uos *UserOrService, db *database.DBConn) error {
	return db.DB.Create(uos).Error
}

func UpdateUserOrService(uos *UserOrService, db *database.DBConn) error {
	newUos := UserOrService{}
	if result := db.DB.First(&newUos); result.Error != nil {
		return result.Error
	}

	newUos.Username = uos.Username
	newUos.Password = uos.Password
	return db.DB.Save(&newUos).Error
}

func RemoveUserOrService(id uint64, db *database.DBConn) error {
	return db.DB.Delete(id).Error
}

func GetUserOrService(id uint64, uos *UserOrService, db *database.DBConn) error {
	return db.DB.First(uos).Error
}

func GetAllUserOrServices(uos *[]UserOrService, db *database.DBConn) error {
	return db.DB.Find(uos).Error
}

func GetAll(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		uoss := []UserOrService{}
		if err := GetAllUserOrServices(&uoss, db); err != nil {
			log.Print(err)
		}
		_ = json.NewEncoder(w).Encode(uoss)
	}
}

func AddUosAndCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto UserOrServiceDTO
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		err := dto.loadCertDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Lemara Request: "+err.Error(), http.StatusBadRequest)
			return
		}
		uos := UserOrService{ID: rand.Uint64(), Username: dto.Username, Password: dto.Password}

		_ = AddUserOrServiceToDB(&uos, db)
	}
}

func (dto *UserOrServiceDTO) loadCertDTO(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}

func (a UserOrService) ToString() string {
	return fmt.Sprintf("Admin {ID: %d, Username: %s, Password: %s}",
		a.ID, a.Username, a.Password)
}
