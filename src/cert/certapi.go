package cert

import (
	"bytes"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"saltsec/database"
	"saltsec/middleware"
	"saltsec/userOrService"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.step.sm/crypto/x509util"
)

type CertDTO struct {
	Type         CertType  `json:"type"`
	Issuer       UserDnDTO `json:"issuer"`
	Country      string    `json:"country"`
	Organization string    `json:"organization"`
	CommonName   string    `json:"commonName"`
	Password     string    `json:"password"`
	IsCA         bool      `json:"isCA"`
	KeyUsages    []string  `json:"keyUsages"`
	ExtKeyUsages []string  `json:"extKeyUsages"`
	EmailAddress string    `json:"emailAddress"`
}

type ParamsDTO struct {
	KeyUsages    []string `json:"keyUsages"`
	ExtKeyUsages []string `json:"extKeyUsages"`
}

type UserDnDTO struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	CommonName string `json:"commonName"`
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
		rootTemplate.Issuer.SerialNumber = rootTemplate.SerialNumber.String()

		cert, err := GenCARootCert(rootTemplate)
		if err != nil {
			middleware.JSONResponse(w, "Internal Server Error failed to generate certificate", http.StatusInternalServerError)
			return
		}
		cert.Type = GetType(cert.Cert)
		log.Printf("Generated cert: %s\n", cert.Cert.SerialNumber.String())
		json.NewEncoder(w).Encode(cert.Cert.SerialNumber.String())
		cert.Save(dto.EmailAddress, dto.Password)
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
		entity := userOrService.UserOrService{}
		err = userOrService.GetUosByUsername(dto.Issuer.Username, &entity, db)
		if err != nil {
				middleware.JSONResponse(w, "Internal Server Error failed to get user", http.StatusInternalServerError)
				return
		}
		cert, err := GenCAIntermediateCert(caTemplate, dto.Issuer, dto.Password)
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
		cert, err := GenEndEntityCert(eeTemplate, dto.Issuer, dto.Password)
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

func DownloadCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		filename := findFileName(params["cn"], params["un"])

		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
		w.Header().Set("Content-Type", "application/octet-stream")
		http.ServeFile(w, r, filename)
	}
}

func GetCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cert := Certificate{}
		params := mux.Vars(r)
		dto := UserDnDTO{Username: params["us"], Password: params["ps"], CommonName: params["cn"]}
		cert.Load(dto)
		json.NewEncoder(w).Encode(cert)
	}
}

func GetAllCerts(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		certs := []Certificate{}
		LoadAll(db, &certs)
		json.NewEncoder(w).Encode(certs)
	}
}

func CheckIfArchived(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		serialNumber := params["sn"]
		log.Printf("Looking for serial: %s", serialNumber)
		if !IsArchived(db, serialNumber) {
			middleware.JSONResponse(w, "Not Found Certificate is not archived", http.StatusNotFound)
			return
		}
		middleware.JSONResponse(w, "OK Certificate is archived", http.StatusOK)
	}
}

func AddToArchive(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var dto UserDnDTO
		err := dto.loadUserDnDTO(r)
		if err != nil {
			middleware.JSONResponse(w, "Internal Server Error "+err.Error(), http.StatusInternalServerError)
			return
		}

		if err := ArchiveCert(db, dto); err != nil {
			middleware.JSONResponse(w, "Bad Request "+err.Error(), http.StatusNotFound)
			return
		}
		middleware.JSONResponse(w, "OK Certificate archived", http.StatusOK)
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
		BasicConstraintsValid: dto.IsCA,
		IsCA:                  dto.IsCA,
		// NOTE(Jovan): -1 = unset -> No limit for how many certs can be
		// "under" current CA
		MaxPathLen: -1,
		//IPAddresses: []net.IP{net.ParseIP(dto.IPAddress)},
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

func (dto *UserDnDTO) loadUserDnDTO(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}
