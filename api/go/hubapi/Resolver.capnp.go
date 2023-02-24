// Code generated by capnpc-go. DO NOT EDIT.

package hubapi

import (
	capnp "capnproto.org/go/capnp/v3"
	text "capnproto.org/go/capnp/v3/encoding/text"
	fc "capnproto.org/go/capnp/v3/flowcontrol"
	schemas "capnproto.org/go/capnp/v3/schemas"
	server "capnproto.org/go/capnp/v3/server"
	context "context"
	fmt "fmt"
)

const AuthType = uint64(0xf41ce091c8668088)

// Constants defined in Resolver.capnp.
const (
	ResolverServiceName     = "resolver"
	DefaultResolverAddress  = "/tmp/hiveot-resolver.socket"
	AuthTypeUnauthenticated = "unauthenticated"
	AuthTypeAdmin           = "admin"
	AuthTypeIotDevice       = "iotdevice"
	AuthTypeUser            = "user"
	AuthTypeService         = "service"
)

type CapabilityInfo capnp.Struct

// CapabilityInfo_TypeID is the unique identifier for the type CapabilityInfo.
const CapabilityInfo_TypeID = 0xae50171d46b6aaf1

func NewCapabilityInfo(s *capnp.Segment) (CapabilityInfo, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7})
	return CapabilityInfo(st), err
}

func NewRootCapabilityInfo(s *capnp.Segment) (CapabilityInfo, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7})
	return CapabilityInfo(st), err
}

func ReadRootCapabilityInfo(msg *capnp.Message) (CapabilityInfo, error) {
	root, err := msg.Root()
	return CapabilityInfo(root.Struct()), err
}

func (s CapabilityInfo) String() string {
	str, _ := text.Marshal(0xae50171d46b6aaf1, capnp.Struct(s))
	return str
}

func (s CapabilityInfo) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (CapabilityInfo) DecodeFromPtr(p capnp.Ptr) CapabilityInfo {
	return CapabilityInfo(capnp.Struct{}.DecodeFromPtr(p))
}

func (s CapabilityInfo) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s CapabilityInfo) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s CapabilityInfo) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s CapabilityInfo) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s CapabilityInfo) InterfaceID() uint64 {
	return capnp.Struct(s).Uint64(0)
}

func (s CapabilityInfo) SetInterfaceID(v uint64) {
	capnp.Struct(s).SetUint64(0, v)
}

func (s CapabilityInfo) MethodID() uint16 {
	return capnp.Struct(s).Uint16(8)
}

func (s CapabilityInfo) SetMethodID(v uint16) {
	capnp.Struct(s).SetUint16(8, v)
}

func (s CapabilityInfo) InterfaceName() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s CapabilityInfo) HasInterfaceName() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s CapabilityInfo) InterfaceNameBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetInterfaceName(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

func (s CapabilityInfo) MethodName() (string, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.Text(), err
}

func (s CapabilityInfo) HasMethodName() bool {
	return capnp.Struct(s).HasPtr(1)
}

func (s CapabilityInfo) MethodNameBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(1)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetMethodName(v string) error {
	return capnp.Struct(s).SetText(1, v)
}

func (s CapabilityInfo) AuthTypes() (capnp.TextList, error) {
	p, err := capnp.Struct(s).Ptr(2)
	return capnp.TextList(p.List()), err
}

func (s CapabilityInfo) HasAuthTypes() bool {
	return capnp.Struct(s).HasPtr(2)
}

func (s CapabilityInfo) SetAuthTypes(v capnp.TextList) error {
	return capnp.Struct(s).SetPtr(2, v.ToPtr())
}

// NewAuthTypes sets the authTypes field to a newly
// allocated capnp.TextList, preferring placement in s's segment.
func (s CapabilityInfo) NewAuthTypes(n int32) (capnp.TextList, error) {
	l, err := capnp.NewTextList(capnp.Struct(s).Segment(), n)
	if err != nil {
		return capnp.TextList{}, err
	}
	err = capnp.Struct(s).SetPtr(2, l.ToPtr())
	return l, err
}
func (s CapabilityInfo) Protocol() (string, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return p.Text(), err
}

func (s CapabilityInfo) HasProtocol() bool {
	return capnp.Struct(s).HasPtr(3)
}

func (s CapabilityInfo) ProtocolBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(3)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetProtocol(v string) error {
	return capnp.Struct(s).SetText(3, v)
}

func (s CapabilityInfo) ServiceID() (string, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.Text(), err
}

func (s CapabilityInfo) HasServiceID() bool {
	return capnp.Struct(s).HasPtr(4)
}

func (s CapabilityInfo) ServiceIDBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(4)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetServiceID(v string) error {
	return capnp.Struct(s).SetText(4, v)
}

func (s CapabilityInfo) Network() (string, error) {
	p, err := capnp.Struct(s).Ptr(5)
	return p.Text(), err
}

func (s CapabilityInfo) HasNetwork() bool {
	return capnp.Struct(s).HasPtr(5)
}

func (s CapabilityInfo) NetworkBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(5)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetNetwork(v string) error {
	return capnp.Struct(s).SetText(5, v)
}

func (s CapabilityInfo) Address() (string, error) {
	p, err := capnp.Struct(s).Ptr(6)
	return p.Text(), err
}

func (s CapabilityInfo) HasAddress() bool {
	return capnp.Struct(s).HasPtr(6)
}

func (s CapabilityInfo) AddressBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(6)
	return p.TextBytes(), err
}

