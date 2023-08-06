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
	"os/exec"
	"path"
	"strconv"
	"strings"
	"sync"
	"syscall"
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
var systemPassword *string
var srsFFmpeg *string
var srsFFprobe *string
var srsFFmpegStderr *bool
var srsDVRStderr *bool
var srsFFprobeStdout *bool
var srsFFprobeDuration *int
var srsFFprobeTimeout *int

func options() string {
	return fmt.Sprintf("log=%v, timeout=%vms, secret=%vB, checkApiSecret=%v, endpoint=%v, waitReady=%v, initPassword=%v, systemPassword=%vB",
		*srsLog, *srsTimeout, len(*apiSecret), *checkApiSecret, *endpoint, *waitReady, *initPassword, len(*systemPassword))
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
	endpoint = flag.String("endpoint", "http://localhost:2022", "The endpoint for api")
	waitReady = flag.Bool("wait-ready", false, "Whether wait for the service ready")
	apiReadyimeout = flag.Int("api-ready-timeout", 30000, "Check when startup, the timeout in ms")
	initPassword = flag.Bool("init-password", false, "Whether init the system and set password")
	systemPassword = flag.String("system-password", os.Getenv("MGMT_PASSWORD"), "The system password for login")
	srsFFmpeg = flag.String("srs-ffmpeg", "ffmpeg", "The FFmpeg tool")
	srsFFmpegStderr = flag.Bool("srs-ffmpeg-stderr", false, "Whether enable the FFmpeg stderr log")
	srsDVRStderr = flag.Bool("srs-dvr-stderr", false, "Whether enable the DVR stderr log")
	srsFFprobeStdout = flag.Bool("srs-ffprobe-stdout", false, "Whether enable the FFprobe stdout log")
	srsFFprobeDuration = flag.Int("srs-ffprobe-duration", 30000, "For each case, the duration for ffprobe in ms")
	srsFFprobeTimeout = flag.Int("srs-ffprobe-timeout", 40000, "For each case, the timeout for ffprobe in ms")
	srsFFprobe = flag.String("srs-ffprobe", "ffprobe", "The FFprobe tool")

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
		if _, err := os.Stat(nFilename); err == nil {
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
		return err
	}
	if *srsFFprobe, err = tryOpenFile(*srsFFprobe); err != nil {
		return err
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

	m := http.MethodPost
	if body == nil {
		m = http.MethodGet
	}

	u := fmt.Sprintf("%s%s", *endpoint, apiPath)
	req, err := http.NewRequestWithContext(ctx, m, u, body)
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

type backendService struct {
	// The context for case.
	caseCtx       context.Context
	caseCtxCancel context.CancelFunc

	// When SRS process started.
	readyCtx       context.Context
	readyCtxCancel context.CancelFunc

	// Whether already closed.
	closedCtx       context.Context
	closedCtxCancel context.CancelFunc

	// All goroutines
	wg sync.WaitGroup

	// The name, args and env for cmd.
	name string
	args []string
	env  []string
	// If timeout, kill the process.
	duration time.Duration

	// The process stdout and stderr.
	stdout bytes.Buffer
	stderr bytes.Buffer
	// The process error.
	r0 error
	// The process pid.
	pid int
	// Whether ignore process exit status error.
	ignoreExitStatusError bool

	// Hooks for owner.
	// Before start the process.
	onBeforeStart func(ctx context.Context, bs *backendService, cmd *exec.Cmd) error
	// After started the process.
	onAfterStart func(ctx context.Context, bs *backendService, cmd *exec.Cmd) error
	// Before kill the process, when case is done.
	onBeforeKill func(ctx context.Context, bs *backendService, cmd *exec.Cmd) error
	// After stopped the process. Always callback when run is called.
	onStop func(ctx context.Context, bs *backendService, cmd *exec.Cmd, r0 error, stdout, stderr *bytes.Buffer) error
	// When dispose the process. Always callback when run is called.
	onDispose func(ctx context.Context, bs *backendService) error
}

func newBackendService(opts ...func(v *backendService)) *backendService {
	v := &backendService{}

	v.readyCtx, v.readyCtxCancel = context.WithCancel(context.Background())
	v.closedCtx, v.closedCtxCancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *backendService) Close() error {
	if v.closedCtx.Err() != nil {
		return v.r0
	}
	v.closedCtxCancel()

	if v.caseCtxCancel != nil {
		v.caseCtxCancel()
	}
	if v.readyCtxCancel != nil {
		v.readyCtxCancel()
	}

	v.wg.Wait()

	if v.onDispose != nil {
		v.onDispose(v.caseCtx, v)
	}

	logger.Tf(v.caseCtx, "Process is closed, pid=%v, r0=%v", v.pid, v.r0)
	return nil
}

func (v *backendService) ReadyCtx() context.Context {
	return v.readyCtx
}

func (v *backendService) Run(ctx context.Context, cancel context.CancelFunc) error {
	// Always dispose resource of process.
	defer v.Close()

	// Start SRS with -e, which only use environment variables.
	cmd := exec.Command(v.name, v.args...)

	// If not started, we also need to callback the onStop.
	var processStarted bool
	defer func() {
		if v.onStop != nil && !processStarted {
			v.onStop(ctx, v, cmd, v.r0, &v.stdout, &v.stderr)
		}
	}()

	// Ignore if already error.
	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Save the context of case.
	v.caseCtx, v.caseCtxCancel = ctx, cancel

	// Setup stdout and stderr.
	cmd.Stdout = &v.stdout
	cmd.Stderr = &v.stderr
	cmd.Env = v.env
	if v.onBeforeStart != nil {
		if err := v.onBeforeStart(ctx, v, cmd); err != nil {
			return errors.Wrapf(err, "onBeforeStart failed")
		}
	}

	// Try to start the SRS server.
	if err := cmd.Start(); err != nil {
		return err
	}

	// Now process started, query the pid.
	v.pid = cmd.Process.Pid
	v.readyCtxCancel()
	processStarted = true
	if v.onAfterStart != nil {
		if err := v.onAfterStart(ctx, v, cmd); err != nil {
			return errors.Wrapf(err, "onAfterStart failed")
		}
	}

	// The context for SRS process.
	processDone, processDoneCancel := context.WithCancel(context.Background())

	// If exceed timeout, kill the process.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()
		if v.duration <= 0 {
			return
		}

		select {
		case <-ctx.Done():
		case <-time.After(v.duration):
			logger.Tf(ctx, "Process killed duration=%v, pid=%v, name=%v, args=%v", v.duration, v.pid, v.name, v.args)
			cmd.Process.Kill()
		}
	}()

	// If SRS process terminated, notify case to stop.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		// When SRS quit, also terminate the case.
		defer cancel()

		// Notify other goroutine, SRS already done.
		defer processDoneCancel()

		if err := cmd.Wait(); err != nil && !v.ignoreExitStatusError {
			v.r0 = errors.Wrapf(err, "Process wait err, pid=%v, name=%v, args=%v", v.pid, v.name, v.args)
		}
		if v.onStop != nil {
			if err := v.onStop(ctx, v, cmd, v.r0, &v.stdout, &v.stderr); err != nil {
				if v.r0 == nil {
					v.r0 = errors.Wrapf(err, "Process onStop err, pid=%v, name=%v, args=%v", v.pid, v.name, v.args)
				} else {
					logger.Ef(ctx, "Process onStop err %v", err)
				}
			}
		}
	}()

	// If case terminated, notify SRS process to stop.
	v.wg.Add(1)
	go func() {
		defer v.wg.Done()

		select {
		case <-ctx.Done():
			// Notify owner that we're going to kill the process.
			if v.onBeforeKill != nil {
				v.onBeforeKill(ctx, v, cmd)
			}

			// When case terminated, also terminate the SRS process.
			cmd.Process.Signal(syscall.SIGINT)
		case <-processDone.Done():
			// Ignore if already done.
			return
		}

		// Start a goroutine to ensure process killed.
		go func() {
			time.Sleep(3 * time.Second)
			if processDone.Err() == nil { // Ignore if already done.
				cmd.Process.Signal(syscall.SIGKILL)
			}
		}()
	}()

	// Wait for SRS or case done.
	select {
	case <-ctx.Done():
	case <-processDone.Done():
	}

	return v.r0
}

