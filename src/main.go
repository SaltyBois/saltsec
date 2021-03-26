package main

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"saltsec/cert"
	"saltsec/database"
	"saltsec/globals"
	"saltsec/router"
	"saltsec/seeder"
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
	//var rootTemplate = x509.Certificate{
	//	SerialNumber: big.NewInt(1),
	//	Subject: pkix.Name{
	//		Country:      []string{"SE"},
	//		Organization: []string{"Company Co."},
	//		CommonName:   "Root CA",
	//	},
	//	NotBefore:             time.Now().Add(-10 * time.Second),
	//	NotAfter:              time.Now().AddDate(10, 0, 0),
	//	KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
	//	ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
	//	BasicConstraintsValid: true,
	//	IsCA:                  true,
	//	MaxPathLen:            2,
	//	IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	//}
	//c, _ := cert.GenCARootCert(&rootTemplate)
	//c.Type = cert.Root
	//c.Save("test123")
	cert.Init()

	r := router.Router{}
	r.R = mux.NewRouter()
	r.InitRouter(&db)

	handler := cors.Default().Handler(r.R)
	log.Printf("Starting main backend on port %s...\n", globals.PORT)
	log.Fatal(http.ListenAndServe(globals.PORT, handler))
}
