package cert

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
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
	Type         CertType `json:"type"`
	Country      string   `json:"country"`
	Organization string   `json:"organization"`
	CommonName   string   `json:"commonName"`
	IPAddress    string   `json:"ipAddress"`
	IsCA         bool     `json:"isCA"`
	KeyUsages    []string `json:"keyUsages"`
	ExtKeyUsages []string `json:"extKeyUsages"`
	IssuerSerial string   `json:"issuerSerial"`
}

type ParamsDTO struct {
	KeyUsages    []string `json:"keyUsages"`
	ExtKeyUsages []string `json:"extKeyUsages"`
}

func AddCARootCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var dto CertDTO
		err := dto.loadCertDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		if dto.Type != Root {
			middleware.JSONResponse(w, "Bad Request cert type not of type 'Root'", http.StatusBadRequest)
			return
		}
		rootTemplate := dto.parseCertDTO()
		setKeyUsages(rootTemplate, dto.KeyUsages)
		setExtKeyUsages(rootTemplate, dto.ExtKeyUsages)

		cert, err := GenCARootCert(rootTemplate)
		if err != nil {
			middleware.JSONResponse(w, "Internal Server Error failed to generate certificate", http.StatusInternalServerError)
			return
		}
		log.Printf("Generated cert: %s\n", cert.Cert.SerialNumber.String())
		json.NewEncoder(w).Encode(cert.Cert.SerialNumber.String())
	}
}

func AddCACert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto CertDTO
		err := dto.loadCertDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		if dto.Type != Intermediary {
			middleware.JSONResponse(w, "Bad Request cert type not of type 'Intermediary'", http.StatusBadRequest)
			return
		}
		caTemplate := dto.parseCertDTO()
		setKeyUsages(caTemplate, dto.KeyUsages)
		setExtKeyUsages(caTemplate, dto.ExtKeyUsages)
		cert, err := GenCAIntermediateCert(caTemplate, dto.IssuerSerial)
		if err != nil {
			middleware.JSONResponse(w, "Internal Server Error failed to generate certificate", http.StatusInternalServerError)
			return
		}
		log.Printf("Generated cert: %s\n", cert.Cert.SerialNumber.String())
		json.NewEncoder(w).Encode(cert.Cert.SerialNumber.String())
	}
}

func AddEECert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto CertDTO
		err := dto.loadCertDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Request cert type not of type 'EndEntity", http.StatusBadRequest)
			return
		}
		eeTemplate := dto.parseCertDTO()
		setKeyUsages(eeTemplate, dto.KeyUsages)
		setExtKeyUsages(eeTemplate, dto.ExtKeyUsages)
		cert, err := GenEndEntityCert(eeTemplate, dto.IssuerSerial)
		if err != nil {
			middleware.JSONResponse(w, "Internal Server Error failed to generate certificate", http.StatusInternalServerError)
			return
		}
		log.Printf("Generated cert: %s\n", cert.Cert.SerialNumber.String())
		json.NewEncoder(w).Encode(cert.Cert.SerialNumber.String())
	}
}

func GetCertParams() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto ParamsDTO
		keyUsages := []string{
			"KeyUsageDigitalSignature",
			"KeyUsageContentCommitment",
			"KeyUsageKeyEncipherment",
			"KeyUsageDataEncipherment",
			"KeyUsageKeyAgreement",
			"KeyUsageCertSign",
			"KeyUsageCRLSign",
			"KeyUsageEncipherOnly",
			"KeyUsageDecipherOnly",
		}

		extKeyUsages := []string{
			"ExtKeyUsageAny",
			"ExtKeyUsageServerAuth",
			"ExtKeyUsageClientAuth",
			"ExtKeyUsageCodeSigning",
			"ExtKeyUsageEmailProtection",
			"ExtKeyUsageIPSECEndSystem",
			"ExtKeyUsageIPSECTunnel",
			"ExtKeyUsageIPSECUser",
			"ExtKeyUsageTimeStamping",
			"ExtKeyUsageOCSPSigning",
			"ExtKeyUsageMicrosoftServerGatedCrypto",
			"ExtKeyUsageNetscapeServerGatedCrypto",
			"ExtKeyUsageMicrosoftCommercialCodeSigning",
			"ExtKeyUsageMicrosoftKernelCodeSigning",
		}

		dto.KeyUsages = keyUsages
		dto.ExtKeyUsages = extKeyUsages
		json.NewEncoder(w).Encode(dto)
	}
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

func (dto *CertDTO) parseCertDTO() *x509.Certificate {
	rootTemplate := x509.Certificate{
		SerialNumber: GetRandomSerial(),
		Subject: pkix.Name{
			Country:      []string{dto.Country},
			Organization: []string{dto.Organization},
			CommonName:   dto.CommonName,
		},
		NotBefore: time.Now().Add(-10 * time.Second),
		NotAfter:  time.Now().AddDate(int(dto.Type), 0, 0),
		// NOTE(Jovan): Used for MaxPathLen
		BasicConstraintsValid: false,
		IsCA:                  dto.IsCA,
		// NOTE(Jovan): -1 = unset -> No limit for how many certs can be
		// "under" current CA
		MaxPathLen:  -1,
		IPAddresses: []net.IP{net.ParseIP(dto.IPAddress)},
	}
	return &rootTemplate
}

func (dto *CertDTO) loadCertDTO(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}
