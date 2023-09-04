package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

// initPassword is used to init the password for SRS Stack.
//
// SRS Stack requires initializing and setting up a password by default when the system is initialized.
// This password is utilized for logging into the system through the UI and API. Typically, for testing
// purposes, we generate a random password and proceed with system initialization. However, when developing
// and debugging tests for an already initialized system, there is no need to initialize a new password,
// as it has already been established.
var initPassword *bool

// systemPassword is the password sued by SRS Stack.
//
// Users can specify a dedicated password; if not provided, a temporary password will be created to
// initialize the SRS Stack. This password will be used to test the login feature.
var systemPassword *string

// initSelfSignedCert is used to init the self-signed cert for HTTPS.
//
// We should not initialize self-signed certificates by default, as this sets up SSL files for nginx,
// making it impossible for BT or aaPanel to configure SSL. Therefore, we disable it by default, and users
// can enable it if necessary, such as for testing the HTTPS API. It is important to note that we support
// requesting SSL certificates through BT or aaPanel, or via SRS Stack, all of which share the same
// .well-known web directory and nginx configurations. However, SRS Stack does not rely on the NGINX SSL
// files.
//
// When this feature is enabled and there are no certificate files (nginx.key and nginx.crt) in the nginx
// directory, SRS Stack will create a self-signed certificate, save it to the nginx SSL file, and then
// generate the nginx configuration in nginx.server.conf.
var initSelfSignedCert *bool

// domainLetsEncrypt is the domain used to test the letsencrypt feature.
//
// In a development environment, you typically don't have a public IP address, which means you won't be
// able to test the Let's Encrypt feature by default. However, when testing in an online environment, you
// can enable this feature, and it will initiate a genuine request to Let's Encrypt to obtain a valid HTTPS
// certificate.
//
// If empty string, don't test the letsencrypt feature.
var domainLetsEncrypt *string

// httpsInsecureVerify is used to verify the HTTPS certificate.
//
// When using a self-signed certificate for local HTTPS testing, you should not attempt to verify the HTTP
// certificate, as it is self-signed and verification will fail. For testing in an online environment, you
// can enable the domainLetsEncrypt feature to request a legitimate and valid certificate from Let's Encrypt,
// at which point you can enable verification for the HTTPS certificate.
var httpsInsecureVerify *bool

// noMediaTest is used to disable the media test.
//
// The media test, which involves using ffmpeg to publish and record a media stream, can be time-consuming
// and may occasionally fail. Therefore, it is necessary to retry the process if a failure occurs.
var noMediaTest *bool

// noBilibiliTest is used to disable the bilibili test.
var noBilibiliTest *bool

var srsLog *bool
var srsTimeout *int
var endpoint *string
var endpointRTMP *string
var endpointHTTP *string
var endpointSRT *string
var forceHttps *bool
var apiSecret *string
var checkApiSecret *bool
var waitReady *bool
var apiReadyimeout *int
var srsFFmpeg *string
var srsFFprobe *string
var srsFFmpegStderr *bool
var srsDVRStderr *bool
var srsFFprobeStdout *bool
var srsFFprobeDuration *int
var srsFFprobeTimeout *int
var srsInputFile *string

func options() string {
	return fmt.Sprintf("log=%v, timeout=%vms, secret=%vB, checkApiSecret=%v, endpoint=%v, forceHttps=%v, "+
		"waitReady=%v, initPassword=%v, initSelfSignedCert=%v, systemPassword=%vB, domainLetsEncrypt=%v, "+
		"httpsInsecureVerify=%v, srsInputFile=%v, noMediaTest=%v, noBilibiliTest=%v, endpointRTMP=%v, endpointHTTP=%v, "+
		"endpointSRT=%v",
		*srsLog, *srsTimeout, len(*apiSecret), *checkApiSecret, *endpoint, *forceHttps, *waitReady, *initPassword,
		*initSelfSignedCert, len(*systemPassword), *domainLetsEncrypt, *httpsInsecureVerify, *srsInputFile,
		*noMediaTest, *noBilibiliTest, *endpointRTMP, *endpointHTTP, *endpointSRT,
	)
}

