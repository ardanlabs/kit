package tcp

import (
	"time"

	"github.com/ardanlabs/kit/pool"
)

// OptUserPool declares fields for the user to pass their own
// work pools for configuration.
type OptUserPool struct {
	RecvPool *pool.Pool // User provided work pool for the receive work.
	SendPool *pool.Pool // User provided work pool for the send work.
}

// OptIntPool declares fields for the user to provide configuration
// for an internally configured pool.
type OptIntPool struct {
	RecvMinPoolSize func() int // Min number of routines the recv pool must have.
	RecvMaxPoolSize func() int // Max number of routines the recv pool can have.
	SendMinPoolSize func() int // Min number of routines the send pool must have.
	SendMaxPoolSize func() int // Max number of routines the send pool can have.
}

// OptRateLimit declares fields for the user to provide configuration
// for connection rate limit.
type OptRateLimit struct {
	RateLimit func() time.Duration // Connection rate limit per single connection.
}

// Config provides a data structure of required configuration parameters.
type Config struct {
	NetType string // "tcp", tcp4" or "tcp6"
	Addr    string // "host:port" or "[ipv6-host%zone]:port"

	ConnHandler ConnHandler // Support for binding new connections to a reader and writer.
	ReqHandler  ReqHandler  // Support for handling the specific request workflow.
	RespHandler RespHandler // Support for handling the specific response workflow.

	// *************************************************************************
	// ** Required, choose one option.                                        **
	// *************************************************************************

	// Decide if you want to pass in your own work pool for configuration options
	// for the tcp value to create its own. Pass in your own pool if you want to
	// share a single pool across multiple tcp values.

	OptUserPool
	OptIntPool

	// *************************************************************************
	// ** Not Required, optional                                              **
	// *************************************************************************

	OptRateLimit
}

// Validate checks the configuration to required items.
func (cfg *Config) Validate() error {
	if cfg == nil {
		return ErrInvalidConfiguration
	}

	if cfg.NetType != "tcp" && cfg.NetType != "tcp4" && cfg.NetType != "tcp6" {
		return ErrInvalidNetType
	}

	if cfg.ConnHandler == nil {
		return ErrInvalidConnHandler
	}

	if cfg.ReqHandler == nil {
		return ErrInvalidReqHandler
	}

	if cfg.RespHandler == nil {
		return ErrInvalidRespHandler
	}

	if (cfg.RecvPool != nil && cfg.SendPool == nil) || (cfg.RecvPool == nil && cfg.SendPool != nil) {
		return ErrInvalidPoolConfiguration
	}

	return nil
}