func (s CapabilityInfo) SetAddress(v string) error {
	return capnp.Struct(s).SetText(6, v)
}

// CapabilityInfo_List is a list of CapabilityInfo.
type CapabilityInfo_List = capnp.StructList[CapabilityInfo]

// NewCapabilityInfo creates a new list of CapabilityInfo.
func NewCapabilityInfo_List(s *capnp.Segment, sz int32) (CapabilityInfo_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 16, PointerCount: 7}, sz)
	return capnp.StructList[CapabilityInfo](l), err
}

// CapabilityInfo_Future is a wrapper for a CapabilityInfo promised by a client call.
type CapabilityInfo_Future struct{ *capnp.Future }

func (f CapabilityInfo_Future) Struct() (CapabilityInfo, error) {
	p, err := f.Future.Ptr()
	return CapabilityInfo(p.Struct()), err
}

type CapResolverService capnp.Client

// CapResolverService_TypeID is the unique identifier for the type CapResolverService.
const CapResolverService_TypeID = 0xab76eb2c88343a05

func (c CapResolverService) ListCapabilities(ctx context.Context, params func(CapProvider_listCapabilities_Params) error) (CapProvider_listCapabilities_Results_Future, capnp.ReleaseFunc) {
	s := capnp.Send{
		Method: capnp.Method{
			InterfaceID:   0xacf14758b95cf892,
			MethodID:      0,
			InterfaceName: "hubapi/Resolver.capnp:CapProvider",
			MethodName:    "listCapabilities",
		},
	}
	if params != nil {
		s.ArgsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		s.PlaceArgs = func(s capnp.Struct) error { return params(CapProvider_listCapabilities_Params(s)) }
	}
	ans, release := capnp.Client(c).SendCall(ctx, s)
	return CapProvider_listCapabilities_Results_Future{Future: ans.Future()}, release
}

// String returns a string that identifies this capability for debugging
// purposes.  Its format should not be depended on: in particular, it
// should not be used to compare clients.  Use IsSame to compare clients
// for equality.
func (c CapResolverService) String() string {
	return fmt.Sprintf("%T(%v)", c, capnp.Client(c))
}

// AddRef creates a new Client that refers to the same capability as c.
// If c is nil or has resolved to null, then AddRef returns nil.
func (c CapResolverService) AddRef() CapResolverService {
	return CapResolverService(capnp.Client(c).AddRef())
}

// Release releases a capability reference.  If this is the last
// reference to the capability, then the underlying resources associated
// with the capability will be released.
//
// Release will panic if c has already been released, but not if c is
// nil or resolved to null.
func (c CapResolverService) Release() {
	capnp.Client(c).Release()
}

