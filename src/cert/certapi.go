package cert

import (
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"log"
	"net/http"
	"saltsec/database"
	"saltsec/middleware"
	"time"

	"github.com/gorilla/mux"
	"go.step.sm/crypto/x509util"
)

type RequestCertDTO struct {
	Subject      pkix.Name `json:"subject"`
	Issuer       pkix.Name `json:"issuer"`
	KeyUsage     x509util.KeyUsage  `json:"keyUsage"`
	ExtKeyUsages *x509util.ExtKeyUsage  `json:"extKeyUsages"`
	IssuerSerial string    `json:"issuerSerial"`
	SubjectEmail string    `json:"subjectEmail"`
	IssuerEmail  string    `json:"issuerEmail"`
	Type 		 CertType  `json:"type"`
}

type ParamsDTO struct {
	KeyUsages    []string `json:"keyUsages"`
	ExtKeyUsages []string `json:"extKeyUsages"`
}

func AddCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		var dto RequestCertDTO
		err := dto.decode(r)
		if err != nil {
			middleware.JSONResponse(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		c := dto.parse()
		var cert *Certificate
		if dto.Type == Root {
			cert, err = GenCARootCert(c)
		} else {
			lookupDTO := LookupDTO {
				Username: dto.IssuerEmail,
				SerialNumber: dto.Issuer.SerialNumber,
			}
			cert, err = GenCert(c, lookupDTO)
		}

		if err != nil {
			middleware.JSONResponse(w, "Bad Request: "+err.Error(), http.StatusBadRequest)
			return
		}

		// cert.Verify()
		if err := cert.Save(); err != nil {
			middleware.JSONResponse(w, "Internal Server Error: " + err.Error(), http.StatusInternalServerError)
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

// func DownloadCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		params := mux.Vars(r)
// 		filename := findFileName(params["cn"], params["un"])

// 		w.Header().Set("Content-Disposition", "attachment; filename="+strconv.Quote(filename))
// 		w.Header().Set("Content-Type", "application/octet-stream")
// 		http.ServeFile(w, r, filename)
// 	}
// }

func GetCert(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cert := Certificate{}
		params := mux.Vars(r)
		dto := LookupDTO{Username: params["us"], SerialNumber: params["sn"]}
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

func GetArchived(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		archived := []ArchivedCert{}
		db.DB.Find(&archived)
		json.NewEncoder(w).Encode(archived)
	}
}

func CheckIfArchived(db *database.DBConn) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Checking archives...")
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
		var dto LookupDTO
		err := dto.decode(r)
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

func (dto *RequestCertDTO) parse() *x509.Certificate {
	serialNumber := GetRandomSerial()
	dto.Subject.SerialNumber = serialNumber.String()
	template := x509.Certificate{
		SerialNumber:   serialNumber,
		Subject:        dto.Subject,
		Issuer:         dto.Issuer,
		EmailAddresses: []string{dto.SubjectEmail},
		NotBefore:      time.Now().Add(-10 * time.Second),
		NotAfter:       time.Now().AddDate(int(dto.Type), 0, 0),
	}
	dto.KeyUsage.Set(&template)
	dto.ExtKeyUsages.Set(&template)

	if dto.Type == Root || dto.Type == Intermediary {
		template.BasicConstraintsValid = true
		template.IsCA = true
		template.MaxPathLen = -1
	}

	if dto.Type == Root {
		template.Issuer = dto.Subject
	}

	return &template
}

func (dto *RequestCertDTO) decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}

func (dto *LookupDTO) decode(r *http.Request) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(dto); err != nil {
		return err
	}
	return nil
}
