package keystore

import (
	"log"
	"os"
	"github.com/pavel-v-chernykh/keystore-go/v4"
)

func readKeyStore(filename string, password []byte) keystore.KeyStore {
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	keyStore := keystore.New()
	if err := keyStore.Load(f, password); err != nil {
		log.Fatal(err) // nolint: gocritic
	}

	return keyStore
}

func writeKeyStore(keyStore keystore.KeyStore, filename string, password []byte) {
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	err = keyStore.Store(f, password)
	if err != nil {
		log.Fatal(err) // nolint: gocritic
	}
}

