package adapter

import (
	"context"
	"net"

	"capnproto.org/go/capnp/v3"

	"github.com/hiveot/hub.capnp/go/hubapi"
	"github.com/hiveot/hub/internal/caphelp"
	"github.com/hiveot/hub/pkg/directorystore/thingkvstore"
)

// DirectoryStoreCapnpAdapter for the directory store.
// This implements the capnproto generated interface DirectoryStore_Server
// See hub.capnp/go/hubapi/DirectoryStore.capnp.go for the interface.
type DirectoryStoreCapnpAdapter struct {
	store *thingkvstore.ThingKVStoreServer
}

func (adr *DirectoryStoreCapnpAdapter) GetTD(ctx context.Context, call hubapi.DirectoryStore_getTD) (err error) {
	var thingID string
	var td string

	args := call.Args()
	thingID, _ = args.ThingID()
	td, err = adr.store.GetTD(ctx, thingID)
	if err == nil {
		res, _ := call.AllocResults()
		err = res.SetTdJson(td)
	}
	return err
}

func (adr *DirectoryStoreCapnpAdapter) QueryTDs(ctx context.Context, call hubapi.DirectoryStore_queryTDs) (err error) {
	var jsonPath string
	var tdList []string

	args := call.Args()
	limit := args.Limit()
	offset := args.Offset()
	jsonPath, err = args.JsonPath()
	if err == nil {
		tdList, err = adr.store.QueryTDs(ctx, jsonPath, int(limit), int(offset), nil)
	}
	if err == nil {
		res, _ := call.AllocResults()
		textList, _ := res.NewTds(int32(len(tdList)))
		for i := 0; i < len(tdList); i++ {
			textList.Set(i, tdList[i])
		}
	}
	return err
}

func (adr *DirectoryStoreCapnpAdapter) ListTDs(ctx context.Context, call hubapi.DirectoryStore_listTDs) (err error) {
	var tdList []string

	args := call.Args()
	limit := args.Limit()
	offset := args.Offset()
	tdList, err = adr.store.ListTDs(ctx, int(limit), int(offset), nil)
	if err == nil {
		res, _ := call.AllocResults()
		textList, _ := res.NewTds(int32(len(tdList)))
		for i := 0; i < len(tdList); i++ {
			textList.Set(i, tdList[i])
		}
	}
	return err
}

func (adr *DirectoryStoreCapnpAdapter) UpdateTD(ctx context.Context, call hubapi.DirectoryStore_updateTD) (err error) {

	args := call.Args()
	thingID, _ := args.ThingID()
	tdDoc, _ := args.TdDoc()
	err = adr.store.UpdateTD(ctx, thingID, tdDoc)
	return err
}

// StartDirectoryStoreCapnpAdapter starts the directory store capnp protocol server
func StartDirectoryStoreCapnpAdapter(ctx context.Context,
	lis net.Listener,
	store *thingkvstore.ThingKVStoreServer) error {

	main := hubapi.DirectoryStore_ServerToClient(&DirectoryStoreCapnpAdapter{
		store: store,
	})
	return caphelp.CapServe(ctx, lis, capnp.Client(main))
}
