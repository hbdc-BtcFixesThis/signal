package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"net"
	"os"
	"strings"
	"time"
)

var (
	validFrom  = ""
	validFor   = 365 * 24 * time.Hour
	isCA       = true
	rsaBits    = 2048
	ecdsaCurve = RSA
)

type EllipticCurve uint8

const (
	RSA EllipticCurve = iota
	P224
	P256
	P384
	P521
)

func (ec EllipticCurve) String() string { return [...]string{"", "P224", "P256", "P384", "P521"}[ec] }

func (ec EllipticCurve) genKey() (interface{}, error) {
	switch ec {
	case RSA:
		return rsa.GenerateKey(rand.Reader, rsaBits)
	case P224:
		return ecdsa.GenerateKey(elliptic.P224(), rand.Reader)
	case P256:
		return ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case P384:
		return ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case P521:
		return ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		return nil, fmt.Errorf("Unrecognized elliptic curve: %s", ec.String())
	}
}

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

func CheckTLSKeyCertPath(certPath string, keyPath string) error {
	if _, err := os.Stat(certPath); os.IsNotExist(err) {
		return err
	} else if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return err
	}
	return nil
}

func notBeforeAfter() (time.Time, time.Time, error) {
	var notBefore time.Time
	var err error
	if len(validFrom) == 0 {
		notBefore = time.Now()
	} else {
		notBefore, err = time.Parse("Jan 3 00:00:00 2009", validFrom)
		if err != nil {
			// err will be picked up and returned below
			return notBefore, notBefore, err
		}
	}

	return notBefore, notBefore.Add(validFor), nil
}

func serialNumber() (*big.Int, error) {
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	return rand.Int(rand.Reader, serialNumberLimit)
}

func save(path string, pemBlock *pem.Block) error {
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Printf("failed to open %s for writing: %s", path, err)
		return err
	}
	pem.Encode(out, pemBlock)
	out.Close()
	log.Printf("written %s\n", path)
	return nil
}

func GenerateTLSKeyCert(certPath string, keyPath string, host string) error {
	priv, err := ecdsaCurve.genKey()
	if err != nil {
		log.Printf("failed to generate private key: %s", err)
		return err
	}

	notBefore, notAfter, err := notBeforeAfter()
	if err != nil {
		log.Printf("Failed to parse creation date: %s\n", err)
		return err
	}

	serialNumber, err := serialNumber()
	if err != nil {
		log.Printf("failed to generate serial number: %s", err)
		return err
	}

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               pkix.Name{Organization: []string{"Bitcoin Signal"}},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	hosts := strings.Split(host, ",")
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	if isCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		log.Printf("Failed to create certificate: %s", err)
		return err
	}

	err = save(keyPath, pemBlockForKey(priv))
	if err != nil {
		return err
	}
	return save(certPath, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
}
