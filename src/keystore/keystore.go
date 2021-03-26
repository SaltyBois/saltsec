package keystore

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"software.sslmate.com/src/go-pkcs12"
)

func ReadPFX(filename string, password string) (*rsa.PrivateKey, *x509.Certificate, error) {
	fmt.Println(filename)
	data, err := ioutil.ReadFile(filepath.FromSlash(filename))
	if err != nil {
		log.Println("Faild to load PFX file: ", err)
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

func WritePFX(cert *x509.Certificate, PrivateKey *rsa.PrivateKey, password string, filename string) error {
	pfxBytes, err := pkcs12.Encode(rand.Reader, PrivateKey, cert, []*x509.Certificate{}, password)
	if err != nil {
		log.Printf("Faild to encode PFX file\n")
		return err
	}
	log.Println(filename)
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
