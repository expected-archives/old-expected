package certs

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/docker/libtrust"
	"github.com/sirupsen/logrus"
)

type Config struct {
	PublicKey  string `envconfig:"public_key" default:"./certs/server.crt"`
	PrivateKey string `envconfig:"private_key" default:"./certs/server.key"`
}

var (
	publicKey  libtrust.PublicKey
	privateKey libtrust.PrivateKey
)

func Init(config Config) error {
	var err error

	publicKey, privateKey, err = loadCertAndKey(config.PublicKey, config.PrivateKey)
	if err != nil {
		logrus.
			WithField("init certs", err).
			WithField("publicKey", publicKey).
			WithField("privateKey", privateKey).
			Error()
		return err
	}
	return nil
}

func loadCertAndKey(certFile, keyFile string) (pk libtrust.PublicKey, prk libtrust.PrivateKey, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, nil, err
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return nil, nil, err
	}
	pk, err = libtrust.FromCryptoPublicKey(x509Cert.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	prk, err = libtrust.FromCryptoPrivateKey(cert.PrivateKey)
	return pk, prk, nil
}

func GetPublicKey() libtrust.PublicKey {
	return publicKey
}

func GetPrivateKey() libtrust.PrivateKey {
	return privateKey
}
