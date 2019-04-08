package token

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/docker/libtrust"
	"github.com/sirupsen/logrus"
	"os"
)

var publicKey libtrust.PublicKey
var privateKey libtrust.PrivateKey

func Init(publicK string, privateK string) {
	var err error

	publicKey, privateKey, err = loadCertAndKey(publicK, privateK)
	if err != nil {
		dir, _ := os.Getwd()
		logrus.
			WithField("init certs", err).
			WithField("cwd", dir).
			WithField("publicKey", publicKey).
			WithField("privateKey", privateKey).
			Fatalln()
	}
}

func loadCertAndKey(certFile, keyFile string) (pk libtrust.PublicKey, prk libtrust.PrivateKey, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}
	pk, err = libtrust.FromCryptoPublicKey(x509Cert.PublicKey)
	if err != nil {
		return
	}
	prk, err = libtrust.FromCryptoPrivateKey(cert.PrivateKey)
	return
}
