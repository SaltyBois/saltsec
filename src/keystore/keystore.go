package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"log"
	"os"

	"software.sslmate.com/src/go-pkcs12"
)

func ReadPFX(filename string, password string) (*rsa.PrivateKey, *x509.Certificate, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("Faild to load PFX file")
		return nil, nil, err
	}
	privateKey, cert, err := pkcs12.Decode(data, password)
	if err != nil {
		log.Printf("Faild decoding PFX data")
		return nil, nil, err
	}

	PKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, nil, errors.New("could not convert to rsa.PrivateKey")
	}
	return PKey, cert, nil
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