// ServiceRunner is an interface to run backend service.
type ServiceRunner interface {
	Run(ctx context.Context, cancel context.CancelFunc) error
}

// ServiceReadyQuerier is an interface to detect whether service is ready.
type ServiceReadyQuerier interface {
	ReadyCtx() context.Context
}

type FFmpegClient interface {
	ServiceRunner
	ServiceReadyQuerier
}

type ffmpegClient struct {
	// The backend service process.
	process *backendService

	// FFmpeg cli args, without ffmpeg binary.
	args []string
	// Let the process quit, do not cancel the case.
	cancelCaseWhenQuit bool
	// When timeout, stop FFmpeg, sometimes the '-t' does not work.
	ffmpegDuration time.Duration
}

func NewFFmpeg(opts ...func(v *ffmpegClient)) FFmpegClient {
	v := &ffmpegClient{
		process:            newBackendService(),
		cancelCaseWhenQuit: true,
	}

	// Do cleanup.
	v.process.onDispose = func(ctx context.Context, bs *backendService) error {
		return nil
	}

	// We ignore any exit error, because FFmpeg might exit with error even publish ok.
	v.process.ignoreExitStatusError = true

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *ffmpegClient) ReadyCtx() context.Context {
	return v.process.ReadyCtx()
}

func (v *ffmpegClient) Run(ctx context.Context, cancel context.CancelFunc) error {
	logger.Tf(ctx, "Starting FFmpeg by %v", strings.Join(v.args, " "))

	v.process.name = *srsFFmpeg
	v.process.args = v.args
	v.process.env = os.Environ()
	v.process.duration = v.ffmpegDuration

	v.process.onStop = func(ctx context.Context, bs *backendService, cmd *exec.Cmd, r0 error, stdout, stderr *bytes.Buffer) error {
		logger.Tf(ctx, "FFmpeg process pid=%v exit, r0=%v, stdout=%v", bs.pid, r0, stdout.String())
		if *srsFFmpegStderr && stderr.Len() > 0 {
			logger.Tf(ctx, "FFmpeg process pid=%v, stderr is \n%v", bs.pid, stderr.String())
		}
		return nil
	}

	// We might not want to cancel the case, for example, when check DVR by session, we just let the FFmpeg process to
	// quit and we should check the callback and DVR file.
	ffCtx, ffCancel := context.WithCancel(ctx)
	go func() {
		select {
		case <-ctx.Done():
		case <-ffCtx.Done():
			if v.cancelCaseWhenQuit {
				cancel()
			}
		}
	}()

	return v.process.Run(ffCtx, ffCancel)
}

