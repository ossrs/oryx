package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/ossrs/go-oryx-lib/errors"
	"github.com/ossrs/go-oryx-lib/logger"
)

var srsLog *bool
var srsTimeout *int
var endpoint *string
var apiSecret *string
var checkApiSecret *bool
var waitReady *bool
var apiReadyimeout *int
var initPassword *bool

func options() string {
	return fmt.Sprintf("log=%v, timeout=%vms, secret=%vB, checkApiSecret=%v, endpoint=%v, waitReady=%v, initPassword=%v",
		*srsLog, *srsTimeout, len(*apiSecret), *checkApiSecret, *endpoint, *waitReady, *initPassword)
}

func prepareTest(ctx context.Context) (err error) {
	envFile := "../platform/containers/data/config/.env"
	if _, err := os.Stat(envFile); err == nil {
		if err := godotenv.Load(envFile); err != nil {
			return errors.Wrapf(err, "load %v", envFile)
		}
	}

	srsLog = flag.Bool("srs-log", false, "Whether enable the detail log")
	srsTimeout = flag.Int("srs-timeout", 5000, "For each case, the timeout in ms")
	apiSecret = flag.String("api-secret", os.Getenv("SRS_PLATFORM_SECRET"), "The secret for api")
	checkApiSecret = flag.Bool("check-api-secret", true, "Whether check the api secret")
	endpoint = flag.String("endpoint", "http://localhost:2022", "The endpoint for api")
	waitReady = flag.Bool("wait-ready", false, "Whether wait for the service ready")
	apiReadyimeout = flag.Int("api-ready-timeout", 30000, "Check when startup, the timeout in ms")
	initPassword = flag.Bool("init-password", false, "Whether init the system and set password")

	// Should parse it first.
	flag.Parse()

	if *checkApiSecret && *apiSecret == "" {
		return errors.Errorf("empty api secret")
	}

	return nil
}

func waitForServiceReady(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*apiReadyimeout)*time.Millisecond)
	logger.Tf(ctx, "Wait for API ready with %v, apiReadyimeout=%vms", options(), *apiReadyimeout)
	defer cancel()

	for {
		if ctx.Err() != nil {
			logger.Ef(ctx, "Wait for API ready timeout, err %v", ctx.Err())
			return ctx.Err()
		}

		err := apiRequest(ctx, "/terraform/v1/host/versions", nil, nil)
		if err == nil {
			logger.T(ctx, "API ready")
			break
		}

		logger.Tf(ctx, "Wait for API ready, err %v", err)
		time.Sleep(1 * time.Second)
	}

	return nil
}

func initSystemPassword(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(*srsTimeout)*time.Millisecond)
	logger.Tf(ctx, "Wait for API ready with %v", options())
	defer cancel()

	// Initialize the system by password.
	password := fmt.Sprintf("%x", rand.Uint64())
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
	}

	if *initPassword {
		if err := initSystemPassword(ctx); err != nil {
			logger.Ef(ctx, "Init system fail, err %+v", err)
			os.Exit(-1)
		}
	}

	os.Exit(m.Run())
}

// Filter the test error, ignore context.Canceled
func filterTestError(errs ...error) error {
	var filteredErrors []error

	for _, err := range errs {
		if err == nil || errors.Cause(err) == context.Canceled {
			continue
		}

		// If url error, server maybe error, do not print the detail log.
		if r0 := errors.Cause(err); r0 != nil {
			if r1, ok := r0.(*url.Error); ok {
				err = r1
			}
		}

		filteredErrors = append(filteredErrors, err)
	}

	if len(filteredErrors) == 0 {
		return nil
	}
	if len(filteredErrors) == 1 {
		return filteredErrors[0]
	}

	var descs []string
	for i, err := range filteredErrors[1:] {
		descs = append(descs, fmt.Sprintf("err #%d, %+v", i, err))
	}
	return errors.Wrapf(filteredErrors[0], "with %v", strings.Join(descs, ","))
}

func apiRequest(ctx context.Context, apiPath string, data interface{}, response interface{}) (err error) {
	var body io.Reader
	if data != nil {
		if b, err := json.Marshal(data); err != nil {
			return errors.Wrapf(err, "marshal data")
		} else {
			body = bytes.NewReader(b)
		}
	}

	u := fmt.Sprintf("%s%s", *endpoint, apiPath)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, body)
	if err != nil {
		return errors.Wrapf(err, "new request")
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %v", *apiSecret))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "do request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("invalid status code %v", resp.StatusCode)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "read body")
	}

	obj := &struct {
		Code int         `json:"code"`
		Data interface{} `json:"data"`
	}{
		Data: response,
	}
	if err = json.Unmarshal(b, obj); err != nil {
		return errors.Wrapf(err, "unmarshal %s", b)
	}

	if obj.Code != 0 {
		return errors.Errorf("invalid code %v of %s", obj.Code, b)
	}

	return nil
}
