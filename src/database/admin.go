package database

import (
	"saltsec/admin"
)

func (db *DBConn) AddAdmin(a *admin.Admin) error {
	return db.DB.Create(a).Error
}

func (db *DBConn) UpdateAdmin(a *admin.Admin) error {
	newAdmin := admin.Admin{}
	if result := db.DB.First(&newAdmin); result.Error != nil {
		return result.Error
	}

	newAdmin.Username = a.Username
	newAdmin.Email = a.Email
	newAdmin.Password = a.Password
	return db.DB.Save(&newAdmin).Error
}

func (db *DBConn) RemoveAdmin(id uint64) error {
	return db.DB.Delete(id).Error
}

func (db *DBConn) GetAdmin(id uint64, a *admin.Admin) error {
	return db.DB.First(a).Error
}

func (db *DBConn) GetAllAdmins(a *[]admin.Admin) error {
	return db.DB.Find(a).Error
}