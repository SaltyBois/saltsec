package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"saltsec/database"
	"saltsec/keystore"
	"strings"
	"time"
)

var (
	ROOT_DIR       = filepath.FromSlash("../certs/")
	EE_CERT_DIR    = filepath.FromSlash(ROOT_DIR + "ee/")
	INTER_CERT_DIR = filepath.FromSlash(ROOT_DIR + "inter/")
	ROOT_CERT_DIR  = filepath.FromSlash(ROOT_DIR + "root/")
)

const FILE_EXTENSION = ".pfx"

type Certificate struct {
	Cert       *x509.Certificate   `json:"cert"`
	CertChain  []*x509.Certificate `json:"certChain"`
	PrivateKey *rsa.PrivateKey     `json:"-"`
	Type       CertType            `json:"type"`
}

type ArchivedCert struct {
	SerialNumber string    `gorm:"primaryKey" json:"serialNumber"`
	ArchiveDate  time.Time `json:"archiveDate"`
}

type LookupDTO struct {
	Username     string `json:"username"`
	SerialNumber string `json:"serialNumber"`
}

func Init() {
	if _, err := os.Stat("../certs/"); os.IsNotExist(err) {
		os.Mkdir("../certs/", 0755)
	}
	if _, err := os.Stat(EE_CERT_DIR); os.IsNotExist(err) {
		os.Mkdir(EE_CERT_DIR, 0755)
	}
	if _, err := os.Stat(INTER_CERT_DIR); os.IsNotExist(err) {
		os.Mkdir(INTER_CERT_DIR, 0755)
	}
	if _, err := os.Stat(ROOT_CERT_DIR); os.IsNotExist(err) {
		os.Mkdir(ROOT_CERT_DIR, 0755)
	}
}

func GetRandomSerial() *big.Int {
	z := new(big.Int)
	b, err := genRandomBytes(4)
	if err != nil {
		log.Fatalf("Failed to generate random serial, returned error: %s\n", err)
	}
	z.SetBytes(b)
	return z
}

func GenCARootCert(rootTemplate *x509.Certificate) (*Certificate, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Failed to generate key, returned error: %s\n", err)
		return nil, err
	}
	// NOTE(Jovan): Generate a self-signed cert
	rootCert, _, err := genCert(rootTemplate, rootTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	rootCert.Issuer.SerialNumber = rootCert.SerialNumber.String()
	cert := Certificate{Cert: rootCert, PrivateKey: privateKey, Type: Root}
	return &cert, nil
}

func GenCert(template *x509.Certificate, issuerDTO LookupDTO) (*Certificate, error) {

	issuerCert := Certificate{}
	if err := issuerCert.Load(issuerDTO); err != nil {
		log.Printf("Failed to load issuer cert, returned error: %s\n", err)
		return nil, err
	}

	// TODO(Jovan): Validate issuer

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Failed to generate key, returned error: %s\n", err)
		return nil, err
	}
	c, _, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerCert.PrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	certChain := issuerCert.CertChain
	certChain = append([]*x509.Certificate{issuerCert.Cert}, certChain...)

	cert := Certificate{
		Cert: c,
		CertChain: certChain,
		PrivateKey: privateKey,
		Type: GetType(c),
	}
	return &cert, err
}

func LoadCert(dto LookupDTO) (*Certificate, error) {
	cert := &Certificate{}
	err := cert.Load(dto)
	return cert, err
}

func (cert *Certificate) Verify() error {
	// TODO(Jovan): Implement
	return errors.New("not implemented")
}

func (cert *Certificate) Save() error {
	filename := cert.Cert.EmailAddresses[0] + cert.Cert.SerialNumber.String()
	err := keystore.WritePFX(cert.Cert, cert.CertChain, cert.PrivateKey, filename)
	if err != nil {
		return err
	}
	log.Printf("Written file: %s\n", filename)
	return nil
}

func (cert *Certificate) Load(dto LookupDTO) error {
	if err := cert.loadCertAndKey(dto.Username + dto.SerialNumber); err != nil {
		return err
	}
	return nil
}

func LoadAll(db *database.DBConn, certs *[]Certificate) error {
	paths := []string{
		keystore.ROOT_DIR,
	}

	for _, root := range paths {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				filename := strings.ReplaceAll(path, keystore.ROOT_DIR, "")
				filename = strings.ReplaceAll(filename, keystore.FILE_EXT, "")
				_, c, cChain, err := keystore.ReadPFX(filename)
				if err != nil {
					return err
				}
				cert := Certificate{
					Cert:       c,
					CertChain:  cChain,
					Type:       GetType(c),
					PrivateKey: nil,
				}
				*certs = append(*certs, cert)
			}
			return nil
		})
	}
	return nil
}

func ArchiveCert(db *database.DBConn, dto LookupDTO) error {
	cert := Certificate{}
	cert.Load(dto)
	archivedCert := ArchivedCert{SerialNumber: cert.Cert.SerialNumber.String(), ArchiveDate: time.Now()}
	return db.DB.Create(&archivedCert).Error
}

func IsArchived(db *database.DBConn, serialNumber string) bool {
	archive := ArchivedCert{}
	return db.DB.First(&archive, serialNumber).Error == nil
}

func GetType(c *x509.Certificate) CertType {
	log.Println(c.Issuer.SerialNumber)
	log.Println(c.SerialNumber.String())
	if !c.IsCA {
		return EndEntity
	} else if c.SerialNumber.String() != c.Issuer.SerialNumber {
		return Intermediary
	} else {
		return Root
	}
}

func findCertFile(serialNumber string) (string, error) {
	paths := []string{
		ROOT_CERT_DIR,
		INTER_CERT_DIR,
		EE_CERT_DIR,
	}

	for _, path := range paths {
		filename := path + serialNumber + ".pem"
		if _, err := os.Stat(filename); os.IsNotExist(err) {
			log.Printf("Looking for %s", filename)
			continue
		}
		return filename, nil
	}
	return "", errors.New("file does not exist")
}

func (cert *Certificate) loadCertAndKey(filename string) error {
	privateKey, c, cChain, err := keystore.ReadPFX(filename)
	if err != nil {
		return err
	}

	cert.Cert = c
	cert.PrivateKey = privateKey
	cert.CertChain = cChain
	cert.Type = GetType(cert.Cert)

	return nil
}

func validateCert(db *database.DBConn, cert *x509.Certificate) error {
	today := time.Now()
	if cert.NotAfter.Before(today) {
		return errors.New("certificate date is not valid")
	}
	if !IsArchived(db, cert.SerialNumber.String()) {
		return errors.New("certificate no longer valid; tagged as archived")
	}
	return nil
}

func genRandomBytes(length int) ([]byte, error) {
	if length <= 0 {
		log.Printf("Error generating random bytes, length less than 0\n")
		return nil, errors.New("failed generating random bytes, length less than 0")
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		log.Printf("Failed reading random bytes, returned error: %s\n", err)
		return nil, err
	}
	return b, nil
}

func genCert(template, parent *x509.Certificate, subjectKey *rsa.PublicKey, issuerKey *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, subjectKey, issuerKey)
	if err != nil {
		log.Printf("Failed to create certificate, returned error: %s\n", err)
		return nil, nil, err
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate, returned error: %s\n", err)
		return nil, nil, err
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	certPEM := pem.EncodeToMemory(&b)

	return cert, certPEM, nil
}
