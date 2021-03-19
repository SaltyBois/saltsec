package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
	"strconv"
)

// TODO(Jovan): Use envvar for port
const (
	DEFAULT_PORT        = ":8081"
	DEFAULT_HOST_DB     = "localhost"
	DEFAULT_PORT_DB     = 5432
	DEFAULT_USER_DB     = "postgres"
	DEFAULT_PASSWORD_DB = "root"
	DEFAULT_NAME_DB     = "SaltyData"
)

var (
	PORT        string
	HOST_DB     string
	PORT_DB     int
	USER_DB     string
	PASSWORD_DB string
	NAME_DB     string
)

// TODO(Jovan): Make it cleaner???
func loadGlobalVars() {
	if val, present := os.LookupEnv("PORT_SALT"); present {
		PORT = val
	} else {
		PORT = DEFAULT_PORT
	}

	if val, present := os.LookupEnv("HOST_DB_SALT"); present {
		HOST_DB = val
	} else {
		HOST_DB = DEFAULT_HOST_DB
	}

	if val, present := os.LookupEnv("PORT_DB_SALT"); present {
		PORT_DB, _ = strconv.Atoi(val)
	} else {
		PORT_DB = DEFAULT_PORT_DB
	}

	if val, present := os.LookupEnv("USER_DB_SALT"); present {
		USER_DB = val
	} else {
		USER_DB = DEFAULT_USER_DB
	}

	if val, present := os.LookupEnv("PASSWORD_DB_SALT"); present {
		PASSWORD_DB = val
	} else {
		PASSWORD_DB = DEFAULT_PASSWORD_DB
	}

	if val, present := os.LookupEnv("NAME_DB_SALT"); present {
		NAME_DB = val
	} else {
		NAME_DB = DEFAULT_NAME_DB
	}
}

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

func connectToDb() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
		HOST_DB, USER_DB, PASSWORD_DB, NAME_DB, PORT_DB)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	return db, err
}

type DBTestStruct struct {
	ID   uint64
	Name string
}

func migrateData(db *gorm.DB) {
	/* TODO(Jovan): Add migrations; eg:
	   db.AutoMigrate(&package.Struct{})
	*/
	db.AutoMigrate(&DBTestStruct{})
}

func seedData(db *gorm.DB) {
	/* TODO(Jovan): Add seeding; eg:
	   // Check for existing data
	   data := package.Struct{}
	   result := db.First(&data)
	   if errors.Is(result.Error, gorm.ErrRecordNotFound) {
	       newData := package.Struct{ <fill_with_data>... }
	       ....
	       // Insert into DB
	       db.Create(&newData)
	   }
	*/
	data := DBTestStruct{}
	result := db.First(&data)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		newData := DBTestStruct{ID: 1234, Name: "Test One"}
		db.Create(&newData)
	}
}

func main() {
	loadGlobalVars()
	if _, dbpresent := os.LookupEnv("DB_DEV"); dbpresent {
		log.Println("DB_DEV set, using database...")
		db, err := connectToDb()
		if err != nil {
			panic("Failed to connect to database")
		}

		migrateData(db)
		seedData(db)
		log.Println("Successfully connected to Salty Database!")
	} else {
		log.Println("DB_DEV not set, not using database...")
	}

	router := mux.NewRouter()
	router.Use(jsonMiddleware)
	router.HandleFunc("/ping", Ping()).Methods("GET")

	handler := cors.Default().Handler(router)
	log.Printf("Starting main backend on port %s...\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, handler))
}