type FFprobeClient interface {
	ServiceRunner
	// ProbeDoneCtx indicates the probe is done.
	ProbeDoneCtx() context.Context
	// Result return the raw string and metadata.
	Result() (string, *ffprobeObject)
}

type ffprobeClient struct {
	// The DVR file for ffprobe. Stream should be DVR to file, then use ffprobe to detect it. If DVR by FFmpeg, we will
	// start a FFmpeg process to do the DVR, or the DVR should be done by other tools.
	dvrFile string
	// The timeout to wait for task to done.
	timeout time.Duration

	// Whether do DVR by FFmpeg, if using SRS DVR, please set to false.
	dvrByFFmpeg bool
	// The stream to DVR for probing. Ignore if not DVR by ffmpeg
	streamURL string
	// The duration of video file for DVR and probing.
	duration time.Duration

	// When probe stream metadata object.
	doneCtx    context.Context
	doneCancel context.CancelFunc
	// The metadata object.
	metadata *ffprobeObject
	// The raw string of ffprobe.
	rawString string
}

func NewFFprobe(opts ...func(v *ffprobeClient)) FFprobeClient {
	v := &ffprobeClient{
		metadata:    &ffprobeObject{},
		dvrByFFmpeg: true,
	}
	v.doneCtx, v.doneCancel = context.WithCancel(context.Background())

	for _, opt := range opts {
		opt(v)
	}

	return v
}

func (v *ffprobeClient) ProbeDoneCtx() context.Context {
	return v.doneCtx
}

func (v *ffprobeClient) Result() (string, *ffprobeObject) {
	return v.rawString, v.metadata
}

