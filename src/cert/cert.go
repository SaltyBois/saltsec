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

func GetRandomSerial() *big.Int {
	z := new(big.Int)
	b, err := genRandomBytes(256)
	if err != nil {
		log.Fatalf("Failed to generate random serial, returned error: %s\n", err)
	}
	z.SetBytes(b)
	return z
}

func GenCARootCert(rootTemplate *x509.Certificate) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key, returned error: %s\n", err)
	}
	// NOTE(Jovan): Generate a self-signed cert
	rootCert, rootPEM := genCert(rootTemplate, rootTemplate, &privateKey.PublicKey, privateKey)
	return rootCert, rootPEM, privateKey
}

func GenCAIntermediateCert(template *x509.Certificate, issuerSerialNumber string) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	issuerCert, err := FindCert(issuerSerialNumber)
	if err != nil {
		// TODO(Jovan): Handle error
		return nil, nil, nil
	}

	issuerPrivateKey, err := FindCertKey(issuerSerialNumber)
	if err != nil {
		// TODO(Jovan): Handle error
		return nil, nil, nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key, returned error: %s\n", err)
	}
	cert, certPEM := genCert(template, issuerCert, &issuerPrivateKey.PublicKey, issuerPrivateKey)
	return cert, certPEM, privateKey
}

// TODO(Jovan): Duplicate code, remove?
func GenEndEntityCert(template *x509.Certificate, issuerSerialNumber string) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	issuerCert, err := FindCert(issuerSerialNumber)
	if err != nil {
		// TODO(Jovan): Handle error
		return nil, nil, nil
	}

	issuerPrivateKey, err := FindCertKey(issuerSerialNumber)
	if err != nil {
		// TODO(Jovan): Handle error
		return nil, nil, nil
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key, returned error: %s\n", err)
	}
	cert, certPEM := genCert(template, issuerCert, &issuerPrivateKey.PublicKey, issuerPrivateKey)
	return cert, certPEM, privateKey
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
		log.Fatalf("Error generating random bytes, length less than 0\n")
	}
	b := make([]byte, length)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func genCert(template, parent *x509.Certificate, publicKey *rsa.PublicKey, privateKey *rsa.PrivateKey) (*x509.Certificate, []byte) {
	certBytes, err := x509.CreateCertificate(rand.Reader, template, parent, publicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate, returned error: %s\n", err)
	}
	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		log.Fatalf("Failed to parse certificate, returned error: %s\n", err)
	}

	b := pem.Block{Type: "CERTIFICATE", Bytes: certBytes}
	certPEM := pem.EncodeToMemory(&b)

	return cert, certPEM
}