// Resolve blocks until the capability is fully resolved or the Context
// expires.
func (c CapResolverService) Resolve(ctx context.Context) error {
	return capnp.Client(c).Resolve(ctx)
}

func (c CapResolverService) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Client(c).EncodeAsPtr(seg)
}

func (CapResolverService) DecodeFromPtr(p capnp.Ptr) CapResolverService {
	return CapResolverService(capnp.Client{}.DecodeFromPtr(p))
}

// IsValid reports whether c is a valid reference to a capability.
// A reference is invalid if it is nil, has resolved to null, or has
// been released.
func (c CapResolverService) IsValid() bool {
	return capnp.Client(c).IsValid()
}

// IsSame reports whether c and other refer to a capability created by the
// same call to NewClient.  This can return false negatives if c or other
// are not fully resolved: use Resolve if this is an issue.  If either
// c or other are released, then IsSame panics.
func (c CapResolverService) IsSame(other CapResolverService) bool {
	return capnp.Client(c).IsSame(capnp.Client(other))
}

// Update the flowcontrol.FlowLimiter used to manage flow control for
// this client. This affects all future calls, but not calls already
// waiting to send. Passing nil sets the value to flowcontrol.NopLimiter,
// which is also the default.
func (c CapResolverService) SetFlowLimiter(lim fc.FlowLimiter) {
	capnp.Client(c).SetFlowLimiter(lim)
}

// Get the current flowcontrol.FlowLimiter used to manage flow control
// for this client.
func (c CapResolverService) GetFlowLimiter() fc.FlowLimiter {
	return capnp.Client(c).GetFlowLimiter()
} // A CapResolverService_Server is a CapResolverService with a local implementation.
type CapResolverService_Server interface {
	ListCapabilities(context.Context, CapProvider_listCapabilities) error
}

// CapResolverService_NewServer creates a new Server from an implementation of CapResolverService_Server.
func CapResolverService_NewServer(s CapResolverService_Server) *server.Server {
	c, _ := s.(server.Shutdowner)
	return server.New(CapResolverService_Methods(nil, s), s, c)
}

// CapResolverService_ServerToClient creates a new Client from an implementation of CapResolverService_Server.
// The caller is responsible for calling Release on the returned Client.
func CapResolverService_ServerToClient(s CapResolverService_Server) CapResolverService {
	return CapResolverService(capnp.NewClient(CapResolverService_NewServer(s)))
}

// CapResolverService_Methods appends Methods to a slice that invoke the methods on s.
// This can be used to create a more complicated Server.
func CapResolverService_Methods(methods []server.Method, s CapResolverService_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xacf14758b95cf892,
			MethodID:      0,
			InterfaceName: "hubapi/Resolver.capnp:CapProvider",
			MethodName:    "listCapabilities",
		},
		Impl: func(ctx context.Context, call *server.Call) error {
			return s.ListCapabilities(ctx, CapProvider_listCapabilities{call})
		},
	})

	return methods
}

// CapResolverService_List is a list of CapResolverService.
type CapResolverService_List = capnp.CapList[CapResolverService]

// NewCapResolverService creates a new list of CapResolverService.
func NewCapResolverService_List(s *capnp.Segment, sz int32) (CapResolverService_List, error) {
	l, err := capnp.NewPointerList(s, sz)
	return capnp.CapList[CapResolverService](l), err
}

type CapProvider capnp.Client

// CapProvider_TypeID is the unique identifier for the type CapProvider.
const CapProvider_TypeID = 0xacf14758b95cf892

func (c CapProvider) ListCapabilities(ctx context.Context, params func(CapProvider_listCapabilities_Params) error) (CapProvider_listCapabilities_Results_Future, capnp.ReleaseFunc) {
	s := capnp.Send{
		Method: capnp.Method{
			InterfaceID:   0xacf14758b95cf892,
			MethodID:      0,
			InterfaceName: "hubapi/Resolver.capnp:CapProvider",
			MethodName:    "listCapabilities",
		},
	}
	if params != nil {
		s.ArgsSize = capnp.ObjectSize{DataSize: 0, PointerCount: 1}
		s.PlaceArgs = func(s capnp.Struct) error { return params(CapProvider_listCapabilities_Params(s)) }
	}
	ans, release := capnp.Client(c).SendCall(ctx, s)
	return CapProvider_listCapabilities_Results_Future{Future: ans.Future()}, release
}

