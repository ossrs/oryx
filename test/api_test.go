package main

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	cr "crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

func TestApi_Empty(t *testing.T) {
	ctx := logger.WithContext(context.Background())
	logger.Tf(ctx, "test done")
}

func TestApi_Ready(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	if err := waitForServiceReady(ctx); err != nil {
		t.Errorf("Fail for err %+v", err)
	} else {
		logger.Tf(ctx, "test done")
	}
}

func TestApi_QueryPublishSecret(t *testing.T) {
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

	var publishSecret string
	if err := apiRequest(ctx, "/terraform/v1/hooks/srs/secret/query", nil, &struct {
		Publish *string `json:"publish"`
	}{
		Publish: &publishSecret,
	}); err != nil {
		r0 = err
	} else if publishSecret == "" {
		r0 = errors.New("empty publish secret")
	}
}

func TestApi_LoginByPassword(t *testing.T) {
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

	var token string
	if err := apiRequest(ctx, "/terraform/v1/mgmt/login", &struct {
		Password string `json:"password"`
	}{
		Password: *systemPassword,
	}, &struct {
		Token *string `json:"token"`
	}{
		Token: &token,
	}); err != nil {
		r0 = err
	} else if token == "" {
		r0 = errors.New("empty token")
	}
}

func TestApi_BootstrapQueryEnvs(t *testing.T) {
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

	res := struct {
		MgmtDocker bool `json:"mgmtDocker"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/envs", nil, &res); err != nil {
		r0 = err
	} else if !res.MgmtDocker {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryInit(t *testing.T) {
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

	res := struct {
		Init bool `json:"init"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/init", nil, &res); err != nil {
		r0 = err
	} else if !res.Init {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryCheck(t *testing.T) {
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

	res := struct {
		Upgrading bool `json:"upgrading"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/check", nil, &res); err != nil {
		r0 = err
	} else if res.Upgrading {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

func TestApi_BootstrapQueryVersions(t *testing.T) {
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

	res := struct {
		Version string `json:"version"`
	}{}
	if err := apiRequest(ctx, "/terraform/v1/mgmt/versions", nil, &res); err != nil {
		r0 = err
	} else if res.Version == "" {
		r0 = errors.Errorf("invalid response %v", res)
	}
}

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

func TestApi_PublishRtmpPlayFlv_SecretQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
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
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("rtmp://localhost/live/%v?secret=%v", streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re",
			"-f", "lavfi", "-i", "testsrc=size=1280x720", "-f", "lavfi", "-i", "sine=frequency=440",
			"-pix_fmt", "yuv420p", "-vcodec", "libx264", "-profile:v", "baseline", "-r", "25", "-g", "50",
			"-acodec", "aac", "-ar", "44100", "-ac", "2",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("http://localhost:8080/live/%v.flv", streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	if ts := 90; m.Format.ProbeScore < ts {
		r4 = errors.Errorf("low score=%v < %v, %v, %v", m.Format.ProbeScore, ts, m.String(), str)
	}
	if dv := m.Duration(); dv < duration {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestApi_PublishRtmpPlayFlv_SecretStream(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
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
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v-%v", pubSecret, os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("rtmp://localhost/live/%v", streamID)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re",
			"-f", "lavfi", "-i", "testsrc=size=1280x720", "-f", "lavfi", "-i", "sine=frequency=440",
			"-pix_fmt", "yuv420p", "-vcodec", "libx264", "-profile:v", "baseline", "-r", "25", "-g", "50",
			"-acodec", "aac", "-ar", "44100", "-ac", "2",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("http://localhost:8080/live/%v.flv", streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	if ts := 90; m.Format.ProbeScore < ts {
		r4 = errors.Errorf("low score=%v < %v, %v, %v", m.Format.ProbeScore, ts, m.String(), str)
	}
	if dv := m.Duration(); dv < duration {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}

func TestApi_PublishRtmpPlayHls_SecretQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
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
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("rtmp://localhost/live/%v?secret=%v", streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re",
			"-f", "lavfi", "-i", "testsrc=size=1280x720", "-f", "lavfi", "-i", "sine=frequency=440",
			"-pix_fmt", "yuv420p", "-vcodec", "libx264", "-profile:v", "baseline", "-r", "25", "-g", "50",
			"-acodec", "aac", "-ar", "44100", "-ac", "2",
			"-f", "flv", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("http://localhost:8080/live/%v.m3u8", streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	// Note that HLS score is low, so we only check duration. Note that only check half of duration, because we
	// might get only some pieces of segments.
	if dv := m.Duration(); dv < duration/2 {
		r4 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration/2, m.String(), str)
	}
}

func TestApi_PublishSrtPlayFlv_SecretQuery(t *testing.T) {
	ctx, cancel := context.WithTimeout(logger.WithContext(context.Background()), time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	var r0, r1, r2, r3, r4, r5 error
	defer func(ctx context.Context) {
		if err := filterTestError(ctx.Err(), r0, r1, r2, r3, r4, r5); err != nil {
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
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	// Start FFmpeg to publish stream.
	streamID := fmt.Sprintf("stream-%v-%v", os.Getpid(), rand.Int())
	streamURL := fmt.Sprintf("srt://localhost:10080#!::r=live/%v?secret=%v,m=publish", streamID, pubSecret)
	ffmpeg := NewFFmpeg(func(v *ffmpegClient) {
		v.args = []string{
			"-re",
			"-f", "lavfi", "-i", "testsrc=size=1280x720", "-f", "lavfi", "-i", "sine=frequency=440",
			"-pix_fmt", "yuv420p", "-vcodec", "libx264", "-profile:v", "baseline", "-r", "25", "-g", "50",
			"-acodec", "aac", "-ar", "44100", "-ac", "2",
			"-f", "mpegts", streamURL,
		}
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r1 = ffmpeg.Run(ctx, cancel)
	}()

	// Start FFprobe to detect and verify stream.
	duration := time.Duration(*srsFFprobeDuration) * time.Millisecond
	ffprobe := NewFFprobe(func(v *ffprobeClient) {
		v.dvrFile = fmt.Sprintf("srs-ffprobe-%v.flv", streamID)
		v.streamURL = fmt.Sprintf("http://localhost:8080/live/%v.flv", streamID)
		v.duration, v.timeout = duration, time.Duration(*srsFFprobeTimeout)*time.Millisecond
	})
	wg.Add(1)
	go func() {
		defer wg.Done()
		r2 = ffprobe.Run(ctx, cancel)
	}()

	// Fast quit for probe done.
	select {
	case <-ctx.Done():
	case <-ffprobe.ProbeDoneCtx().Done():
		cancel()
	}

	str, m := ffprobe.Result()
	if len(m.Streams) != 2 {
		r3 = errors.Errorf("invalid streams=%v, %v, %v", len(m.Streams), m.String(), str)
	}

	if ts := 90; m.Format.ProbeScore < ts {
		r4 = errors.Errorf("low score=%v < %v, %v, %v", m.Format.ProbeScore, ts, m.String(), str)
	}
	if dv := m.Duration(); dv < duration {
		r5 = errors.Errorf("short duration=%v < %v, %v, %v", dv, duration, m.String(), str)
	}
}
