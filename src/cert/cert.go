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
	"time"
)

const (
	EE_CERT_DIR    = "../certs/ee/"
	INTER_CERT_DIR = "../certs/inter/"
	ROOT_CERT_DIR  = "../certs/root/"
)

type Certificate struct {
	Cert       *x509.Certificate
	PEM        []byte
	PrivateKey *rsa.PrivateKey
	Type       CertType
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
	rootCert, rootPEM, err := genCert(rootTemplate, rootTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: rootCert, PEM: rootPEM, PrivateKey: privateKey}
	return &cert, nil
}

func GenCAIntermediateCert(template *x509.Certificate, issuerSerialNumber string) (*Certificate, error) {
	issuerCert, err := FindCert(issuerSerialNumber)
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
	caCert, certPEM, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerPrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: caCert, PEM: certPEM, PrivateKey: privateKey}
	return &cert, nil
}

func GenEndEntityCert(template *x509.Certificate, issuerSerialNumber string) (*Certificate, error) {
	issuerCert, err := FindCert(issuerSerialNumber)
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
	eeCert, certPEM, err := genCert(template, issuerCert.Cert, &privateKey.PublicKey, issuerPrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: eeCert, PEM: certPEM, PrivateKey: privateKey}
	return &cert, err
}

func ValidateCertChain(cert *x509.Certificate) error {
	// NOTE(Jovan): While issuer.SerialNumber != cert.SerialNumber, traverse
	if cert.SerialNumber.String() == cert.Issuer.SerialNumber {
		return validateCert(cert)
	}
	issuerCert, err := FindCert(cert.Issuer.SerialNumber)
	if err != nil {
		return err
	}
	return ValidateCertChain(issuerCert.Cert)
}

func FindCert(serialNumber string) (*Certificate, error) {
	cert := &Certificate{}
	err := cert.Load(serialNumber)
	return cert, err
}

func FindCertKey(serialNumber string) (*rsa.PrivateKey, error) {
	// TODO(Jovan): Get by serial
	return nil, errors.New("not implemented")
}

func (cert *Certificate) Save() error {
	filename := cert.Cert.SerialNumber.String() + ".pem"
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
	pemFile, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed creating PEM file, returned error: %s\n", err)
		return err
	}
	pemBlock, _ := pem.Decode(cert.PEM)
	if pemBlock == nil {
		log.Println("Failed decoding PEM block")
		return errors.New("failed decoding PEM block")
	}
	err = pem.Encode(pemFile, pemBlock)
	if err != nil {
		log.Printf("Failed writing to PEM file, returned error: %s\n", err)
		return err
	}
	return nil
}

func (cert *Certificate) Load(serialNumber string) error {
	filename := ROOT_CERT_DIR + serialNumber + ".pem"
	if err := cert.loadCertAndKey(filename); err == nil {
		return nil
	}
	filename = INTER_CERT_DIR + serialNumber + ".pem"
	if err := cert.loadCertAndKey(filename); err == nil {
		return nil
	}
	filename = EE_CERT_DIR + serialNumber + ".pem"
	if err := cert.loadCertAndKey(filename); err == nil {
		return nil
	}
	return errors.New("certificate/key file does not exist")
}

func (cert *Certificate) loadCertAndKey(filename string) error {
	certtmp, pemBlock, err := loadCertFile(filename)
	if err != nil {
		return err
	}
	cert.Cert = certtmp
	cert.Type = Root
	cert.PEM = pem.EncodeToMemory(pemBlock)
	key, err := loadKeyFile(filename)
	cert.PrivateKey = key
	return err
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

func validateCert(cert *x509.Certificate) error {
	today := time.Now()
	if cert.NotAfter.Before(today) {
		return errors.New("certificate date is not valid")
	}
	// TODO(Jovan) OCSP
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