// String returns a string that identifies this capability for debugging
// purposes.  Its format should not be depended on: in particular, it
// should not be used to compare clients.  Use IsSame to compare clients
// for equality.
func (c CapProvider) String() string {
	return fmt.Sprintf("%T(%v)", c, capnp.Client(c))
}

// AddRef creates a new Client that refers to the same capability as c.
// If c is nil or has resolved to null, then AddRef returns nil.
func (c CapProvider) AddRef() CapProvider {
	return CapProvider(capnp.Client(c).AddRef())
}

// Release releases a capability reference.  If this is the last
// reference to the capability, then the underlying resources associated
// with the capability will be released.
//
// Release will panic if c has already been released, but not if c is
// nil or resolved to null.
func (c CapProvider) Release() {
	capnp.Client(c).Release()
}

// Resolve blocks until the capability is fully resolved or the Context
// expires.
func (c CapProvider) Resolve(ctx context.Context) error {
	return capnp.Client(c).Resolve(ctx)
}

func (c CapProvider) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Client(c).EncodeAsPtr(seg)
}

func (CapProvider) DecodeFromPtr(p capnp.Ptr) CapProvider {
	return CapProvider(capnp.Client{}.DecodeFromPtr(p))
}

// IsValid reports whether c is a valid reference to a capability.
// A reference is invalid if it is nil, has resolved to null, or has
// been released.
func (c CapProvider) IsValid() bool {
	return capnp.Client(c).IsValid()
}

// IsSame reports whether c and other refer to a capability created by the
// same call to NewClient.  This can return false negatives if c or other
// are not fully resolved: use Resolve if this is an issue.  If either
// c or other are released, then IsSame panics.
func (c CapProvider) IsSame(other CapProvider) bool {
	return capnp.Client(c).IsSame(capnp.Client(other))
}

// Update the flowcontrol.FlowLimiter used to manage flow control for
// this client. This affects all future calls, but not calls already
// waiting to send. Passing nil sets the value to flowcontrol.NopLimiter,
// which is also the default.
func (c CapProvider) SetFlowLimiter(lim fc.FlowLimiter) {
	capnp.Client(c).SetFlowLimiter(lim)
}

// Get the current flowcontrol.FlowLimiter used to manage flow control
// for this client.
func (c CapProvider) GetFlowLimiter() fc.FlowLimiter {
	return capnp.Client(c).GetFlowLimiter()
} // A CapProvider_Server is a CapProvider with a local implementation.
type CapProvider_Server interface {
	ListCapabilities(context.Context, CapProvider_listCapabilities) error
}

// CapProvider_NewServer creates a new Server from an implementation of CapProvider_Server.
func CapProvider_NewServer(s CapProvider_Server) *server.Server {
	c, _ := s.(server.Shutdowner)
	return server.New(CapProvider_Methods(nil, s), s, c)
}

// CapProvider_ServerToClient creates a new Client from an implementation of CapProvider_Server.
// The caller is responsible for calling Release on the returned Client.
func CapProvider_ServerToClient(s CapProvider_Server) CapProvider {
	return CapProvider(capnp.NewClient(CapProvider_NewServer(s)))
}

// CapProvider_Methods appends Methods to a slice that invoke the methods on s.
// This can be used to create a more complicated Server.
func CapProvider_Methods(methods []server.Method, s CapProvider_Server) []server.Method {
	if cap(methods) == 0 {
		methods = make([]server.Method, 0, 1)
	}

	methods = append(methods, server.Method{
		Method: capnp.Method{
			InterfaceID:   0xacf14758b95cf892,
			MethodID:      0,
			InterfaceName: "hubapi/Resolver.capnp:CapProvider",
			MethodName:    "listCapabilities",
		},
		Impl: func(ctx context.Context, call *server.Call) error {
			return s.ListCapabilities(ctx, CapProvider_listCapabilities{call})
		},
	})

	return methods
}

