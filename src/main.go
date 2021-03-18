package main

import (
    "encoding/json"
    "log"
    "fmt"
    "os"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
)

// TODO(Jovan): Use envvar for port
const (
    PORT        = ":8081"
    HOST_DB     = "localhost"
    PORT_DB     = 5432
    USER_DB     = "postgres"
    PASSWORD_DB = "root"
    NAME_DB     = "SaltyData"
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

func connectToDb() (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d",
        HOST_DB, USER_DB, PASSWORD_DB, NAME_DB, PORT_DB)
        db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
        return db, err
}

func migrateData(db *gorm.DB) {
    /* TODO(Jovan): Add migrations; eg: 
    db.AutoMigrate(&package.Struct{})
    */
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
}

func main() {

    if _, dbpresent := os.LookupEnv("DB_DEV"); !dbpresent {
        log.Println("DB_DEV set to true, not using database...")
        db, err := connectToDb()
        if err != nil {
            panic("Failed to connect to database")
        }

        migrateData(db)
        seedData(db)
        log.Println("Successfully connected to Salty Database!")
    }


    router := mux.NewRouter()
    router.Use(jsonMiddleware)
    router.HandleFunc("/ping", Ping()).Methods("GET")

    handler := cors.Default().Handler(router)
    log.Printf("Starting main backend on port %s...\n", PORT)
    log.Fatal(http.ListenAndServe(PORT, handler))
}