func prepareTest(ctx context.Context) (err error) {
	// Try to load the .env file.
	for _, envFile := range []string{
		".env",
		"test/.env",
		"/data/config/.env",
		"../platform/containers/data/config/.env",
	} {
		if _, err := os.Stat(envFile); err == nil {
			if err := godotenv.Overload(envFile); err == nil {
				break
			}
		}
	}

	// Parse the options.
	srsLog = flag.Bool("srs-log", false, "Whether enable the detail log")
	srsTimeout = flag.Int("srs-timeout", 60000, "For each case, the timeout in ms")
	apiSecret = flag.String("api-secret", os.Getenv("SRS_PLATFORM_SECRET"), "The secret for api")
	checkApiSecret = flag.Bool("check-api-secret", true, "Whether check the api secret")
	endpoint = flag.String("endpoint", "http://localhost:2022", "The endpoint for api, can be http or https")
	endpointRTMP = flag.String("endpoint-rtmp", "rtmp://localhost:1935", "The endpoint for rtmp")
	endpointHTTP = flag.String("endpoint-http", "http://localhost:8080", "The endpoint for http")
	endpointSRT = flag.String("endpoint-srt", "srt://localhost:10080", "The endpoint for srt")
	forceHttps = flag.Bool("force-https", false, "Force to use HTTPS api")
	waitReady = flag.Bool("wait-ready", false, "Whether wait for the service ready")
	apiReadyimeout = flag.Int("api-ready-timeout", 30000, "Check when startup, the timeout in ms")
	initPassword = flag.Bool("init-password", false, "Whether init the system and set password")
	initSelfSignedCert = flag.Bool("init-self-signed-cert", false, "Whether init self-signed cert for HTTPS")
	domainLetsEncrypt = flag.String("domain-lets-encrypt", "", "Use the domain to test the letsencrypt feature")
	httpsInsecureVerify = flag.Bool("https-insecure-verify", false, "Whether verify the HTTPS certificate")
	systemPassword = flag.String("system-password", os.Getenv("MGMT_PASSWORD"), "The system password for login")
	srsFFmpeg = flag.String("srs-ffmpeg", "ffmpeg", "The FFmpeg tool")
	srsFFmpegStderr = flag.Bool("srs-ffmpeg-stderr", false, "Whether enable the FFmpeg stderr log")
	srsDVRStderr = flag.Bool("srs-dvr-stderr", false, "Whether enable the DVR stderr log")
	srsFFprobeStdout = flag.Bool("srs-ffprobe-stdout", false, "Whether enable the FFprobe stdout log")
	srsFFprobeDuration = flag.Int("srs-ffprobe-duration", 30000, "For each case, the duration for ffprobe in ms")
	srsFFprobeTimeout = flag.Int("srs-ffprobe-timeout", 40000, "For each case, the timeout for ffprobe in ms")
	srsFFprobe = flag.String("srs-ffprobe", "ffprobe", "The FFprobe tool")
	srsInputFile = flag.String("srs-input-file", "source.200kbps.768x320.flv", "The input file")
	noMediaTest = flag.Bool("no-media-test", false, "Whether disable the media test")
	noBilibiliTest = flag.Bool("no-bilibili-test", false, "Whether disable the bilibili test")

	// Should parse it first.
	flag.Parse()
	logger.Tf(ctx, "Test with %v", options())

	if *checkApiSecret && *apiSecret == "" {
		return errors.Errorf("empty api secret")
	}

	// Try to locate file.
	tryOpenFile := func(filename string) (string, error) {
		// Match if file exists.
		if _, err := os.Stat(filename); err == nil {
			return filename, nil
		}

		// If we run in GoLand, the current directory is in blackbox, so we use parent directory.
		nFilename := path.Join("../", filename)
		if stat, err := os.Stat(nFilename); err == nil && !stat.IsDir() {
			return nFilename, nil
		}

		// Try to open files in test directory.
		nFilename = path.Join("test", filename)
		if stat, err := os.Stat(nFilename); err == nil && !stat.IsDir() {
			return nFilename, nil
		}

		// Try to find file by which if it's a command like ffmpeg.
		cmd := exec.Command("which", filename)
		cmd.Env = []string{"PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
		if v, err := cmd.Output(); err == nil {
			return strings.TrimSpace(string(v)), nil
		}

		return filename, errors.Errorf("file %v not found", filename)
	}

	// Check and relocate path of tools.
	if *srsFFmpeg, err = tryOpenFile(*srsFFmpeg); err != nil {
		return errors.Wrapf(err, "not found ffmpeg %v", *srsFFmpeg)
	}
	if *srsFFprobe, err = tryOpenFile(*srsFFprobe); err != nil {
		return errors.Wrapf(err, "not found ffprobe %v", *srsFFprobe)
	}
	if *srsInputFile, err = tryOpenFile(*srsInputFile); err != nil {
		return errors.Wrapf(err, "not found input file %v", *srsInputFile)
	}

	return nil
}

func TestMain(m *testing.M) {
	ctx := logger.WithContext(context.Background())

	if err := prepareTest(ctx); err != nil {
		logger.Ef(ctx, "Prepare test fail, err %+v", err)
		os.Exit(-1)
	}

	// Disable the logger during all tests.
	if *srsLog == false {
		olw := logger.Switch(ioutil.Discard)
		defer func() {
			logger.Switch(olw)
		}()
	}

	// Init rand seed.
	rand.Seed(time.Now().UnixNano())

	// Wait for the service ready.
	if *waitReady {
		if err := waitForServiceReady(ctx); err != nil {
			os.Exit(-1)
		}
		logger.Tf(ctx, "Wait for service ready ok")
	}

	if *initPassword {
		if err := initSystemPassword(ctx); err != nil {
			logger.Ef(ctx, "Init system fail, err %+v", err)
			os.Exit(-1)
		}
		logger.Tf(ctx, "Init system password ok")
	}

	if *initSelfSignedCert {
		if err := apiRequest(ctx, "/terraform/v1/mgmt/auto-self-signed-certificate", nil, nil); err != nil {
			logger.Ef(ctx, "Init self-signed cert fail, err %+v", err)
			os.Exit(-1)
		}
		logger.Tf(ctx, "Init self-signed cert ok")
	}

	os.Exit(m.Run())
}

func waitForServiceReady(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*apiReadyimeout)*time.Millisecond)
	defer cancel()

	for {
		if ctx.Err() != nil {
			logger.Ef(ctx, "Wait for API ready timeout, err %v", ctx.Err())
			return ctx.Err()
		}

		err := apiRequest(ctx, "/terraform/v1/host/versions", nil, nil)
		if err == nil {
			break
		}

		logger.Tf(ctx, "Wait for API ready, err %v", err)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func initSystemPassword(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*srsTimeout)*time.Millisecond)
	defer cancel()

	// Set the password.
	password := fmt.Sprintf("%x", rand.Uint64())
	if *systemPassword != "" {
		password = *systemPassword
	} else {
		*systemPassword = password
	}

	// Initialize the system by password.
	var token string
	if err := apiRequest(ctx, "/terraform/v1/mgmt/init", &struct {
		Password string `json:"password"`
	}{
		Password: password,
	}, &struct {
		Token *string `json:"token"`
	}{
		Token: &token,
	}); err != nil {
		return errors.Wrapf(err, "init system")
	}
	if token == "" {
		return errors.Errorf("invalid token")
	}

	// Login the system by password.
	var token2 string
	if err := apiRequest(ctx, "/terraform/v1/mgmt/login", &struct {
		Password string `json:"password"`
	}{
		Password: password,
	}, &struct {
		Token *string `json:"token"`
	}{
		Token: &token2,
	}); err != nil {
		return errors.Wrapf(err, "login system")
	}
	if token2 == "" {
		return errors.Errorf("invalid token")
	}

	return nil
}
