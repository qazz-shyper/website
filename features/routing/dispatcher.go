package routing

import (
	"context"

	"github.com/qazz-shyper/website/common/net"
	"github.com/qazz-shyper/website/features"
	"github.com/qazz-shyper/website/transport"
)

// Dispatcher is a feature that dispatches inbound requests to outbound handlers based on rules.
// Dispatcher is required to be registered in a Website instance to make Website function properly.
//
// xray:api:stable
type Dispatcher interface {
	features.Feature

	// Dispatch returns a Ray for transporting data for the given request.
	Dispatch(ctx context.Context, dest net.Destination) (*transport.Link, error)
}

// DispatcherType returns the type of Dispatcher interface. Can be used to implement common.HasType.
//
// xray:api:stable
func DispatcherType() interface{} {
	return (*Dispatcher)(nil)
}