func (v *ffprobeClient) Run(ctxCase context.Context, cancelCase context.CancelFunc) error {
	if true {
		ctx, cancel := context.WithTimeout(ctxCase, v.timeout)
		defer cancel()

		logger.Tf(ctx, "Starting FFprobe for stream=%v, dvr=%v, duration=%v, timeout=%v",
			v.streamURL, v.dvrFile, v.duration, v.timeout)

		// Try to start a DVR process.
		for ctx.Err() == nil {
			// If not DVR by FFmpeg, we just wait the DVR file to be ready, and it should be done by SRS or other tools.
			if v.dvrByFFmpeg {
				// If error, just ignore and retry, because the stream might not be ready. For example, for HLS, the DVR process
				// might need to wait for a duration of segment, 10s as such.
				_ = v.doDVR(ctx)
			}

			// Check whether DVR file is ok.
			if fs, err := os.Stat(v.dvrFile); err == nil && fs.Size() > 1024 {
				logger.Tf(ctx, "DVR FFprobe file is ok, file=%v, size=%v", v.dvrFile, fs.Size())
				break
			}

			// If not DVR by FFmpeg, must be by other tools, only need to wait.
			if !v.dvrByFFmpeg {
				logger.Tf(ctx, "Waiting stream=%v to be DVR", v.streamURL)
			}

			// Wait for a while and retry. Use larger timeout for HLS.
			retryTimeout := 1 * time.Second
			if strings.Contains(v.streamURL, ".m3u8") || v.dvrFile == "" {
				retryTimeout = 3 * time.Second
			}

			select {
			case <-ctx.Done():
			case <-time.After(retryTimeout):
			}
		}
	}

	// Ignore if case terminated.
	if ctxCase.Err() != nil {
		return nil
	}

	// Start a probe process for the DVR file.
	return v.doProbe(ctxCase, cancelCase)
}

func (v *ffprobeClient) doDVR(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	if !v.dvrByFFmpeg {
		return nil
	}

	process := newBackendService()
	process.name = *srsFFmpeg
	process.args = []string{
		"-t", fmt.Sprintf("%v", int64(v.duration/time.Second)),
		"-i", v.streamURL, "-c", "copy", "-y", v.dvrFile,
	}
	process.env = os.Environ()

	process.onDispose = func(ctx context.Context, bs *backendService) error {
		return nil
	}
	process.onBeforeStart = func(ctx context.Context, bs *backendService, cmd *exec.Cmd) error {
		logger.Tf(ctx, "DVR start %v %v", bs.name, strings.Join(bs.args, " "))
		return nil
	}
	process.onStop = func(ctx context.Context, bs *backendService, cmd *exec.Cmd, r0 error, stdout, stderr *bytes.Buffer) error {
		logger.Tf(ctx, "DVR process pid=%v exit, r0=%v, stdout=%v", bs.pid, r0, stdout.String())
		if *srsDVRStderr && stderr.Len() > 0 {
			logger.Tf(ctx, "DVR process pid=%v, stderr is \n%v", bs.pid, stderr.String())
		}
		return nil
	}

	return process.Run(ctx, cancel)
}

func (v *ffprobeClient) doProbe(ctx context.Context, cancel context.CancelFunc) error {
	process := newBackendService()
	process.name = *srsFFprobe
	process.args = []string{
		"-show_error", "-show_private_data", "-v", "quiet", "-find_stream_info",
		"-analyzeduration", fmt.Sprintf("%v", int64(v.duration/time.Microsecond)),
		"-print_format", "json", "-show_format", "-show_streams", v.dvrFile,
	}
	process.env = os.Environ()

	process.onDispose = func(ctx context.Context, bs *backendService) error {
		if _, err := os.Stat(v.dvrFile); !os.IsNotExist(err) {
			os.Remove(v.dvrFile)
		}
		return nil
	}
	process.onBeforeStart = func(ctx context.Context, bs *backendService, cmd *exec.Cmd) error {
		logger.Tf(ctx, "FFprobe start %v %v", bs.name, strings.Join(bs.args, " "))
		return nil
	}
	process.onStop = func(ctx context.Context, bs *backendService, cmd *exec.Cmd, r0 error, stdout, stderr *bytes.Buffer) error {
		logger.Tf(ctx, "FFprobe process pid=%v exit, r0=%v, stderr=%v", bs.pid, r0, stderr.String())
		if *srsFFprobeStdout && stdout.Len() > 0 {
			logger.Tf(ctx, "FFprobe process pid=%v, stdout is \n%v", bs.pid, stdout.String())
		}

		str := stdout.String()
		v.rawString = str

		if err := json.Unmarshal([]byte(str), v.metadata); err != nil {
			return err
		}

		m := v.metadata
		logger.Tf(ctx, "FFprobe done pid=%v, %v", bs.pid, m.String())

		v.doneCancel()
		return nil
	}

	return process.Run(ctx, cancel)
}

