package hooks

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"

	resty "gopkg.in/resty.v1"
)

const maxBodyLen = 2 << 15

func LoggingRequest() PreRequestHook {
	return LoggingRequestWithLogger(logger.StandardLogger())
}

func LoggingRequestWithLogger(l logger.Logger) PreRequestHook {
	return func(c *resty.Client, r *resty.Request) error {
		l.WithFields(logger.Fields{
			"method":  r.RawRequest.Method,
			"url":     r.RawRequest.URL.String(),
			"headers": composeHeaders(r.RawRequest.Header),
			"body":    requestBody(c, r),
		}).Infof(r.Context(), "httpclient outgoing http request")
		return nil
	}
}

func LoggingResponse() AfterResponseHook {
	return LoggingResponseWithLogger(logger.StandardLogger())
}

func LoggingResponseWithLogger(l logger.Logger) AfterResponseHook {
	return func(c *resty.Client, r *resty.Response) error {
		l.WithFields(logger.Fields{
			"status":  r.Status(),
			"headers": composeHeaders(r.Header()),
			"body":    responseBody(c, r),
		}).Infof(r.Request.Context(), "httpclient incoming http response")
		return nil
	}
}

func composeHeaders(headers http.Header) string {
	pairs := make([]string, 0, len(headers))
	for key, values := range headers {
		pairs = append(pairs, fmt.Sprintf("%s: %s", key, strings.Join(values, ", ")))
	}
	sort.Strings(pairs)
	return strings.Join(pairs, "; ")
}

func requestBody(_ *resty.Client, r *resty.Request) string {
	if r.RawRequest.Body == nil || r.RawRequest.Body == http.NoBody {
		return ""
	}
	body, _ := ioutil.ReadAll(r.RawRequest.Body)
	_ = r.RawRequest.Body.Close()
	r.RawRequest.Body = ioutil.NopCloser(bytes.NewReader(body))
	bodyStr := string(body)
	cnt := len(bodyStr)
	if cnt > maxBodyLen {
		cnt = maxBodyLen
	}
	return bodyStr[:cnt]
}

func responseBody(_ *resty.Client, r *resty.Response) string {
	cnt := len(r.String())
	if cnt > maxBodyLen {
		cnt = maxBodyLen
	}
	return r.String()[:cnt]
}
