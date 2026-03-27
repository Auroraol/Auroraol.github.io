package hooks

import (
	"net/url"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	resty "gopkg.in/resty.v1"
)

const (
	metricLabelHost   string = "host"
	metricLabelMethod string = "method"
	metricLabelPath   string = "path"
	// metricLabelURL      string = "url" // deprecated
	metricLabelHTTPCode string = "http_code"
	metricLabelCode     string = "code" // 业务code
)

func init() {
	prometheus.MustRegister(latencySummary, latencyCounter)
}

var (
	latencySummary = prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Name: "xd_sdk_httpclient_request_latency",
		Help: "Request latency of request that httpclient made (ms)",
		Objectives: map[float64]float64{
			0.5:  0.05,
			0.95: 0.01,
			0.99: 0.001,
		},
	}, []string{metricLabelHost, metricLabelPath})
	latencyCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "xd_sdk_httpclient_request_counter",
		Help: "Request counter of request that httpclient made",
	}, []string{metricLabelHost, metricLabelPath, metricLabelMethod, metricLabelHTTPCode, metricLabelCode})
)

func LatencyMetrics(parsePathF func(originPath string) string, fs ...ParseSvcCodeFunc) AfterResponseHook {
	return func(c *resty.Client, r *resty.Response) error {
		if r.Request == nil {
			return nil
		}
		u, _ := url.Parse(r.Request.URL)
		if u == nil {
			return nil
		}
		metricPath := u.Path
		if parsePathF != nil {
			if p := parsePathF(u.Path); p != "" {
				metricPath = p
			}
		}
		latencySummary.With(prometheus.Labels{
			metricLabelHost: u.Host,
			metricLabelPath: metricPath,
			// metricLabelURL: u.Host + metricPath,
		}).Observe(float64(r.ReceivedAt().Sub(r.Request.Time).Milliseconds()))

		labels := prometheus.Labels{
			// metricLabelPath:      u.Host + metricPath, // exclude param from url
			metricLabelHost:     u.Host,
			metricLabelPath:     metricPath,
			metricLabelMethod:   r.Request.Method,
			metricLabelHTTPCode: strconv.Itoa(r.StatusCode()),
			metricLabelCode:     "",
		}
		for _, f := range fs {
			if f == nil || (r.Result() == nil && len(r.Body()) == 0) {
				continue
			}
			obj := r.Result()
			if obj == nil {
				obj = r.Body()
			}
			if c, o := f(r.StatusCode(), obj); o {
				labels["code"] = strconv.Itoa(c)
				break
			}
		}
		latencyCounter.With(labels).Inc()
		return nil
	}
}
