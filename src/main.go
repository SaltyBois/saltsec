package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"saltsec/database"
	"saltsec/globals"
	"saltsec/router"
	"saltsec/seeder"
	"saltsec/cert"
)

func main() {
	globals.LoadGlobalVars()
	db := database.DBConn{}
	if _, dbpresent := os.LookupEnv("DB_DEV"); dbpresent {
		log.Println("DB_DEV set, using database...")
		err := db.ConnectToDb()
		if err != nil {
			panic("Failed to connect to database")
		}

		seeder.MigrateData(&db)
		seeder.SeedData(&db)
		log.Println("Successfully connected to Salty Database!")
	} else {
		log.Println("DB_DEV not set, not using database...")
	}

	// TODO(Jovan): Move to testing...
	_, rootCertPEM, _ := cert.GenCARoot()
	log.Println("rootCert\n", string(rootCertPEM))

	r := router.Router{}
	r.R = mux.NewRouter()
	r.InitRouter(&db)

	handler := cors.Default().Handler(r.R)
	log.Printf("Starting main backend on port %s...\n", globals.PORT)
	log.Fatal(http.ListenAndServe(globals.PORT, handler))
}
