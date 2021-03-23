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

type CertType int

const (
	Root         CertType = 10
	Intermediary CertType = 5
	EndEntity    CertType = 1
)

func (ct CertType) String() string {
	return toString[ct]
}

var toString = map[CertType]string{
	Root:         "Root",
	Intermediary: "Intermediary",
	EndEntity:    "EndEntity",
}

var toID = map[string]CertType{
	"Root":         Root,
	"Intermediary": Intermediary,
	"EndEntity":    EndEntity,
}

func (ct CertType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toString[ct])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (ct *CertType) UnmarshalJSON(b []byte) error {
	var jsonString string
	err := json.Unmarshal(b, &jsonString)
	if err != nil {
		log.Fatalf("Failed to unmarshal CertType, returned error: %s\n", err)
		return err
	}
	*ct = toID[jsonString]
	return nil
}

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
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&dto); err != nil {
			middleware.JSONResponse(w, "Bad Request"+err.Error(), http.StatusBadRequest)
			return
		}
		rootTemplate := parseCertDTO(&dto)
		setKeyUsages(rootTemplate, dto.KeyUsages)
		setExtKeyUsages(rootTemplate, dto.ExtKeyUsages)

		_, pem, _ := GenCARootCert(rootTemplate)
		log.Printf("Generated cert: %s\n", string(pem))
		json.NewEncoder(w).Encode(base64.StdEncoding.EncodeToString(pem))
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

func parseCertDTO(dto *CertDTO) *x509.Certificate {

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