/*
   "index": 0,
   "codec_name": "h264",
   "codec_long_name": "H.264 / AVC / MPEG-4 AVC / MPEG-4 part 10",
   "profile": "High",
   "codec_type": "video",
   "codec_tag_string": "avc1",
   "codec_tag": "0x31637661",
   "width": 768,
   "height": 320,
   "coded_width": 768,
   "coded_height": 320,
   "closed_captions": 0,
   "film_grain": 0,
   "has_b_frames": 2,
   "sample_aspect_ratio": "1:1",
   "display_aspect_ratio": "12:5",
   "pix_fmt": "yuv420p",
   "level": 32,
   "chroma_location": "left",
   "field_order": "progressive",
   "refs": 1,
   "is_avc": "true",
   "nal_length_size": "4",
   "id": "0x1",
   "r_frame_rate": "25/1",
   "avg_frame_rate": "25/1",
   "time_base": "1/16000",
   "start_pts": 1280,
   "start_time": "0.080000",
   "duration_ts": 160000,
   "duration": "10.000000",
   "bit_rate": "196916",
   "bits_per_raw_sample": "8",
   "nb_frames": "250",
   "extradata_size": 41,
   "disposition": {
       "default": 1,
       "dub": 0,
       "original": 0,
       "comment": 0,
       "lyrics": 0,
       "karaoke": 0,
       "forced": 0,
       "hearing_impaired": 0,
       "visual_impaired": 0,
       "clean_effects": 0,
       "attached_pic": 0,
       "timed_thumbnails": 0,
       "captions": 0,
       "descriptions": 0,
       "metadata": 0,
       "dependent": 0,
       "still_image": 0
   },
   "tags": {
       "language": "und",
       "handler_name": "VideoHandler",
       "vendor_id": "[0][0][0][0]"
   }
*/
/*
   "index": 1,
   "codec_name": "aac",
   "codec_long_name": "AAC (Advanced Audio Coding)",
   "profile": "LC",
   "codec_type": "audio",
   "codec_tag_string": "mp4a",
   "codec_tag": "0x6134706d",
   "sample_fmt": "fltp",
   "sample_rate": "44100",
   "channels": 2,
   "channel_layout": "stereo",
   "bits_per_sample": 0,
   "id": "0x2",
   "r_frame_rate": "0/0",
   "avg_frame_rate": "0/0",
   "time_base": "1/44100",
   "start_pts": 132,
   "start_time": "0.002993",
   "duration_ts": 441314,
   "duration": "10.007120",
   "bit_rate": "29827",
   "nb_frames": "431",
   "extradata_size": 2,
   "disposition": {
       "default": 1,
       "dub": 0,
       "original": 0,
       "comment": 0,
       "lyrics": 0,
       "karaoke": 0,
       "forced": 0,
       "hearing_impaired": 0,
       "visual_impaired": 0,
       "clean_effects": 0,
       "attached_pic": 0,
       "timed_thumbnails": 0,
       "captions": 0,
       "descriptions": 0,
       "metadata": 0,
       "dependent": 0,
       "still_image": 0
   },
   "tags": {
       "language": "und",
       "handler_name": "SoundHandler",
       "vendor_id": "[0][0][0][0]"
   }
*/
type ffprobeObjectMedia struct {
	Index          int    `json:"index"`
	CodecName      string `json:"codec_name"`
	CodecType      string `json:"codec_type"`
	Timebase       string `json:"time_base"`
	Bitrate        string `json:"bit_rate"`
	Profile        string `json:"profile"`
	Duration       string `json:"duration"`
	CodecTagString string `json:"codec_tag_string"`

	// For video codec.
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	CodedWidth   int    `json:"coded_width"`
	CodedHeight  int    `json:"coded_height"`
	RFramerate   string `json:"r_frame_rate"`
	AvgFramerate string `json:"avg_frame_rate"`
	PixFmt       string `json:"pix_fmt"`
	Level        int    `json:"level"`

	// For audio codec.
	Channels      int    `json:"channels"`
	ChannelLayout string `json:"channel_layout"`
	SampleFmt     string `json:"sample_fmt"`
	SampleRate    string `json:"sample_rate"`
}

