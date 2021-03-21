package main

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"saltsec/cert"
	"saltsec/database"
	"saltsec/globals"
	"saltsec/router"
	"saltsec/seeder"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
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

	// TODO(Jovan): Move to testing...
	rootTemplate := x509.Certificate{
		SerialNumber: big.NewInt(1), // TODO(Jovan): Random?
		Subject: pkix.Name{
			Country:      []string{"RS"},
			Organization: []string{"SaltyBois Inc."},
			CommonName:   "Root CA",
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		KeyUsage:              x509.KeyUsageCertSign,// | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageOCSPSigning,
		},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            2,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}
	_, rootCertPEM, _ := cert.GenCARootCert(&rootTemplate)
	log.Println("rootCert\n", string(rootCertPEM))

	r := router.Router{}
	r.R = mux.NewRouter()
	r.InitRouter(&db)

	handler := cors.Default().Handler(r.R)
	log.Printf("Starting main backend on port %s...\n", globals.PORT)
	log.Fatal(http.ListenAndServe(globals.PORT, handler))
}
