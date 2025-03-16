// Package certs пакет для создания сертификатов для https
package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"

	"math/big"
	"net"
	"os"
	"time"
)

// CertPair - содержит пару ключей для ассиметричного шифрования
type CertPair struct {
	CertPem       bytes.Buffer
	PrivateKeyPEM bytes.Buffer
}

// CertPairFiles - содержит пути к файлам с парой ключей для ассиметричного шифрования
type CertPairFiles struct {
	CertPem       string
	PrivateKeyPEM string
}

// Generate - генерация пары ключей
func Generate() (CertPair, error) {
	// создаём шаблон сертификата
	result := CertPair{}
	cert := &x509.Certificate{
		// указываем уникальный номер сертификата
		SerialNumber: big.NewInt(1658),
		// заполняем базовую информацию о владельце сертификата
		Subject: pkix.Name{
			Organization: []string{"Yandex.Praktikum"},
			Country:      []string{"RU"},
		},
		// разрешаем использование сертификата для 127.0.0.1 и ::1
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		// сертификат верен, начиная со времени создания
		NotBefore: time.Now(),
		// время жизни сертификата — 10 лет
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		// устанавливаем использование ключа для цифровой подписи,
		// а также клиентской и серверной авторизации
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature,
	}

	// создаём новый приватный RSA-ключ длиной 4096 бит
	// обратите внимание, что для генерации ключа и сертификата
	// используется rand.Reader в качестве источника случайных данных
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return result, err
	}

	// создаём сертификат x.509
	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &privateKey.PublicKey, privateKey)
	if err != nil {
		return result, err
	}

	// кодируем сертификат и ключ в формате PEM, который
	// используется для хранения и обмена криптографическими ключами

	pem.Encode(&result.CertPem, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	pem.Encode(&result.PrivateKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	return result, nil
}

// GetCertFiles - генерация и сохранение в файлы пары ключей
func GetCertFiles() (CertPairFiles, error) {
	files := CertPairFiles{}
	certPair, err := Generate()

	if err != nil {
		return files, err
	}

	certPemFile, err := os.CreateTemp(os.TempDir(), "*")

	if err != nil {
		return files, err
	}

	_, err = certPemFile.Write(certPair.CertPem.Bytes())
	if err != nil {
		return files, err
	}
	certPemFile.Close()
	privKeyFile, err := os.CreateTemp(os.TempDir(), "*")

	if err != nil {
		return files, err
	}

	_, err = privKeyFile.Write(certPair.PrivateKeyPEM.Bytes())
	privKeyFile.Close()
	if err != nil {
		return files, err
	}

	files.CertPem = certPemFile.Name()
	files.PrivateKeyPEM = privKeyFile.Name()

	return files, nil
}
