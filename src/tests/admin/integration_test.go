package tests

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"saltsec/admin"
	"saltsec/database"
	"saltsec/globals"
	"saltsec/router"
	"testing"
)

type Application struct {
	db      database.DBConn
	handler http.Handler
}

func seedAdmins(db *database.DBConn) (*gorm.DB, error) {
	admins := []admin.Admin{
		admin.Admin{
			ID:       123,
			Username: "AdminOneTwoThree",
			Email:    "adminonetwothree@email.com",
			Password: "AdminOneTwoThree",
		},
	}
	db.DB.AutoMigrate(admin.Admin{})
	tx := db.DB.Begin()
	tx.SavePoint("sp1")
	for _, a := range admins {
		if err := db.DB.Create(&a).Error; err != nil {
			return tx, err
		}
	}
	return tx, nil
}

func (a *Application) serveHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

func (a *Application) startServer() {
	router := router.Router{}
	router.R = mux.NewRouter()
	router.InitRouter(&a.db)
	a.handler = cors.Default().Handler(router.R)
	// log.Fatal(http.ListenAndServe(globals.PORT, a.handler))
}

func TestAdminGet(t *testing.T) {
	t.Log("Running Admin integration test...")
	globals.LoadGlobalVars()
	a := Application{}
	a.db.ConnectToDb()
	a.startServer()
	tx, err := seedAdmins(&a.db)
	if err != nil {
		tx.RollbackTo("sp1")
		t.Fatalf("Seeding returned error: %s\n", err)
	}

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/admin"), nil)
	if err != nil {
		t.Errorf("Error creating request: %s\n", err)
	}

	w := httptest.NewRecorder()
	a.serveHTTP(w, req)

	if want, got := http.StatusOK, w.Code; want != got {
		t.Errorf("Expected status code: %v, got: %v\n", want, got)
	}

	t.Log("Rolling back...")
	tx.RollbackTo("sp1")
}
