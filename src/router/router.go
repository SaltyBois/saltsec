package router

import (
	"saltsec/admin"
	"saltsec/cert"
	"saltsec/database"
	"saltsec/middleware"
	"saltsec/user"

	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
}

func (r *Router) initRoutes(db *database.DBConn) {
	// NOTE(Jovan): Admin
	r.R.HandleFunc("/api/admin", admin.GetFirstAdmin(db)).Methods("GET")
	r.R.HandleFunc("/api/cert", cert.AddCert(db)).Methods("POST")
	r.R.HandleFunc("/api/cert/params", cert.GetCertParams()).Methods("GET")
	r.R.HandleFunc("/api/cert", cert.GetAllCerts(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/", cert.GetCert(db)).Queries("us", "{us}", "sn", "{sn}").Methods("GET")
	r.R.HandleFunc("/api/cert/archive/check/{sn}", cert.CheckIfArchived(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/archive/add", cert.AddToArchive(db)).Methods("POST")
	r.R.HandleFunc("/api/cert/archive", cert.GetArchived(db)).Methods("GET")
	// r.R.HandleFunc("/api/cert/download/{sn}", cert.DownloadCert(db)).Methods("GET")
	r.R.HandleFunc("/api/uos", user.GetAllUsers(db)).Methods("GET")
	r.R.HandleFunc("/api/uos/{username}", user.GetUser(db)).Methods("GET")
	r.R.HandleFunc("/api/uos/add", user.AddUserAndCert(db)).Methods("POST")

}

func (r *Router) InitRouter(db *database.DBConn) {
	r.R.Use(middleware.JSONMiddleware)
	r.initRoutes(db)
}
