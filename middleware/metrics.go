package middleware

import (
	"bufio"
	"net"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricsMiddleware struct {
	// counter of process, will be a counter of requests into the server
	opsProcessed *prometheus.CounterVec
}

type logger interface {
	Infof(format string, args ...interface{})
}

// A warpper to intercept set data on the ResponseWriter
type resposeWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

// initialize the counter
// should be more dynamic
func NewMetricsMiddleware() *MetricsMiddleware {
	opsProcessed := promauto.NewCounterVec(prometheus.CounterOpts{
		// name of the metric
		Name: "myapp_processed_ops_total",
		Help: "The total number of processed events",
		// labels of the metric, will use to filter the metric
	}, []string{"method", "path", "statuscode"})
	return &MetricsMiddleware{
		opsProcessed: opsProcessed,
	}
}

func (lm *MetricsMiddleware) Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// will happen before the request

		wi := &resposeWriterInterceptor{
			statusCode:     http.StatusOK,
			ResponseWriter: w,
		}
		// the request
		next.ServeHTTP(wi, r)
		// will happen after the request fished executing
		// increase the counter
		lm.opsProcessed.With(prometheus.Labels{"method": r.Method, "path": r.RequestURI,
			"statuscode": strconv.Itoa(wi.statusCode)}).Inc()
	})
}

func (w *resposeWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *resposeWriterInterceptor) Write(p []byte) (int, error) {
	return w.ResponseWriter.Write(p)
}

func (w *resposeWriterInterceptor) Flush() {
	f, _ := w.ResponseWriter.(http.Flusher)
	// if !ok :
	f.Flush()
}

func (w *resposeWriterInterceptor) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, _ := w.ResponseWriter.(http.Hijacker)
	// if !ok : nil, nil, errors.New("type assertion failed http.ResponseWriter not http.Hijack")
	return h.Hijack()
}
