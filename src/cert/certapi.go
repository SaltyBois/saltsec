package cert

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/gob"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"saltsec/database"
	"saltsec/middleware"
	"time"

	"go.step.sm/crypto/x509util"
)

type CertDTO struct {
	Country, Organization, CommonName string
	KeyUsages                         []string
	ExtKeyUsages                      []string
	IsCA                              bool
	IPAddress                         string
}

func setExtKeyUsages(cert *x509.Certificate, usages []string) {
	var extKeyUsage x509util.ExtKeyUsage
	extKeyUsageBytes := &bytes.Buffer{}
	gob.NewEncoder(extKeyUsageBytes).Encode(usages)
	extKeyUsage.UnmarshalJSON(extKeyUsageBytes.Bytes())
	extKeyUsage.Set(cert)
}

func setKeyUsages(cert *x509.Certificate, usages []string) {
	var keyUsage x509util.KeyUsage
	keyUsageBytes := &bytes.Buffer{}
	gob.NewEncoder(keyUsageBytes).Encode(usages)
	keyUsage.UnmarshalJSON(keyUsageBytes.Bytes())
	keyUsage.Set(cert)
}

func parseCARootDTO(dto *CertDTO) *x509.Certificate {

	rootTemplate := x509.Certificate{
		SerialNumber: GetRandomSerial(),
		Subject: pkix.Name{
			Country:      []string{dto.Country},
			Organization: []string{dto.Organization},
			CommonName:   dto.CommonName,
		},
		NotBefore:             time.Now().Add(-10 * time.Second),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  dto.IsCA,
		MaxPathLen:            2,
		IPAddresses:           []net.IP{net.ParseIP(dto.IPAddress)},
	}
	return &rootTemplate
}

func GetCARootCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var dto CertDTO
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&dto); err != nil {
			middleware.JSONResponse(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		rootTemplate := parseCARootDTO(&dto)
		setKeyUsages(rootTemplate, dto.KeyUsages)
		setExtKeyUsages(rootTemplate, dto.ExtKeyUsages)

		_, pem, _ := GenCARootCert(rootTemplate)
		log.Print(string(pem))
		middleware.JSONResponse(w, "Success", http.StatusOK)
		json.NewEncoder(w).Encode(base64.StdEncoding.EncodeToString(pem))
	}
}