// CapProvider_listCapabilities holds the state for a server call to CapProvider.listCapabilities.
// See server.Call for documentation.
type CapProvider_listCapabilities struct {
	*server.Call
}

// Args returns the call's arguments.
func (c CapProvider_listCapabilities) Args() CapProvider_listCapabilities_Params {
	return CapProvider_listCapabilities_Params(c.Call.Args())
}

// AllocResults allocates the results struct.
func (c CapProvider_listCapabilities) AllocResults() (CapProvider_listCapabilities_Results, error) {
	r, err := c.Call.AllocResults(capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return CapProvider_listCapabilities_Results(r), err
}

// CapProvider_List is a list of CapProvider.
type CapProvider_List = capnp.CapList[CapProvider]

// NewCapProvider creates a new list of CapProvider.
func NewCapProvider_List(s *capnp.Segment, sz int32) (CapProvider_List, error) {
	l, err := capnp.NewPointerList(s, sz)
	return capnp.CapList[CapProvider](l), err
}

type CapProvider_listCapabilities_Params capnp.Struct

// CapProvider_listCapabilities_Params_TypeID is the unique identifier for the type CapProvider_listCapabilities_Params.
const CapProvider_listCapabilities_Params_TypeID = 0xbb0d5b68acfa1d84

func NewCapProvider_listCapabilities_Params(s *capnp.Segment) (CapProvider_listCapabilities_Params, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return CapProvider_listCapabilities_Params(st), err
}

func NewRootCapProvider_listCapabilities_Params(s *capnp.Segment) (CapProvider_listCapabilities_Params, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return CapProvider_listCapabilities_Params(st), err
}

func ReadRootCapProvider_listCapabilities_Params(msg *capnp.Message) (CapProvider_listCapabilities_Params, error) {
	root, err := msg.Root()
	return CapProvider_listCapabilities_Params(root.Struct()), err
}

func (s CapProvider_listCapabilities_Params) String() string {
	str, _ := text.Marshal(0xbb0d5b68acfa1d84, capnp.Struct(s))
	return str
}

func (s CapProvider_listCapabilities_Params) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (CapProvider_listCapabilities_Params) DecodeFromPtr(p capnp.Ptr) CapProvider_listCapabilities_Params {
	return CapProvider_listCapabilities_Params(capnp.Struct{}.DecodeFromPtr(p))
}

func (s CapProvider_listCapabilities_Params) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s CapProvider_listCapabilities_Params) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s CapProvider_listCapabilities_Params) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s CapProvider_listCapabilities_Params) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s CapProvider_listCapabilities_Params) AuthType() (string, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.Text(), err
}

func (s CapProvider_listCapabilities_Params) HasAuthType() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s CapProvider_listCapabilities_Params) AuthTypeBytes() ([]byte, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return p.TextBytes(), err
}

func (s CapProvider_listCapabilities_Params) SetAuthType(v string) error {
	return capnp.Struct(s).SetText(0, v)
}

// CapProvider_listCapabilities_Params_List is a list of CapProvider_listCapabilities_Params.
type CapProvider_listCapabilities_Params_List = capnp.StructList[CapProvider_listCapabilities_Params]

// NewCapProvider_listCapabilities_Params creates a new list of CapProvider_listCapabilities_Params.
func NewCapProvider_listCapabilities_Params_List(s *capnp.Segment, sz int32) (CapProvider_listCapabilities_Params_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return capnp.StructList[CapProvider_listCapabilities_Params](l), err
}

// CapProvider_listCapabilities_Params_Future is a wrapper for a CapProvider_listCapabilities_Params promised by a client call.
type CapProvider_listCapabilities_Params_Future struct{ *capnp.Future }

func (f CapProvider_listCapabilities_Params_Future) Struct() (CapProvider_listCapabilities_Params, error) {
	p, err := f.Future.Ptr()
	return CapProvider_listCapabilities_Params(p.Struct()), err
}

type CapProvider_listCapabilities_Results capnp.Struct

