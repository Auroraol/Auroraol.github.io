package hooks

import (
	"net/http/httptrace"

	"gitlab.xiaoduoai.com/golib/xd_sdk/logger"
	resty "gopkg.in/resty.v1"
)

// DNSTrace DNSTrace
func DNSTrace() PreRequestHook {
	return DNSTraceWithLogger(logger.StandardLogger())
}

// DNSTraceWithLogger DNSTraceWithLogger
func DNSTraceWithLogger(l logger.Logger) PreRequestHook {
	return func(_ *resty.Client, req *resty.Request) error {
		trace := &httptrace.ClientTrace{
			DNSStart: func(dnsInfo httptrace.DNSStartInfo) {
				l.WithField("host", req.RawRequest.URL.Host).Debugf(req.Context(), "doing dns resolve for host %s", dnsInfo.Host)
			},
			DNSDone: func(dnsDoneInfo httptrace.DNSDoneInfo) {
				l.WithField("host", req.RawRequest.URL.Host).Debugf(req.Context(), "done dns resolve for host %v", dnsDoneInfo.Addrs)
			},
		}
		req.RawRequest = req.RawRequest.WithContext(httptrace.WithClientTrace(req.RawRequest.Context(), trace))
		return nil
	}
}
