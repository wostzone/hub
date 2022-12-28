package capnpclient

import (
	"context"
	"net"

	"capnproto.org/go/capnp/v3/rpc"

	"github.com/hiveot/hub.capnp/go/hubapi"
	"github.com/hiveot/hub/pkg/provisioning"
)

// ProvisioningCapnpClient provides a POGS wrapper around the generated provisioning capnp client
// This implements the IProvisioning interface
type ProvisioningCapnpClient struct {
	connection *rpc.Conn              // connection to the capnp server
	capability hubapi.CapProvisioning // capnp client
}

// CapManageProvisioning provides the capability to manage provisioning requests
func (cl *ProvisioningCapnpClient) CapManageProvisioning(
	ctx context.Context, clientID string) provisioning.IManageProvisioning {

	getCap, release := cl.capability.CapManageProvisioning(ctx,
		func(params hubapi.CapProvisioning_capManageProvisioning_Params) error {
			err2 := params.SetClientID(clientID)
			return err2
		})
	defer release()
	capability := getCap.Cap()
	newCap := NewManageProvisioningCapnpClient(capability.AddRef())
	return newCap
}

// CapRequestProvisioning provides the capability to provision IoT devices
func (cl *ProvisioningCapnpClient) CapRequestProvisioning(
	ctx context.Context, clientID string) provisioning.IRequestProvisioning {

	getCap, release := cl.capability.CapRequestProvisioning(ctx,
		func(params hubapi.CapProvisioning_capRequestProvisioning_Params) error {
			err2 := params.SetClientID(clientID)
			return err2
		})
	defer release()
	capability := getCap.Cap()
	newCap := NewRequestProvisioningCapnpClient(capability.AddRef())
	return newCap
}

// CapRefreshProvisioning provides the capability for IoT devices to refresh
func (cl *ProvisioningCapnpClient) CapRefreshProvisioning(
	ctx context.Context, clientID string) provisioning.IRefreshProvisioning {

	getCap, release := cl.capability.CapRefreshProvisioning(ctx,
		func(params hubapi.CapProvisioning_capRefreshProvisioning_Params) error {
			err2 := params.SetClientID(clientID)
			return err2
		})
	defer release()
	capability := getCap.Cap()
	newCap := NewRefreshProvisioningCapnpClient(capability.AddRef())
	return newCap
}

// Release the client capability
func (cl *ProvisioningCapnpClient) Release() {
	cl.capability.Release()
	cl.connection.Close()
}

// NewProvisioningCapnpClient returns a provisioning service client using the capnp protocol
//
//	ctx is the context for this client's connection. Release it to release the client.
//	conn is the connection with the provisioning capnp RPC server
func NewProvisioningCapnpClient(ctx context.Context, connection net.Conn) *ProvisioningCapnpClient {
	var cl *ProvisioningCapnpClient

	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)
	capability := hubapi.CapProvisioning(rpcConn.Bootstrap(ctx))

	cl = &ProvisioningCapnpClient{
		connection: rpcConn,
		capability: capability,
	}
	return cl
}