// CapProvider_listCapabilities_Results_TypeID is the unique identifier for the type CapProvider_listCapabilities_Results.
const CapProvider_listCapabilities_Results_TypeID = 0xf5cba8f2960769af

func NewCapProvider_listCapabilities_Results(s *capnp.Segment) (CapProvider_listCapabilities_Results, error) {
	st, err := capnp.NewStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return CapProvider_listCapabilities_Results(st), err
}

func NewRootCapProvider_listCapabilities_Results(s *capnp.Segment) (CapProvider_listCapabilities_Results, error) {
	st, err := capnp.NewRootStruct(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1})
	return CapProvider_listCapabilities_Results(st), err
}

func ReadRootCapProvider_listCapabilities_Results(msg *capnp.Message) (CapProvider_listCapabilities_Results, error) {
	root, err := msg.Root()
	return CapProvider_listCapabilities_Results(root.Struct()), err
}

func (s CapProvider_listCapabilities_Results) String() string {
	str, _ := text.Marshal(0xf5cba8f2960769af, capnp.Struct(s))
	return str
}

func (s CapProvider_listCapabilities_Results) EncodeAsPtr(seg *capnp.Segment) capnp.Ptr {
	return capnp.Struct(s).EncodeAsPtr(seg)
}

func (CapProvider_listCapabilities_Results) DecodeFromPtr(p capnp.Ptr) CapProvider_listCapabilities_Results {
	return CapProvider_listCapabilities_Results(capnp.Struct{}.DecodeFromPtr(p))
}

func (s CapProvider_listCapabilities_Results) ToPtr() capnp.Ptr {
	return capnp.Struct(s).ToPtr()
}
func (s CapProvider_listCapabilities_Results) IsValid() bool {
	return capnp.Struct(s).IsValid()
}

func (s CapProvider_listCapabilities_Results) Message() *capnp.Message {
	return capnp.Struct(s).Message()
}

func (s CapProvider_listCapabilities_Results) Segment() *capnp.Segment {
	return capnp.Struct(s).Segment()
}
func (s CapProvider_listCapabilities_Results) InfoList() (CapabilityInfo_List, error) {
	p, err := capnp.Struct(s).Ptr(0)
	return CapabilityInfo_List(p.List()), err
}

func (s CapProvider_listCapabilities_Results) HasInfoList() bool {
	return capnp.Struct(s).HasPtr(0)
}

func (s CapProvider_listCapabilities_Results) SetInfoList(v CapabilityInfo_List) error {
	return capnp.Struct(s).SetPtr(0, v.ToPtr())
}

// NewInfoList sets the infoList field to a newly
// allocated CapabilityInfo_List, preferring placement in s's segment.
func (s CapProvider_listCapabilities_Results) NewInfoList(n int32) (CapabilityInfo_List, error) {
	l, err := NewCapabilityInfo_List(capnp.Struct(s).Segment(), n)
	if err != nil {
		return CapabilityInfo_List{}, err
	}
	err = capnp.Struct(s).SetPtr(0, l.ToPtr())
	return l, err
}

// CapProvider_listCapabilities_Results_List is a list of CapProvider_listCapabilities_Results.
type CapProvider_listCapabilities_Results_List = capnp.StructList[CapProvider_listCapabilities_Results]

// NewCapProvider_listCapabilities_Results creates a new list of CapProvider_listCapabilities_Results.
func NewCapProvider_listCapabilities_Results_List(s *capnp.Segment, sz int32) (CapProvider_listCapabilities_Results_List, error) {
	l, err := capnp.NewCompositeList(s, capnp.ObjectSize{DataSize: 0, PointerCount: 1}, sz)
	return capnp.StructList[CapProvider_listCapabilities_Results](l), err
}

// CapProvider_listCapabilities_Results_Future is a wrapper for a CapProvider_listCapabilities_Results promised by a client call.
type CapProvider_listCapabilities_Results_Future struct{ *capnp.Future }

func (f CapProvider_listCapabilities_Results_Future) Struct() (CapProvider_listCapabilities_Results, error) {
	p, err := f.Future.Ptr()
	return CapProvider_listCapabilities_Results(p.Struct()), err
}

