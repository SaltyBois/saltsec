package router

import (
	"github.com/gorilla/mux"
	"net/http"
	"saltsec/admin"
	"saltsec/database"
)

type Router struct {
	R *mux.Router
}

// TODO(Jovan): Move to common middleware?
func jsonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}


func (r *Router) initRoutes(db *database.DBConn) {
	// NOTE(Jovan): Admin
	r.R.HandleFunc("/api/admin", admin.GetAll(db))
}

func (r *Router) InitRouter(db *database.DBConn) {
	r.R.Use(jsonMiddleware)
	r.initRoutes(db)
}