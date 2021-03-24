package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"math/big"
	"time"
)

type Certificate struct {
	Cert       *x509.Certificate
	PEM        []byte
	PrivateKey *rsa.PrivateKey
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
	caCert, certPEM, err := genCert(template, issuerCert, &privateKey.PublicKey, issuerPrivateKey)
	if err != nil {
		log.Printf("Failed to generate certificate, returned error: %s\n", err)
		return nil, err
	}
	cert := Certificate{Cert: caCert, PEM: certPEM, PrivateKey: privateKey}
	return &cert, nil
}

// TODO(Jovan): Duplicate code, remove?
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
	eeCert, certPEM, err := genCert(template, issuerCert, &privateKey.PublicKey, issuerPrivateKey)
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
	return ValidateCertChain(issuerCert)
}

func FindCert(serialNumber string) (*x509.Certificate, error) {
	// TODO(Jovan): Get by serial
	return nil, errors.New("not implemented")
}

func FindCertKey(serialNumber string) (*rsa.PrivateKey, error) {
	// TODO(Jovan): Get by serial
	return nil, errors.New("not implemented")
}

func validateCert(cert *x509.Certificate) error {
	today := time.Now()
	if cert.NotAfter.Before(today) {
		return errors.New("nertificate date is not valid")
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
