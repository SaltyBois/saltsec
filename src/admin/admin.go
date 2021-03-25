package admin

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"saltsec/database"
)

type Admin struct {
	ID       uint64 `json: id`
	Username string `json: username`
	Email    string `json: email`
	// TODO(Jovan): Hash
	Password string `json: password`
}

func AddAdmin(admin *Admin, db *database.DBConn) error {
	return db.DB.Create(admin).Error
}

func UpdateAdmin(admin *Admin, db *database.DBConn) error {
	newAdmin := Admin{}
	if result := db.DB.First(&newAdmin); result.Error != nil {
		return result.Error
	}

	newAdmin.Username = admin.Username
	newAdmin.Email = admin.Email
	newAdmin.Password = admin.Password
	return db.DB.Save(&newAdmin).Error
}

func RemoveAdmin(id uint64, db *database.DBConn) error {
	return db.DB.Delete(id).Error
}

func GetAdmin(id uint64, admin *Admin, db *database.DBConn) error {
	return db.DB.First(admin).Error
}

func GetAllAdmins(admin *[]Admin, db *database.DBConn) error {
	return db.DB.Find(admin).Error
}

func GetAll(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		admins := []Admin{}
		if err := GetAllAdmins(&admins, db); err != nil {
			log.Print(err)
		}
		json.NewEncoder(w).Encode(admins)
	}
}

func GetFirstAdmin(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		admin := Admin{}
		if err := GetAdmin(1, &admin, db); err != nil {
			log.Print(err)
		}
		json.NewEncoder(w).Encode(admin)
	}
}

func (a Admin) ToString() string {
	return fmt.Sprintf("Admin {ID: %d, Username: %s, Email: %s, Password: %s}",
		a.ID, a.Username, a.Email, a.Password)
}
