package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "log"
    "net/http"
    _ "github.com/lib/pq"
)

// TODO(Jovan): Use envvar for port
const (
    PORT     = ":8081"
    HOST     = "localhost"
    PORT_DB  = 5432
    USER     = "postgres"
    PASSWORD = "root"
    DBNAME   = "SaltyData"
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

func main() {

    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
        HOST, PORT_DB, USER, PASSWORD, DBNAME)

    db, err := sql.Open("postgres", psqlInfo)
    if err != nil {
        panic(err)
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        panic(err)
    }

    fmt.Println("Successfully connected to Salty Database!")

    router := mux.NewRouter()
    router.Use(jsonMiddleware)
    router.HandleFunc("/ping", Ping()).Methods("GET")

    handler := cors.Default().Handler(router)
    fmt.Printf("Starting main backend on port %s...\n", PORT)
    log.Fatal(http.ListenAndServe(PORT, handler))
}
