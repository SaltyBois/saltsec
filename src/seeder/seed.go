package seeder

import (
	"log"
	"saltsec/admin"
	"saltsec/cert"
	"saltsec/database"
	"saltsec/user"
)

type Seed struct {
	Name string
	Run  func(*database.DBConn) error
}

func MigrateData(db *database.DBConn) {
	db.DB.AutoMigrate(&admin.Admin{})
	db.DB.AutoMigrate(&cert.ArchivedCert{})
	db.DB.AutoMigrate(&user.User{})
}

func SeedData(db *database.DBConn) {
	for _, seed := range allSeeds() {
		if err := seed.Run(db); err != nil {
			log.Printf("Seed: '%s' failed with error: '%s'", seed.Name, err)
		}
	}
}

func allSeeds() []Seed {
	return []Seed{
		Seed{
			Name: "CreateAdmin1",
			Run: func(db *database.DBConn) error {
				a := admin.Admin{ID: 1, Username: "admin1", Email: "admin@email.com", Password: "admin1"}
				return admin.AddAdmin(&a, db)
			},
		},
		//Seed{
		//	Name: "User1",
		//	Run: func(db *database.DBConn) error {
		//		uos := userOrService.UserOrService{ID: 1, Username: "user1", Password: "user1"}
		//		return userOrService.AddUserOrServiceToDB(&uos, db)
		//	},
		//},
		//Seed{
		//	Name: "Service2",
		//	Run: func(db *database.DBConn) error {
		//		uos := userOrService.UserOrService{ID: 2, Username: "service2", Password: "service2"}
		//		return userOrService.AddUserOrServiceToDB(&uos, db)
		//	},
		//},
	}
}
