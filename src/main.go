package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "github.com/rs/cors"
    "github.com/gorilla/mux"
)

// TODO(Jovan): Use envvar for port
const (
    PORT = ":8081"
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
    router := mux.NewRouter()
    router.Use(jsonMiddleware)
    router.HandleFunc("/ping", Ping()).Methods("GET")

    handler := cors.Default().Handler(router)
    fmt.Printf("Starting main backend on port %s...\n", PORT)
    log.Fatal(http.ListenAndServe(PORT, handler))
}
