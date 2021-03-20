package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"saltsec/globals"
)

// NOTE(Jovan): Required for new method definition of non-local
// type gorm.DB
type DBConn struct {
	DB *gorm.DB
}

func ConnectToDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		globals.HOST_DB, globals.USER_DB, globals.PASSWORD_DB, globals.NAME_DB, globals.PORT_DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}
