//
// Copyright (c) 2022-2023 Winlin
//
// SPDX-License-Identifier: MIT
//
package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
)

var certManager *CertManager

type CertManager struct {
	// httpsCertificate is the certificate for the server.
	httpsCertificate *tls.Certificate

	// httpCertificateReload is used to reload the certificate.
	httpCertificateReload chan bool
}

func NewCertManager() *CertManager {
	return &CertManager{
		httpCertificateReload: make(chan bool, 1),
	}
}

func (v *CertManager) Initialize(ctx context.Context) error {
	return nil
}

func (v *CertManager) ReloadCertificate(ctx context.Context) {
	select {
	case certManager.httpCertificateReload <- true:
	case <-ctx.Done():
	default:
	}
}

func (v *CertManager) reloadCertificateFile(ctx context.Context) error {
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
		logger.Tf(ctx, "crontab: ignore for no cert file")
		return nil
	}

	cert, err := tls.LoadX509KeyPair(crtFile, keyFile)
	if err != nil {
		logger.Tf(ctx, "crontab: ignore load cert %v, key %v failed", crtFile, keyFile)
		return nil
	}

	v.httpsCertificate = &cert
	logger.Tf(ctx, "crontab: reload certificate file ok")

	return nil
}

// updateSslFiles update the ssl files.
func (v *CertManager) updateSslFiles(ctx context.Context, key, crt string) error {
	keyFile := path.Join(conf.Pwd, "containers/data/config/nginx.key")
	crtFile := path.Join(conf.Pwd, "containers/data/config/nginx.crt")

	if err := exec.CommandContext(ctx, "rm", "-f", keyFile, crtFile).Run(); err != nil {
		return errors.Wrapf(err, "rm -f %v %v", keyFile, crtFile)
	}

	if err := ioutil.WriteFile(keyFile, []byte(key), 0644); err != nil {
		return errors.Wrapf(err, "write key %vB to %v", len(key), keyFile)
	}

	if err := ioutil.WriteFile(crtFile, []byte(crt), 0644); err != nil {
		return errors.Wrapf(err, "write crt %vB to %v", len(crt), crtFile)
	}

	return nil
}

// updateLetsEncrypt request letsencrypt and update the ssl files.
func (v *CertManager) updateLetsEncrypt(ctx context.Context, domain string) error {
	defer v.ReloadCertificate(ctx)

	if true {
		args := []string{
			"--email", "srs.stack@gmail.com", "--domains", domain,
			"--http.webroot", path.Join(conf.Pwd, "containers/data"), "--http", "--accept-tos",
			"run",
		}
		cmd := exec.CommandContext(ctx, "lego", args...)
		cmd.Dir = path.Join(conf.Pwd, "containers/data/lego")

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "run lego %v in %v, stdout %v, stderr %v",
				args, cmd.Dir, stdout.String(), stderr.String())
		}
		logger.Tf(ctx, "run lego %v in %v ok, stdout %v, stderr %v",
			args, cmd.Dir, stdout.String(), stderr.String(),
		)
	}

	keyFile := path.Join(conf.Pwd, fmt.Sprintf("containers/data/lego/.lego/certificates/%v.key", domain))
	if _, err := os.Stat(keyFile); err != nil {
		return errors.Wrapf(err, "stat %v", keyFile)
	}

	crtFile := path.Join(conf.Pwd, fmt.Sprintf("containers/data/lego/.lego/certificates/%v.crt", domain))
	if _, err := os.Stat(crtFile); err != nil {
		return errors.Wrapf(err, "stat %v", crtFile)
	}

	targetKeyFile := path.Join(conf.Pwd, "containers/data/config/nginx.key")
	targetCrtFile := path.Join(conf.Pwd, "containers/data/config/nginx.crt")
	if err := exec.CommandContext(ctx, "rm", "-f", targetKeyFile, targetCrtFile).Run(); err != nil {
		return errors.Wrapf(err, "rm -f %v %v", targetKeyFile, targetCrtFile)
	}

	if true {
		source := fmt.Sprintf("../lego/.lego/certificates/%v.key", domain)
		cmd := exec.CommandContext(ctx, "ln", "-sf", source, "nginx.key")
		cmd.Dir = path.Join(conf.Pwd, "containers/data/config")
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "run %v in %v", cmd.Args, cmd.Dir)
		}
	}

	if true {
		source := fmt.Sprintf("../lego/.lego/certificates/%v.crt", domain)
		cmd := exec.CommandContext(ctx, "ln", "-sf", source, "nginx.crt")
		cmd.Dir = path.Join(conf.Pwd, "containers/data/config")
		if err := cmd.Run(); err != nil {
			return errors.Wrapf(err, "run %v in %v", cmd.Args, cmd.Dir)
		}
	}

	return nil
}

// renewLetsEncrypt request letsencrypt and update the ssl files.
func (v *CertManager) renewLetsEncrypt(ctx context.Context, domain string) error {
	defer v.ReloadCertificate(ctx)

	args := []string{
		"--email", "srs.stack@gmail.com", "--domains", domain,
		"--http.webroot", path.Join(conf.Pwd, "containers/data"), "--http",
		"renew", "--days", "30",
	}
	cmd := exec.CommandContext(ctx, "lego", args...)
	cmd.Dir = path.Join(conf.Pwd, "containers/data/lego")

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "run lego %v in %v, stdout %v, stderr %v",
			args, cmd.Dir, stdout.String(), stderr.String())
	}
	logger.Tf(ctx, "run lego %v in %v ok, stdout %v, stderr %v",
		args, cmd.Dir, stdout.String(), stderr.String(),
	)

	return nil
}

func (v *CertManager) refreshSSLCert(ctx context.Context) error {
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

	if err := v.renewLetsEncrypt(ctx, domain); err != nil {
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
