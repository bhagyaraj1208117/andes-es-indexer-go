package core

import (
	"github.com/bhagyaraj1208117/andes-es-indexer-go/core/request"
	"github.com/bhagyaraj1208117/andes-es-indexer-go/metrics"
)

// StatusMetricsHandler defines the behavior of a component that handles status metrics
type StatusMetricsHandler interface {
	AddIndexingData(args metrics.ArgsAddIndexingData)
	GetMetrics() map[string]*request.MetricsResponse
	GetMetricsForPrometheus() string
	IsInterfaceNil() bool
}

// WebServerHandler defines the behavior of a component that handles the web server
type WebServerHandler interface {
	StartHttpServer() error
	Close() error
}
