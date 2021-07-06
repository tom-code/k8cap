package main

import (
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "crypto/x509/pkix"
  "encoding/base64"
  "encoding/pem"
  "fmt"
  "io/ioutil"
  "math/big"
  "time"
  )


func main() {
  priv, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    panic(err)
  }

  privPKCS1 := x509.MarshalPKCS1PrivateKey(priv)
  privBlock := pem.Block {
    Type: "RSA PRIVATE KEY",
    Bytes: privPKCS1,
  }
  err = ioutil.WriteFile("private.pem", pem.EncodeToMemory(&privBlock), 0600)
  if err != nil {
    panic(err)
  }

  pubPKIX, err := x509.MarshalPKIXPublicKey(&priv.PublicKey)
  if err != nil {
    panic(err)
  }
  pubBlock := pem.Block {
    Type: "PUBLIC KEY",
    Bytes: pubPKIX,
  }
  err = ioutil.WriteFile("public.pem", pem.EncodeToMemory(&pubBlock), 0600)
  if err != nil {
    panic(err)
  }

  template := x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name {
      Organization: []string{"owner"},
    },
    DNSNames: []string{"hook1.default.svc"},
    NotBefore: time.Now(),
    NotAfter:  time.Now().Add(time.Hour * 24 * 700),

    KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
    BasicConstraintsValid: true,
  }

  certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
  if err != nil {
    panic(err)
  }
  certBlock := pem.Block {
    Type: "CERTIFICATE",
    Bytes: certBytes,
  }
  err = ioutil.WriteFile("cert.pem", pem.EncodeToMemory(&certBlock), 0600)
  if err != nil {
    panic(err)
  }
  pubb64 := base64.StdEncoding.EncodeToString(pem.EncodeToMemory(&certBlock))
  fmt.Printf("caBundle: %s\n", pubb64)
}
