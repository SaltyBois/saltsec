package seeder

import (
	"log"
	"saltsec/admin"
	"saltsec/cert"
	"saltsec/database"
)

type Seed struct {
	Name string
	Run  func(*database.DBConn) error
}

func MigrateData(db *database.DBConn) {
	db.DB.AutoMigrate(&admin.Admin{})
	db.DB.AutoMigrate(&cert.ArchivedCert{})
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
		Seed{
			Name: "CreateADmin2",
			Run: func(db *database.DBConn) error {
				a := admin.Admin{ID: 2, Username: "admin2", Email: "admin2@email.com", Password: "admin2"}
				return admin.AddAdmin(&a, db)
			},
		},
		Seed{
			Name: "Archive1",
			Run: func(db *database.DBConn) error {
				return cert.ArchiveCert(db, "1")
			},
		},
	}
}
