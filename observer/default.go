package observer

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"

	"github.com/txsvc/stdlib/pkg/provider"
)

type (
	// defaultObserverImpl provides a simple implementation in the absence of any configuration
	defaultObserverImpl struct {
		loggingDisabled bool
	}
)

var (
	// Interface guards.

	// This enforces a compile-time check of the provider implmentation,
	// making sure all the methods defined in the interfaces are implemented.

	_ provider.GenericProvider = (*defaultObserverImpl)(nil)

	_ ErrorReportingProvider = (*defaultObserverImpl)(nil)
	_ LoggingProvider        = (*defaultObserverImpl)(nil)
	_ MetricsProvider        = (*defaultObserverImpl)(nil)

	// the instance, a singleton
	theDefaultProvider *defaultObserverImpl
)

func init() {
	reset()
}

func reset() {
	// initialize the observer with a NULL provider that prevents NPEs in case someone forgets to initialize the platform with a real provider
	loggingConfig := provider.WithProvider("observer.default.logger", TypeLogger, NewDefaultProvider)
	errorReportingConfig := provider.WithProvider("observer.default.errorreporting", TypeErrorReporter, NewDefaultProvider)
	metricsConfig := provider.WithProvider("observer.default.metrics", TypeMetrics, NewDefaultProvider)

	NewConfig(loggingConfig, errorReportingConfig, metricsConfig)
}

// a default provider that does nothing but prevents NPEs in case someone forgets to actually initializa the 'real' provider
func NewDefaultProvider() interface{} {
	if theDefaultProvider == nil {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		// log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

		theDefaultProvider = &defaultObserverImpl{
			loggingDisabled: false,
		}
		theDefaultProvider.EnableLogging()
	}
	return theDefaultProvider
}

// IF GenericProvider

func (np *defaultObserverImpl) Close() error {
	return nil
}

// IF ErrorReportingProvider

func (np *defaultObserverImpl) ReportError(e error) {
	log.Error().Stack().Err(e).Msg("")
}

// IF LoggingProvider

func (np *defaultObserverImpl) Log(msg string, keyValuePairs ...string) {
	if np.loggingDisabled {
		return // just do nothing
	}
	log.Info().Msg(msg)
}

func (np *defaultObserverImpl) LogWithLevel(lvl Severity, msg string, keyValuePairs ...string) {
	if np.loggingDisabled {
		return // just do nothing
	}

	switch lvl {
	case LevelInfo:
		log.Info().Msg(msg)
	case LevelWarn:
		log.Warn().Msg(msg)
	case LevelError:
		log.Error().Msg(msg)
	case LevelDebug:
		log.Debug().Msg(msg)
	}
}

func (np *defaultObserverImpl) EnableLogging() {
	np.loggingDisabled = false
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func (np *defaultObserverImpl) DisableLogging() {
	np.loggingDisabled = true
	zerolog.SetGlobalLevel(zerolog.Disabled)
}

// IF MetricsProvider

func (np *defaultObserverImpl) Meter(ctx context.Context, metric string, args ...string) {
}
