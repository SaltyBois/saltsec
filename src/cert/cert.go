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
	"saltsec/userOrService"
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
	Cert       *x509.Certificate
	PrivateKey *rsa.PrivateKey
	Type       CertType
}

type ArchivedCert struct {
	SerialNumber string    `gorm:"primaryKey" json:"serialNumber"`
	ArchiveDate  time.Time `json:"archiveDate"`
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
	b, err := genRandomBytes(256)
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

func GenCAIntermediateCert(template *x509.Certificate, issuerSerialNumber, password string) (*Certificate, error) {
	issuerCert, err := FindCert(issuerSerialNumber, password)
	if err != nil {
		log.Printf("Failed to find issuer cert, returned error: %s\n", err)
		return nil, err
	}

	issuerPrivateKey, err := FindCertKey(issuerSerialNumber)
	if err != nil {
		log.Printf("Failed to find issuer PK, returned error: %s\n", err)
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Failed to generate key, returned error: %s\n", err)
		return nil, err
	}
	caCert, _, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerPrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: caCert, PrivateKey: privateKey}
	return &cert, nil
}

func GenEndEntityCert(template *x509.Certificate, issuerSerialNumber, password string) (*Certificate, error) {
	issuerCert, err := FindCert(issuerSerialNumber, password)
	if err != nil {
		return nil, err
	}

	issuerPrivateKey, err := FindCertKey(issuerSerialNumber)
	if err != nil {
		log.Printf("Failed to get issuer PK, returned error: %s\n", err)
		return nil, err
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Printf("Failed to generate key, returned error: %s\n", err)
		return nil, err
	}
	eeCert, _, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerPrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: eeCert, PrivateKey: privateKey}
	return &cert, err
}

func ValidateCertChain(db *database.DBConn, cert *x509.Certificate) error {
	// NOTE(Jovan): While issuer.SerialNumber != cert.SerialNumber, traverse
	if cert.SerialNumber.String() == cert.Issuer.SerialNumber {
		return validateCert(db, cert)
	}
	// TODO(Jovan): get password
	password := "test123"
	issuerCert, err := FindCert(cert.Issuer.SerialNumber, password)
	if err != nil {
		return err
	}
	return ValidateCertChain(db, issuerCert.Cert)
}

func FindCert(serialNumber, password string) (*Certificate, error) {
	cert := &Certificate{}
	err := cert.Load(serialNumber, password)
	return cert, err
}

func FindCertKey(serialNumber string) (*rsa.PrivateKey, error) {
	// TODO(Jovan): Get by serial
	return nil, errors.New("not implemented")
}

func (cert *Certificate) Save(username, password string) error {
	// TODO(jovan): DN as filename
	filename := cert.Cert.Subject.CommonName + username
	//cert.Cert.SerialNumber.String() + ".pem"
	switch cert.Type {
	case Root:
		filename = ROOT_CERT_DIR + filename
	case Intermediary:
		filename = INTER_CERT_DIR + filename
	case EndEntity:
		filename = EE_CERT_DIR + filename
	default:
		return errors.New("invalid certificate type")
	}
	filename = filename + ".pfx"
	err := keystore.WritePFX(cert.Cert, cert.PrivateKey, password, filename)
	if err != nil {
		return err
	}
	return nil
}

func (cert *Certificate) Load(serialNumber, password string) error {
	filename := ROOT_CERT_DIR + serialNumber + FILE_EXTENSION
	if err := cert.loadCertAndKey(filename, password); err == nil {
		return nil
	}
	filename = INTER_CERT_DIR + serialNumber + FILE_EXTENSION
	if err := cert.loadCertAndKey(filename, password); err == nil {
		return nil
	}
	filename = EE_CERT_DIR + serialNumber + FILE_EXTENSION
	if err := cert.loadCertAndKey(filename, password); err == nil {
		return nil
	}
	return errors.New("certificate/key file does not exist")
}

