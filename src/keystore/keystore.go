package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"io/ioutil"
	"log"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)


func ReadPFX(filename string, password string) (*rsa.PrivateKey ,*x509.Certificate) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Faild to load PFX file")
	}
	privateKey, cert, err := pkcs12.Decode(data, password)
	if err != nil {
		log.Printf("Faild decoding PFX data")
	}

	PKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil
	}
	return PKey, cert
}


func WritePFX(cert *x509.Certificate, PrivateKey *rsa.PrivateKey, password string, filepath string) error {
	filename := filepath + cert.SerialNumber.String() + ".pfx"
	pfxBytes, err := pkcs12.Encode(rand.Reader, PrivateKey, cert, []*x509.Certificate{}, password)
	if err != nil {
		log.Printf("Faild to encode PFX file\n")
		return err
	}
	if err := os.WriteFile(
		filename,
		pfxBytes,
		os.ModePerm,
	); err != nil {
		log.Printf("Faild writing to PFX file, returned error: #{err}\n")
		return err
	}
	return nil
}

