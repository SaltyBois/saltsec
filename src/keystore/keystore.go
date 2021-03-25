package keystore

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"

	"saltsec/cert"

	"software.sslmate.com/src/go-pkcs12"
)

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func encrypt(data []byte, password string) []byte {
	block, _ := aes.NewCipher([]byte(createHash(password)))
	gcm, _ := cipher.NewGCM(block)
	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)
	cipher := gcm.Seal(nonce, nonce, data, nil)
	return cipher
}

func ReadPFX(filename string) (*rsa.PrivateKey ,*x509.Certificate) {
	data, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Printf("Faild to load PFX file")
	}

	//Todo change password
	privateKey, cert, err := pkcs12.Decode(data, pkcs12.DefaultPassword)
	if err != nil {
		log.Printf("Faild decoding PFX data")
	}

	return privateKey, cert

}


func WritePFX(cert *x509.Certificate, PrivateKey *rsa.PrivateKey) {

	dir := "../keys/rooot/"
	file := path.Join(dir, "pri.pfx")

	// Todo change password
	pfxBytes, err := pkcs12.Encode(rand.Reader, PrivateKey, cert, []*x509.Certificate{}, pkcs12.DefaultPassword)

	if err != nil {
		log.Printf("Faild writing to PFX file, returned error: #{err}\n")
	}

	if _, _, err := pkcs12.Decode(pfxBytes, pkcs12.DefaultPassword); err != nil {
		log.Printf("Error decoding: #{err}\n")
	}

	if err := ioutil.WriteFile(
		file,
		pfxBytes,
		os.ModePerm,
	); err != nil {
		log.Printf("Error decoding: #{err}\n")
	}

}

