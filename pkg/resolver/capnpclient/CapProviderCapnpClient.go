package capnpclient

import (
	"context"

	"capnproto.org/go/capnp/v3"

	"github.com/hiveot/hub.capnp/go/hubapi"
	"github.com/hiveot/hub/internal/caphelp"
	"github.com/hiveot/hub/pkg/resolver"
	"github.com/hiveot/hub/pkg/resolver/capserializer"
)

// CapProviderCapnpClient is a POGS wrapper around the provider capnp capability.
// This implements the IProvider interface as used by the service calling getCapability on the provider.

type CapProviderCapnpClient struct {
	capProvider hubapi.CapProvider // capnp client of the capability provider
}

// GetCapability obtains the capability with the given name.
// The caller must release the capability when done.
func (cl *CapProviderCapnpClient) GetCapability(ctx context.Context,
	clientID string, clientType string, capabilityName string, args []string) (
	capability capnp.Client, err error) {

	method, release := cl.capProvider.GetCapability(ctx,
		func(params hubapi.CapProvider_getCapability_Params) error {
			_ = params.SetClientID(clientID)
			_ = params.SetClientType(clientType)
			_ = params.SetCapabilityName(capabilityName)
			if args != nil {
				err = params.SetArgs(caphelp.MarshalStringList(args))
			}
			return err
		})
	defer release()
	// return a future. Caller must release
	// this does not detect a broken connection until the capability is used
	capability = method.Capability().AddRef()
	return capability, err
}

// ListCapabilities lists the available capabilities of the service
// Returns a list of capabilities that can be obtained through the service
func (cl *CapProviderCapnpClient) ListCapabilities(
	ctx context.Context) (infoList []resolver.CapabilityInfo, err error) {

	infoList = make([]resolver.CapabilityInfo, 0)
	method, release := cl.capProvider.ListCapabilities(ctx, nil)
	defer release()
	resp, err := method.Struct()
	if err == nil {
		infoListCapnp, err2 := resp.InfoList()
		if err = err2; err == nil {
			infoList = capserializer.UnmarshalCapabilyInfoList(infoListCapnp)
		}
	}
	return infoList, err
}

// Release this client
func (cl *CapProviderCapnpClient) Release() {
	cl.capProvider.Release()
}

// NewCapProviderCapnpClient create a new provider client for obtaining capnp capabilities.
func NewCapProviderCapnpClient(capProvider hubapi.CapProvider) (cl *CapProviderCapnpClient) {

	cl = &CapProviderCapnpClient{
		capProvider: capProvider,
	}
	return cl
}
