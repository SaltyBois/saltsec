package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"saltsec/database"
	"saltsec/globals"
	"saltsec/seeder"
	"saltsec/admin"
)

// TODO(Jovan): Move to common middleware?
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

type Pong struct {
	Message string `json: message`
}

func Ping() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		pong := Pong{Message: "pong"}
		json.NewEncoder(w).Encode(&pong)
	}
}

func initRoutes(router *mux.Router, db *database.DBConn) {
	router.HandleFunc("/ping", Ping()).Methods("GET")

	// NOTE(Jovan): Admin
	router.HandleFunc("/api/admin", admin.GetAll(db))
}

func main() {
	globals.LoadGlobalVars()
	db := database.DBConn{}
	if _, dbpresent := os.LookupEnv("DB_DEV"); dbpresent {
		log.Println("DB_DEV set, using database...")
		conn, err := database.ConnectToDb()
		if err != nil {
			panic("Failed to connect to database")
		}
		db.DB = conn

		seeder.MigrateData(&db)
		seeder.SeedData(&db)
		log.Println("Successfully connected to Salty Database!")
	} else {
		log.Println("DB_DEV not set, not using database...")
	}

	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	initRoutes(router, &db)

	handler := cors.Default().Handler(router)
	log.Printf("Starting main backend on port %s...\n", globals.PORT)
	log.Fatal(http.ListenAndServe(globals.PORT, handler))
}
