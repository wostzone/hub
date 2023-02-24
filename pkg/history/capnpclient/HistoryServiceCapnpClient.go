// Package capnpclient that wraps the capnp generated client with a POGS API
package capnpclient

import (
	"context"
	"net"

	"capnproto.org/go/capnp/v3/rpc"

	"github.com/hiveot/hub/api/go/hubapi"
	"github.com/hiveot/hub/pkg/history"
)

// HistoryServiceCapnpClient provides a POGS wrapper around the capnp client API
// This implements the IHistoryService interface
type HistoryServiceCapnpClient struct {
	connection *rpc.Conn                // connection to capnp server
	capability hubapi.CapHistoryService // capnp client
}

func (cl *HistoryServiceCapnpClient) CapAddHistory(
	ctx context.Context, clientID string, ignoreRetention bool) (history.IAddHistory, error) {

	getCap, release := cl.capability.CapAddHistory(ctx,
		func(params hubapi.CapHistoryService_capAddHistory_Params) error {
			err2 := params.SetClientID(clientID)
			params.SetIgnoreRetention(ignoreRetention)
			return err2
		})
	defer release()
	capability := getCap.Cap().AddRef()
	// reuse the add history capability
	newCap := NewAddHistoryCapnpClient(capability)
	return newCap, nil
}

//
//// CapAddHistory provides the capability to add to the history
//func (cl *HistoryServiceCapnpClient) CapAddHistory(
//	ctx context.Context, clientID string, publisherID, thingID string) (history.IAddHistory, error) {
//
//	// The use of a result 'future' avoids a round trip, making this more efficient
//	getCap, release := cl.capability.CapAddHistory(ctx,
//		func(params hubapi.CapHistoryService_capAddHistory_Params) error {
//			err2 := params.SetClientID(clientID)
//			_ = params.SetPublisherID(publisherID)
//			_ = params.SetThingID(thingID)
//			return err2
//		})
//
//	defer release()
//	capability := getCap.Cap().AddRef()
//
//	newCap := NewAddHistoryCapnpClient(capability)
//	return newCap, nil
//}

func (cl *HistoryServiceCapnpClient) CapManageRetention(
	ctx context.Context, clientID string) (history.IManageRetention, error) {

	getCap, release := cl.capability.CapManageRetention(ctx,
		func(params hubapi.CapHistoryService_capManageRetention_Params) error {
			err2 := params.SetClientID(clientID)
			return err2
		})
	defer release()
	capability := getCap.Cap().AddRef()

	newCap := NewManageRetentionCapnpClient(capability)
	return newCap, nil
}

// CapReadHistory the capability to iterate the history
func (cl *HistoryServiceCapnpClient) CapReadHistory(
	ctx context.Context, clientID string, publisherID, thingID string) (history.IReadHistory, error) {

	getCap, release := cl.capability.CapReadHistory(ctx,
		func(params hubapi.CapHistoryService_capReadHistory_Params) error {
			err2 := params.SetClientID(clientID)
			_ = params.SetThingID(thingID)
			_ = params.SetPublisherID(publisherID)
			return err2
		})
	defer release()
	capability := getCap.Cap().AddRef()

	newCap := NewReadHistoryCapnpClient(capability)
	return newCap, nil
}

func (cl *HistoryServiceCapnpClient) Release() {
	cl.capability.Release()
}

// NewHistoryCapnpClient returns a history service client using the capnp protocol.
// This implements the IHistoryService interface.
//
//	ctx is the context for getting capabilities from the server
//	connection is the connection to the capnp server
func NewHistoryCapnpClient(ctx context.Context, connection net.Conn) *HistoryServiceCapnpClient {
	var cl *HistoryServiceCapnpClient
	transport := rpc.NewStreamTransport(connection)
	rpcConn := rpc.NewConn(transport, nil)
	capability := hubapi.CapHistoryService(rpcConn.Bootstrap(ctx))

	cl = &HistoryServiceCapnpClient{
		connection: rpcConn,
		capability: capability,
	}
	return cl
}
