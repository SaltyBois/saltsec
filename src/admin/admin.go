package admin

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Admin struct {
	ID       uint64 `json: id`
	Username string `json: username`
	Email    string `json: email`
	// TODO(Jovan): Hash
	Password string `json: password`
}

func (db *gorm.DB) AddAdmin(admin *Admin) error {
	return db.Create(admin).Error
}

func (db *gorm.DB) UpdateAdmin(admin *Admin) error {
	newAdmin := Admin{}
	if result := db.First(&newAdmin); result.Error != nil {
		return result.Erorr
	}

	newAdmin.Username = admin.Username
	newAdmin.Email = admin.Email
	newAdmin.Password = admin.Password
	return db.Save(&newAdmin)
}

func (db *gorm.DB) RemoveAdmin(id uint64) error {
	return db.Delete(id).Error
}

func (db *gorm.DB) GetAdmin(id uint64, admin *Admin) error {
	return db.First(admin).Error
}

func (db *gorm.DB) GetAllAdmins(admins *[]Admin) error {
	return db.Find(admins).Error
}

func (a Admin) ToString() string {
	return fmt.Sprintf("Admin {ID: %s, Username: %s, Email: %s, Password: %s}",
		a.ID, a.Username, a.Email, a.Password)
}
