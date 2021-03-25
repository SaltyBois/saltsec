package userOrService

import (
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"saltsec/cert"
	"saltsec/database"
)
 // TODO: MILE
type UserOrService struct {
	ID       		 uint64		`json:"id"`
	Username 		 string		`json:"username"`
	Password 		 string		`json:"password"`
	CertSerialNumber *big.Int	`json:"certSerialNumber"`
	cert.CertType
}

func AddUserOrService(uos *UserOrService, db *database.DBConn) error {
	//if (uos.CertType == cert.Root) {
	//
	//}
	//else if (uos.CertType == cert.Root) {
	//
	//}
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

func (a UserOrService) ToString() string {
	return fmt.Sprintf("Admin {ID: %d, Username: %s, Password: %s}",
		a.ID, a.Username, a.Password)
}