func (v *ffprobeObjectMedia) String() string {
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("index=%v, codec=%v, type=%v, tb=%v, bitrate=%v, profile=%v, duration=%v",
		v.Index, v.CodecName, v.CodecType, v.Timebase, v.Bitrate, v.Profile, v.Duration))
	sb.WriteString(fmt.Sprintf(", codects=%v", v.CodecTagString))

	if v.CodecType == "video" {
		sb.WriteString(fmt.Sprintf(", size=%vx%v, csize=%vx%v, rfr=%v, afr=%v, pix=%v, level=%v",
			v.Width, v.Height, v.CodedWidth, v.CodedHeight, v.RFramerate, v.AvgFramerate, v.PixFmt, v.Level))
	} else if v.CodecType == "audio" {
		sb.WriteString(fmt.Sprintf(", channels=%v, layout=%v, fmt=%v, srate=%v",
			v.Channels, v.ChannelLayout, v.SampleFmt, v.SampleRate))
	}

	return sb.String()
}

/*
"filename": "../objs/srs-ffprobe-stream-84487-8369019999559815097.mp4",
"nb_streams": 2,
"nb_programs": 0,
"format_name": "mov,mp4,m4a,3gp,3g2,mj2",
"format_long_name": "QuickTime / MOV",
"start_time": "0.002993",
"duration": "10.080000",
"size": "292725",
"bit_rate": "232321",
"probe_score": 100,

	"tags": {
	    "major_brand": "isom",
	    "minor_version": "512",
	    "compatible_brands": "isomiso2avc1mp41",
	    "encoder": "Lavf59.27.100"
	}
*/
type ffprobeObjectFormat struct {
	Filename   string `json:"filename"`
	Duration   string `json:"duration"`
	NBStream   int16  `json:"nb_streams"`
	Size       string `json:"size"`
	Bitrate    string `json:"bit_rate"`
	ProbeScore int    `json:"probe_score"`
}

func (v *ffprobeObjectFormat) String() string {
	return fmt.Sprintf("file=%v, duration=%v, score=%v, size=%v, bitrate=%v, streams=%v",
		v.Filename, v.Duration, v.ProbeScore, v.Size, v.Bitrate, v.NBStream)
}

/*
	{
	    "streams": [{ffprobeObjectMedia}, {ffprobeObjectMedia}],
	    "format": {ffprobeObjectFormat}
	}
*/
type ffprobeObject struct {
	Format  ffprobeObjectFormat  `json:"format"`
	Streams []ffprobeObjectMedia `json:"streams"`
}

func (v *ffprobeObject) String() string {
	sb := strings.Builder{}
	sb.WriteString(v.Format.String())
	sb.WriteString(", [")
	for _, stream := range v.Streams {
		sb.WriteString("{")
		sb.WriteString(stream.String())
		sb.WriteString("}")
	}
	sb.WriteString("]")
	return sb.String()
}

func (v *ffprobeObject) Duration() time.Duration {
	dv, err := strconv.ParseFloat(v.Format.Duration, 10)
	if err != nil {
		return time.Duration(0)
	}

	return time.Duration(dv*1000) * time.Millisecond
}

func (v *ffprobeObject) Video() *ffprobeObjectMedia {
	for _, media := range v.Streams {
		if media.CodecType == "video" {
			return &media
		}
	}
	return nil
}

func (v *ffprobeObject) Audio() *ffprobeObjectMedia {
	for _, media := range v.Streams {
		if media.CodecType == "audio" {
			return &media
		}
	}
	return nil
}
