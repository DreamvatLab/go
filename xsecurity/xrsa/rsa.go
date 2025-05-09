package xrsa

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"github.com/DreamvatLab/go/xlog"
)

// GenerateKey generates a new key
func GenerateKey(bits int) (*rsa.PrivateKey, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	return privkey, err
}

// PKCS1PrivateKeyToBytes private key to bytes
func PKCS1PrivateKeyToBytes(priv *rsa.PrivateKey) []byte {
	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	return privBytes
}

// PublicKeyToBytes public key to bytes
func PublicKeyToBytes(pub *rsa.PublicKey) ([]byte, error) {
	pubASN1, err := x509.MarshalPKIXPublicKey(pub)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, err
}

// PKCS1BytesToPrivateKey bytes to private key
func PKCS1BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 检查是否加密
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return nil, fmt.Errorf("encrypted PEM blocks are not supported")
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	return key, err
}

// BytesToPublicKey bytes to public key
func BytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 检查是否加密
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return nil, fmt.Errorf("encrypted PEM blocks are not supported")
	}

	ifc, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	key, ok := ifc.(*rsa.PublicKey)
	if !ok {
		xlog.Error("not ok")
	}
	return key, err
}

// PKCS8PrivateKeyToBytes private key to bytes
func PKCS8PrivateKeyToBytes(priv *rsa.PrivateKey) ([]byte, error) {
	bytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		return nil, err
	}

	privBytes := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: bytes,
		},
	)

	return privBytes, err
}

// PKCS8BytesToPrivateKey bytes to private key
func PKCS8BytesToPrivateKey(priv []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// 检查是否加密
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return nil, fmt.Errorf("encrypted PEM blocks are not supported")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return key.(*rsa.PrivateKey), err
}

// CertificateBytesToPublicKey bytes to public key
func CertificateBytesToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Check if encrypted
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		return nil, fmt.Errorf("encrypted PEM blocks are not supported")
	}

	ifc, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}

	return ifc.PublicKey.(*rsa.PublicKey), err
}

// EncryptWithPublicKey encrypts data with public key
func EncryptWithPublicKey(msg []byte, pub *rsa.PublicKey) ([]byte, error) {
	hash := sha512.New()
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, pub, msg, nil)
	if err != nil {
		return nil, err
	}
	return ciphertext, err
}

// DecryptWithPrivateKey decrypts data with private key
func DecryptWithPrivateKey(ciphertext []byte, priv *rsa.PrivateKey) ([]byte, error) {
	hash := sha512.New()
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, err
}

func ReadCertFromFile(file string) (*x509.Certificate, error) {
	caFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	caBlock, _ := pem.Decode(caFile)
	cert, err := x509.ParseCertificate(caBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return cert, err
}

func ReadPrivateKeyFromFile(file string) (*rsa.PrivateKey, error) {
	keyFile, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	keyBlock, _ := pem.Decode(keyFile)
	praKey, err := x509.ParsePKCS8PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	return praKey.(*rsa.PrivateKey), nil
}
