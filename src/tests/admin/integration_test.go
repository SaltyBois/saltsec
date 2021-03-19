package tests

import (
    "testing"
    "net/http"
    "saltsec/admin"
    "saltsec/database"
    "gorm.io/gorm"
)

type Application struct {
    db database.DBConn
    handler http.handler
}

func seedAdmins(db *database.DBConn) (*gorm.DB, error) {
    admins := []admin.Admin{
        admin.Admin{
            ID: 123,
            Username: "AdminOneTwoThree",
            Email: "adminonetwothree@email.com",
            Password: "AdminOneTwoThree",
        },
    }
    db.DB.AutoMigrate(admin.Admin{})
    tx := db.DB.Begin()
    for _, a := range admins {
        if err := db.DB.Create(&a).Error; err != nil {
            return tx, err
        }
    }
    return tx, nil
}

func TestAdminGet(t *testing.T) {
    a := Application{}
    a.db.ConnectToDb()
	a.handler := cors.Default().Handler(router)


    tx, err := seedAdmins()
}
