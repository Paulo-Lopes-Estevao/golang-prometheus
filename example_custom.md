``` go
package main

import (
	"log"
	"strings"


	"github.com/labstack/echo"
    "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)



func main() {
	route := echo.New()

	 var configMetrics = NewConfig()
	configMetrics.Namespace = "namespace"
	configMetrics.Buckets = []float64{
		0.0005, // 0.5ms
		0.001,  // 1ms
		0.005,  // 5ms
		0.01,   // 10ms
		0.05,   // 50ms
		0.1,    // 100ms
		0.5,    // 500ms
		1,      // 1s
		2,      // 2s
	} 

	route.Use(MiddlewareMetricsConfig(configMetrics))

	
	route.GET("/", func(c echo.Context) error {
		return c.JSON(200, "Hello world!")
	})

	route.GET("/metrics", echo.WrapHandler(promhttp.Handler()))

	if err := route.Start(":9000"); err != nil {
		log.Println("Not Running Server A...", err.Error())
	}

}




type Config struct {
	Namespace           string
	Buckets             []float64
	Subsystem           string
	NormalizeHTTPStatus bool
}

const (
	httpRequestsCount    = "http_requests_total"
	httpRequestsDuration = "request_duration_seconds"
	notFoundPath         = "/not-found"
)

var DefaultConfig = Config{
	Namespace: "echo",
	Subsystem: "http",
	Buckets: []float64{
		0.0005,
		0.001, // 1ms
		0.002,
		0.005,
		0.01, // 10ms
		0.02,
		0.05,
		0.1, // 100 ms
		0.2,
		0.5,
		1.0, // 1s
		2.0,
		5.0,
		10.0, // 10s
		15.0,
		20.0,
		30.0,
	},
	NormalizeHTTPStatus: true,
}

func NewConfig() Config {
	return DefaultConfig
}



func MiddlewareMetricsConfig(config Config) echo.MiddlewareFunc {

	httpRequests := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsCount,
		Help:      "Number of get requests.",
	}, []string{"status", "method", "handler"})

	httpDuration := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: config.Namespace,
		Subsystem: config.Subsystem,
		Name:      httpRequestsDuration,
		Help:      "Spend time by processing a route",
		Buckets:   config.Buckets,
	}, []string{"method", "handler"})

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			request_http := c.Request()
			path := c.Path()

			timer := prometheus.NewTimer(httpDuration.WithLabelValues(request_http.Method, path))
			err := next(c)
			timer.ObserveDuration()

			if err != nil {
				c.Error(err)
			}

			status := strconv.Itoa(c.Response().Status)

			httpRequests.WithLabelValues(status, request_http.Method, path).Inc()

			return err

		}
	}

}


```