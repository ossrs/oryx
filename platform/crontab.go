//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"os"
	"path"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

var crontabWorker *CrontabWorker

type CrontabWorker struct {
	wg sync.WaitGroup
}

func NewCrontabWorker() *CrontabWorker {
	return &CrontabWorker{}
}

func (v *CrontabWorker) Close() error {
	v.wg.Wait()
	return nil
}

func (v *CrontabWorker) Start(ctx context.Context) error {
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(24*3600) * time.Second):
			}

			logger.Tf(ctx, "crontab: start to query latest version")
			if versions, err := queryLatestVersion(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			} else {
				logger.Tf(ctx, "crontab: query version ok, result is %v", versions.String())
			}
		}
	}()

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			logger.Tf(ctx, "crontab: start to refresh ssl cert")
			if err := refreshSSLCert(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Duration(24*3600) * time.Second):
			}
		}
	}()

	httpCertificateReload = make(chan bool, 1)

	if cert, err := generateSelfSignCertificate(ctx); err != nil {
		return errors.Wrapf(err, "generate self-signed certificate")
	} else {
		selfSignedCertificate = cert
	}

	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		for {
			logger.Tf(ctx, "crontab: start to refresh certificate file")
			if err := reloadCertificateFile(ctx); err != nil {
				logger.Wf(ctx, "crontab: ignore err %v", err)
			}

			select {
			case <-ctx.Done():
				return
			case <-httpCertificateReload:
			case <-time.After(time.Duration(1*3600) * time.Second):
			}
		}
	}()

	return nil
}

// httpsCertificate is the certificate for the server.
var httpsCertificate *tls.Certificate

// selfSignedCertificate is the self-signed certificate for the server.
var selfSignedCertificate *tls.Certificate

// httpCertificateReload is used to reload the certificate.
var httpCertificateReload chan bool

func reloadCertificateFile(ctx context.Context) error {
	keyFile := path.Join(conf.Pwd, "containers/data/config/nginx.key")
	crtFile := path.Join(conf.Pwd, "containers/data/config/nginx.crt")

	var noKeyFile, noCrtFile bool
	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
		noKeyFile = true
	}
	if _, err := os.Stat(crtFile); os.IsNotExist(err) {
		noCrtFile = true
	}
	if noKeyFile || noCrtFile {
		httpsCertificate = selfSignedCertificate
		logger.Tf(ctx, "crontab: use self-signed certificate")
		return nil
	}

	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		httpsCertificate = selfSignedCertificate
		return errors.Wrapf(err, "crontab: load cert %v, key %v failed, use self-signed certificate", crtFile, keyFile)
	}

	httpsCertificate = &cert
	logger.Tf(ctx, "crontab: reload certificate file ok")

	return nil
}

func generateSelfSignCertificate(ctx context.Context) (*tls.Certificate, error) {
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, errors.Wrapf(err, "generate key")
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"srs.stack.local"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return nil, errors.Wrapf(err, "create cert")
	}

	privBytes, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return nil, errors.Wrapf(err, "marshal ec private key")
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, errors.Wrapf(err, "x509 key pair")
	}

	return &cert, nil
}

func refreshSSLCert(ctx context.Context) error {
	provider, err := rdb.Get(ctx, SRS_HTTPS).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if provider != "lets" {
		logger.Tf(ctx, "crontab: ignore ssl provider %v", provider)
		return nil
	}

	domain, err := rdb.Get(ctx, SRS_HTTPS_DOMAIN).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if domain == "" {
		logger.Tf(ctx, "crontab: ignore ssl domain empty")
		return nil
	}

	if err := renewLetsEncrypt(ctx, domain); err != nil {
		return err
	} else {
		logger.Tf(ctx, "crontab: renew ssl cert ok")
	}

	if err := nginxGenerateConfig(ctx); err != nil {
		return errors.Wrapf(err, "nginx config and reload")
	}

	logger.Tf(ctx, "crontab: refresh ssl cert ok")
	return nil
}