func LoadAll(db *database.DBConn, certs *[]Certificate) error {

	paths := []string{
		ROOT_CERT_DIR,
		INTER_CERT_DIR,
		EE_CERT_DIR,
	}

	entities := []userOrService.UserOrService{}
	err := userOrService.GetAllUserOrServices(&entities, db)
	if err != nil {
		return err
	}
	certFiles := []string{}
	userdns := []string{}
	for _, root := range paths {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if !info.IsDir() {
				certFiles = append(certFiles, path)
				userdns = append(userdns, getUserDn(path))
			}
			return nil
		})
	}
	for _, path := range certFiles {
		for _, e := range entities {
			if strings.Contains(getUserDn(path), e.Username) {

				privateKey, c, err := keystore.ReadPFX(path, e.Password)

				if err != nil {
					return err
				}
				cert := Certificate{Cert: c, PrivateKey: privateKey, Type: GetType(c)}
				*certs = append(*certs, cert)
			}
		}
	}
	return nil
}

func getUserDn(path string) string {

	paths := []string{
		ROOT_CERT_DIR,
		INTER_CERT_DIR,
		EE_CERT_DIR,
	}

	for _, root := range paths {
		if strings.Contains(path, root) {
			filename := strings.ReplaceAll(path, root, "")
			userdn := strings.ReplaceAll(filename, FILE_EXTENSION, "")
			return userdn
		}
	}
	return ""
}

func ArchiveCert(db *database.DBConn, serialNumber string) error {
	archivedCert := ArchivedCert{SerialNumber: serialNumber, ArchiveDate: time.Now()}
	return db.DB.Create(&archivedCert).Error
}

func IsArchived(db *database.DBConn, serialNumber string) bool {
	certs := []ArchivedCert{}
	if err := GetArchived(db, &certs); err != nil {
		log.Printf("Failed getting archived certificates, returned error: %s\n", err)
		return false
	}
	for _, c := range certs {
		if serialNumber == c.SerialNumber {
			return true
		}
	}
	return false
}

func GetArchived(db *database.DBConn, certificates *[]ArchivedCert) error {
	return db.DB.Find(certificates).Error
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

func (cert *Certificate) loadCertAndKey(filename, password string) error {
	privateKey, c, err := keystore.ReadPFX(filename, password)
	if err != nil {
		return err
	}

	cert.Cert = c
	cert.PrivateKey = privateKey
	if !cert.Cert.IsCA {
		cert.Type = EndEntity
	} else if cert.Cert.SerialNumber.String() == cert.Cert.Issuer.SerialNumber {
		cert.Type = Root
	} else {
		cert.Type = Intermediary
	}

	return nil
}

func loadCertFile(filename string) (*x509.Certificate, *pem.Block, error) {
	certPEMBlock, err := loadPEMFile(filename)
	if err != nil {
		log.Printf("Failed to load PEM file, returned error: %s\n", err)
		return nil, nil, err
	}

	cert, err := x509.ParseCertificate(certPEMBlock.Bytes)
	if err != nil {
		log.Printf("Failed parsing PEM bytes to cert, returned error: %s\n", err)
		return nil, nil, err
	}
	return cert, certPEMBlock, nil
}

func loadKeyFile(filename string) (*rsa.PrivateKey, error) {
	keyPEMBytes, err := loadPEMFile(filename)
	if err != nil {
		log.Printf("Failed to load PEM file, returned error: %s\n", err)
		return nil, err
	}
	key, err := x509.ParsePKCS1PrivateKey(keyPEMBytes.Bytes)
	if err != nil {
		log.Printf("Failed parsing PEM bytes to key, returned error: %s\n", err)
		return nil, err
	}
	return key, nil
}

func loadPEMFile(filename string) (*pem.Block, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		log.Printf("Failed loading cert file, returned error: %s\n", err)
		return nil, err
	}
	pemBytes, _ := pem.Decode(fileBytes)
	return pemBytes, nil
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

func genCert(template, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*x509.Certificate, []byte, error) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
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
