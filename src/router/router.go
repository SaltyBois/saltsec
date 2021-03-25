package router

import (
	"saltsec/admin"
	"saltsec/cert"
	"saltsec/database"
	"saltsec/middleware"

	"github.com/gorilla/mux"
)

type Router struct {
	R *mux.Router
}

func (r *Router) initRoutes(db *database.DBConn) {
	// NOTE(Jovan): Admin
	r.R.HandleFunc("/api/admin", admin.GetAll(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/root", cert.AddCARootCert(db)).Methods("POST")
	r.R.HandleFunc("/api/cert/params", cert.GetCertParams()).Methods("GET")
	r.R.HandleFunc("/api/cert/", cert.GetAllCerts(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/{sn}", cert.GetCert(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/archive/check/{sn}", cert.CheckIfArchived(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/archive/add/{sn}", cert.AddToArchive(db)).Methods("GET")
	r.R.HandleFunc("/api/cert/download/{sn}", cert.DownloadCert(db)).Methods("GET")
}

func (r *Router) InitRouter(db *database.DBConn) {
	r.R.Use(middleware.JSONMiddleware)
	r.initRoutes(db)
}
