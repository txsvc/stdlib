package gcp

import (
	"context"
	"log"

	stackdriver_error "cloud.google.com/go/errorreporting"
	stackdriver_logging "cloud.google.com/go/logging"

	"github.com/txsvc/stdlib/observer"
	"github.com/txsvc/stdlib/pkg/env"
	"github.com/txsvc/stdlib/pkg/provider"
)

const (
	defaultLogId = "default"
	metricsLogId = "metrics"
)

type (
	// stackdriverImpl provides a simple implementation in the absence of any configuration
	stackdriverImpl struct {
		logger  *stackdriver_logging.Logger
		metrics *stackdriver_logging.Logger

		loggingClient *stackdriver_logging.Client
		errorClient   *stackdriver_error.Client
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*stackdriverImpl)(nil)

	_ observer.ErrorReportingProvider = (*stackdriverImpl)(nil)
	_ observer.LoggingProvider        = (*stackdriverImpl)(nil)
	_ observer.MetricsProvider        = (*stackdriverImpl)(nil)

	// the instance, a singleton
	theStackdriverProvider *stackdriverImpl
)

// a default provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	if theStackdriverProvider == nil {
		projectID := env.GetString("PROJECT_ID", "")
		serviceName := env.GetString("SERVICE_NAME", "default")

		// initialize logging
		lc, err := stackdriver_logging.NewClient(context.Background(), projectID)
		if err != nil {
			log.Fatal(err)
		}

		// initialize error reporting
		ec, err := stackdriver_error.NewClient(context.Background(), projectID, stackdriver_error.Config{
			ServiceName: serviceName,
			OnError: func(err error) {
				log.Printf("could not log error: %v", err)
			},
		})
		if err != nil || ec == nil {
			log.Fatal(err)
		}

		theStackdriverProvider = &stackdriverImpl{
			metrics:       lc.Logger(metricsLogId),
			logger:        lc.Logger(defaultLogId),
			loggingClient: lc,
			errorClient:   ec,
		}
	}
	return theStackdriverProvider
}

// IF GenericProvider

func (np *stackdriverImpl) Close() error {
	return nil
}

// IF ErrorReportingProvider

func (np *stackdriverImpl) ReportError(e error) {
	if e != nil {
		np.errorClient.Report(stackdriver_error.Entry{Error: e})
	}
}

// IF LoggingProvider

func (np *stackdriverImpl) Log(msg string, keyValuePairs ...string) {
	LogWithLevel(np.logger, observer.LevelInfo, msg, keyValuePairs...)
}

func (np *stackdriverImpl) LogWithLevel(lvl observer.Severity, msg string, keyValuePairs ...string) {
	LogWithLevel(np.logger, lvl, msg, keyValuePairs...)
}

// IF MetricsProvider

func (np *stackdriverImpl) Meter(ctx context.Context, metric string, args ...string) {
	LogWithLevel(np.metrics, observer.LevelInfo, metric, args...)
}

func LogWithLevel(logger *stackdriver_logging.Logger, lvl observer.Severity, msg string, keyValuePairs ...string) {
	e := stackdriver_logging.Entry{
		Payload:  msg,
		Severity: toSeverity(lvl),
	}

	n := len(keyValuePairs)
	if n > 0 {
		labels := make(map[string]string)
		if n == 1 {
			labels[keyValuePairs[0]] = ""
		} else {
			for i := 0; i < n/2; i++ {
				k := keyValuePairs[i*2]
				v := keyValuePairs[(i*2)+1]
				labels[k] = v
			}
			if n%2 == 1 {
				labels[keyValuePairs[n-1]] = ""
			}
		}
		e.Labels = labels
	}

	logger.Log(e)
}

func toSeverity(severity observer.Severity) stackdriver_logging.Severity {
	switch severity {
	case observer.LevelInfo:
		return stackdriver_logging.Info
	case observer.LevelWarn:
		return stackdriver_logging.Warning
	case observer.LevelError:
		return stackdriver_logging.Error
	case observer.LevelDebug:
		return stackdriver_logging.Debug
	}
	return stackdriver_logging.Info
}
