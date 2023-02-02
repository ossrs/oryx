package common

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	tchttp "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/http"
)

const (
	codeLimitExceeded = "RequestLimitExceeded"
	tplRateLimitRetry = "[WARN] rate limit exceeded, retrying (%d/%d) in %f seconds: %s"
)

func (c *Client) sendWithRateLimitRetry(req *http.Request, retryable bool) (resp *http.Response, err error) {
	// make sure maxRetries is more than 0
	maxRetries := maxInt(c.profile.RateLimitExceededMaxRetries, 0)
	durationFunc := safeDurationFunc(c.profile.RateLimitExceededRetryDuration)

	var shadow []byte
	for idx := 0; idx <= maxRetries; idx++ {
		resp, err = c.sendWithNetworkFailureRetry(req, retryable)
		if err != nil {
			return
		}

		shadow, err = shadowRead(resp)
		if err != nil {
			return resp, err
		}

		err = tchttp.ParseErrorFromHTTPResponse(shadow)
		// should not sleep on last request
		if err, ok := err.(*errors.TencentCloudSDKError); ok && err.Code == codeLimitExceeded && idx < maxRetries {
			duration := durationFunc(idx)
			if c.debug {
				c.logger.Printf(tplRateLimitRetry, idx, maxRetries, duration.Seconds(), err.Error())
			}

			time.Sleep(duration)
			continue
		}

		return resp, err
	}

	return resp, err
}

func shadowRead(resp *http.Response) ([]byte, error) {
	var reader io.ReadCloser
	var err error
	var val []byte

	enc := resp.Header.Get("Content-Encoding")
	switch enc {
	case "":
		reader = resp.Body
	case "deflate":
		reader = flate.NewReader(resp.Body)
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("Content-Encoding not support: %s", enc)
	}

	val, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(val))

	// delete the header in case the caller mistake the body being encoded
	delete(resp.Header, "Content-Encoding")

	return val, nil
}
