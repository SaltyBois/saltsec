package cert

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"math/big"
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

func GenCAIntermediateCert(template, parentCert *x509.Certificate, parentPrivateKey *rsa.PrivateKey) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key, returned error: %s\n", err)
	}
	cert, certPEM := genCert(template, parentCert, &parentPrivateKey.PublicKey, parentPrivateKey)
	return cert, certPEM, privateKey
}

// TODO(Jovan): Duplicate code, remove?
func GenEndEntityCert(template, parentCert *x509.Certificate, parentPrivateKey *rsa.PrivateKey) (*x509.Certificate, []byte, *rsa.PrivateKey) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate key, returned error: %s\n", err)
	}
	cert, certPEM := genCert(template, parentCert, &parentPrivateKey.PublicKey, parentPrivateKey)
	return cert, certPEM, privateKey
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
