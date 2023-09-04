package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	cr "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestApi_SetupWebsiteFooter(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	req := struct {
		Beian string `json:"beian"`
		Text  string `json:"text"`
	}{
		Beian: "icp", Text: "TestFooter",
	}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/update", &req, nil); err != nil {
		r0 = err
		return
	}

	res := struct {
		ICP string `json:"icp"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/query", nil, &res); err != nil {
		r0 = err
	} else if res.ICP != "TestFooter" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_SetupWebsiteTitle(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var title string
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/query", nil, &struct {
		Title *string `json:"title"`
	}{
		Title: &title,
	}); err != nil {
		r0 = err
		return
	} else if title == "" {
		title = "SRS"
	}
	defer func() {
		if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/update", &struct {
			Beian string `json:"beian"`
			Text  string `json:"text"`
		}{
			Beian: "title", Text: title,
		}, nil); err != nil {
			r0 = err
		}
	}()

	req := struct {
		Beian string `json:"beian"`
		Text  string `json:"text"`
	}{
		Beian: "title", Text: "TestTitle",
	}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/update", &req, nil); err != nil {
		r0 = err
		return
	}

	res := struct {
		Title string `json:"title"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/beian/query", nil, &res); err != nil {
		r0 = err
	} else if res.Title != "TestTitle" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

// Never run this in parallel, because it changes the publish
// secret which might cause other cases to fail.
func TestApi_UpdatePublishSecret(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0, r1 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var pubSecret string
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &pubSecret,
	}); err != nil {
		r0 = err
		return
	} else if pubSecret == "" {
		r0 = errors.Errorf("invalid response %v", pubSecret)
		return
	}

	// Reset the publish secret to the original value.
	defer func() {
		logger.Tf(ctx, "Reset publish secret to %v", pubSecret)
		if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/update", &struct {
			Secret string `json:"secret"`
		}{
			Secret: pubSecret,
		}, nil); err != nil {
			r1 = err
		}
	}()

	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/update", &struct {
		Secret string `json:"secret"`
	}{
		Secret: "TestPublish",
	}, nil); err != nil {
		r0 = err
		return
	}

	res := struct {
		Publish string `json:"publish"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &res); err != nil {
		r0 = err
	} else if res.Publish != "TestPublish" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_TutorialsQueryBilibili(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// If we are using letsencrypt, we don't need to test this.
	if *domainLetsEncrypt != "" || *httpsInsecureVerify {
		return
	}

	if *noBilibiliTest {
		return
	}

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	req := struct {
		BVID string `json:"bvid"`
	}{
		BVID: "BV1844y1L7dL",
	}
	res := struct {
		Title string `json:"title"`
		Desc  string `json:"desc"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/bilibili", &req, &res); err != nil {
		r0 = err
	} else if res.Title == "" || res.Desc == "" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_SslUpdateCert(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// If we are using letsencrypt, we don't need to test this.
	if *domainLetsEncrypt != "" || *httpsInsecureVerify {
		return
	}

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	var key, crt string
	if err := func() error {
		privateKey, err := ecdsa.GenerateKey(elliptic.P256(), cr.Reader)
		if err != nil {
			return errors.Wrapf(err, "generate ecdsa key")
		}

		template := x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				CommonName: "srs.stack.local",
			},
			NotBefore: time.Now(),
			NotAfter:  time.Now().AddDate(10, 0, 0),
			KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
			},
			BasicConstraintsValid: true,
		}

		derBytes, err := x509.CreateCertificate(cr.Reader, &template, &template, &privateKey.PublicKey, privateKey)
		if err != nil {
			return errors.Wrapf(err, "create certificate")
		}

		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return errors.Wrapf(err, "marshal ecdsa key")
		}

		privateKeyBlock := pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		}
		key = string(pem.EncodeToMemory(&privateKeyBlock))

		certBlock := pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		}
		crt = string(pem.EncodeToMemory(&certBlock))
		logger.Tf(ctx, "cert: create self-signed certificate ok, key=%vB, crt=%vB", len(key), len(crt))
		return nil
	}(); err != nil {
		r0 = err
		return
	}

	if err := apiRequest(ctx, "/terraform/v1/mgmt/ssl", &struct {
		Key string `json:"key"`
		Crt string `json:"crt"`
	}{
		Key: key, Crt: crt,
	}, nil); err != nil {
		r0 = err
		return
	}

	conf := struct {
		Provider string `json:"provider"`
		Key      string `json:"key"`
		Crt      string `json:"crt"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/cert/query", nil, &conf); err != nil {
		r0 = err
	} else if conf.Provider != "ssl" || conf.Key != key || conf.Crt != crt {
		r0 = errors.Errorf("invalid response %v", conf)
	}
}

func TestApi_LetsEncryptUpdateCert(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if *domainLetsEncrypt == "" {
		return
	}

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	if err := apiRequest(ctx, "/terraform/v1/mgmt/letsencrypt", &struct {
		Domain string `json:"domain"`
	}{
		Domain: *domainLetsEncrypt,
	}, nil); err != nil {
		r0 = err
		return
	}

	conf := struct {
		Provider string `json:"provider"`
		Key      string `json:"key"`
		Crt      string `json:"crt"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/cert/query", nil, &conf); err != nil {
		r0 = err
	} else if conf.Provider != "lets" || conf.Key == "" || conf.Crt == "" {
		r0 = errors.Errorf("invalid response %v", conf)
	}
}

func TestApi_SetupHpHLSNoHlsCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	type Data struct {
		NoHlsCtx bool `json:"noHlsCtx"`
	}

	if true {
		initData := Data{}
		if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/query", nil, &initData); err != nil {
			r0 = err
			return
		}
		defer func() {
			if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/update", &initData, nil); err != nil {
				logger.Tf(ctx, "restore hphls config failed %+v", err)
			}
		}()
	}

	noHlsCtx := Data{NoHlsCtx: true}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/update", &noHlsCtx, nil); err != nil {
		r0 = err
		return
	}

	verifyData := Data{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/query", nil, &verifyData); err != nil {
		r0 = err
		return
	} else if verifyData.NoHlsCtx != true {
		r0 = errors.Errorf("invalid response %+v", verifyData)
	}
}

func TestApi_SetupHpHLSWithHlsCtx(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0); err != nil {
			t.Errorf("Fail for err %+v", err)
		} else {
			logger.Tf(ctx, "test done")
		}
	}(ctx)

	type Data struct {
		NoHlsCtx bool `json:"noHlsCtx"`
	}

	if true {
		initData := Data{}
		if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/query", nil, &initData); err != nil {
			r0 = err
			return
		}
		defer func() {
			if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/update", &initData, nil); err != nil {
				logger.Tf(ctx, "restore hphls config failed %+v", err)
			}
		}()
	}

	noHlsCtx := Data{NoHlsCtx: false}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/update", &noHlsCtx, nil); err != nil {
		r0 = err
		return
	}

	verifyData := Data{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/hphls/query", nil, &verifyData); err != nil {
		r0 = err
		return
	} else if verifyData.NoHlsCtx != false {
		r0 = errors.Errorf("invalid response %+v", verifyData)
	}
}