const schema_f02d0b8fc1fe2004 = "x\xda\x9c\x95]h\x1c\xd5\x1b\xc6\xdf\xf7\x9c93\xe9" +
	"\xff\x9fm:\x9e%M\xa51Z\x13\x0c%_M*" +
	"h\x11\x9a\xa4i\xc3\x06);;\x11\xad\x8a8\xd99" +
	"\xcbN\xb3\xbb\xb3\xce\xcc\xae\xf6\xaaR\x8a\x14o\xb4E" +
	"\x04\xef\xbcP\xf0#h\x14*\x88U\xa8\xf4\xa2E\xe8" +
	"\x85 \x827\x8a\xd6\x8b\xa2B+\x0d\xb4\xc5t\xe4\xcc" +
	"\xeelv\xdd\xac\x05\xef2\xcf\xf9=\xcf9\xe7=\x0f" +
	"\xd9\x89\x07\xc9\xb4\xb2'\xf1#\x03b\xcc15\xfc~" +
	"O\xe5\x7f\x07'\x97\xce\x80\xde\xab\x84\xca\xfdw\xce\xbf" +
	"\xf6\xff\xd1k\x008\xf5\x1d\xc9 \xbfJ4\x00\xf3\x0a" +
	"\xa1h^#\x04\x01B\xb6o\xef\xa9\x91\xdf\xaa\x1f\x82" +
	"\xdeK[\xf8\xabd\x01\xf9\xba\xe4\xf9-\xa2\xf1[d" +
	";@x\xe6\xe6\xb3\x9f?5\x7f}e\x13\xfa\x1e\x94" +
	"\x1c\x00_#\xf3|\x88j\x00\xe1\xf5\x0f>;\xd4\xbf" +
	"=\xfd\x11\x18\xbdH6p\xa6i\x00S\x09\xba\x1by" +
	"\xbf\xe4\xa6v\xd0\x10\x01\xc3\xd5\x1d}\xbf\\N\xe9\x9f" +
	"\xb6\x1d\xfd ;\x8a\xfc\x08\x93G_d\x14\xcd\xe7Y" +
	"t\xf4\x93\xfd\xb7W\xf2\xcf$\xbe\x00}\x17\x020\x94" +
	"Q\x16\xfb\x16\x01y\x85\xed\x07\x0cO\xbc4T\x14\xc3" +
	"\xf3\x17\xda\x02\xdfb\xf7\"\x7f?\x0a|G\x06\xae\xd6" +
	"\x02\x17\xdf>\x9a\xfd\xea\xa1W\xae\xb4\xf1\xef\xb2]\xc8" +
	"\xcfF\xfc\xaa\xe4\xcf\xd5\xf8#\xf7\xcd\x9f|}\xf9\xd0" +
	"\xafm\xfc'l\x12\xf9\xf9\x88?'\xf9\x8b5\xfe\xd4" +
	"\xcb\xb9\x8b\xa7\x7f\xday\x03N\xf72\xd2l\xe0_\xb2" +
	"\xdb\xfcR\xc4_`T\xce\xe2cG{\xf3\xcf\xf7\xbe" +
	"Yk\xbe\xdbY\xf6\x83\xbc\xdb\xa5\xe8n/\xcc8\xcf" +
	"\xbd\xb1rg\xbdm\xefu\xe6!\xd7U\x99\xd5\xadR" +
	"4\xfb\xd4h\xef\xaf\x9fL\x8d]^\x1d\x0f\xdbx]" +
	"\x9dE\xfe@\xc4\xef\x94\xfc\xb0\xe4G\xc2|e\xc9*" +
	";\xe3\x19E\xf8n\xa1*\xbc\xb1\xacU.\x95\xf7y" +
	"\xf5OSxU'+\x0e[E\x14iD\xec\x06\x82" +
	"\xdd\x00:.\x841\x02\x00\x9dR\x0eX\xe5Lk\x10" +
	"@\x1a1M\x99\xa1 n\xb4\xac)\x81\xb6%\xa4=" +
	"\xb7\xea\xd8\x02\xbd4\xa2\xa1P\xd6\xd4\x08\x8c\xc7\xa7\xeb" +
	"\xaf\x02\xd1\x13ZXp\xfc\xe0\x80U\xb6p\xc9)8" +
	"\x81#|\x80iL#\xfeK\xbe\x15\xa1\xc7zR\xa5" +
	"\x9c+\xb7\x18\xa4\x0a\x80\x82\x00\xfa\x1fK\x00\xc6\xef\x14" +
	"\x8d\x9b\x04uT\x92(\xc5\xb5\x05\x00\xe3\x06\xc5\x0c\x12" +
	"D\x92D\x02\xa0\xaf{\x00\xc6_\x14\xcd.$\xa8S" +
	"L\"\x05\xe0\x0c\x9f\x060\x15\xa4hn\x93\xbaB\x92" +
	"\xa8\x00\xf0\x04f\xe4\xa3I}X\xea\x8c&\x91\x01\xf0" +
	"!\\\x000\x07\xa5>!uUI\xa2\x0a\xc0G#" +
	"~D\xea\x8fH]cI\xd9\x13\xfe0\xce\x02\x98\x13" +
	"R\x7fL\xea]j\x12\xbb\x00\xf8\xa3\x91\xbeW\xea\xd3" +
	"H0tJ\x81\xf0rV\x164\x91\x9a\xc3-@p" +
	"\x0b`X\x14A\xde\xb5Ss\x00\x80\x1a\x10\xd4\xa0\x89" +
	"\x1c\x90/.\xea\xaf\x1d\xb3\x87-\xa0M\xa2U\x09\xf2" +
	"\x8b\xc7\xca\x02\xd0\xc7\xad\x80iZk\xc7V\xc0\xb0\xec" +
	"\xb9\x81\x9bu\x0b2;\xc6\xfdZ\x01R\x80s\xb1v" +
	"\xbc$\x82\x17]o\xb9\xf1m\xd9\xb6'|\xbf\xe1\xe9" +
	"P+[\xe4\xacJ!\x88\xab5c\xdb=\xd2\xd6\xd2" +
	"\xcf\x9f\xc3\xf1\xa0X\x1e\xcf;U\"\xdc`4.\xeb" +
	"\x98\xeff\xb5e\x114\xa2\xd5N}\xf3\xc6\xe2.\xc5" +
	"U\x1aL[\x9eUD\xdfP\x1a\x0dI\xc82tS" +
	"4\xfaH\xd3@\x9an\xdd\xa1v1\xfa\x84O\x85\xd7" +
	"r\xee\xdd=\x15_xw\xf3\xcd\xd8Z\xd1)\xb5\x18" +
	"'\x07,\xbb\xe8\x94\xee\xe64\xc5@\xf4\x0e-\xde\xd9" +
	"\xe3\xf5\xd7i\xb8\xc9?\xdd\xfbk\xf6\xba\xed?\x8d/" +
	"#\xfcJ\x81\x06\x9b\xceo8\xeai\xce}\xdc\xf1\x03" +
	"9\xbfz\xa1\xb6m\xfc\xb6\x00F\xd5\xeaP\x89\xc6@" +
	"K\xf2/Q\x0a\x9c\x81\xac\x15\x08\xbb\xe5\x9a'\xc2J" +
	"}Y\x0d\x9ch\x19\xe2\xb8\xcd\xd3Rn0'\xa2\x7f" +
	"[-9\x99\xd0q\x03[.\x00\x8a\xbf\x03\x00\x00\xff" +
	"\xffy\xb6\x16\x8b"

func init() {
	schemas.Register(schema_f02d0b8fc1fe2004,
		0x926232450a7531d7,
		0xab76eb2c88343a05,
		0xacf14758b95cf892,
		0xae50171d46b6aaf1,
		0xb21149cee31819b0,
		0xbb0d5b68acfa1d84,
		0xc44728656d257882,
		0xe48627be636aa054,
		0xe5466b9084471e59,
		0xf41ce091c8668088,
		0xf5cba8f2960769af,
		0xfdfeac945e694171,
		0xff2fb0ce2e4957c2)
}
