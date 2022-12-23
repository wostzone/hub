package service

import (
	"context"
	"fmt"
	"net"
	"sort"
	"sync"

	"capnproto.org/go/capnp/v3"
	"github.com/sirupsen/logrus"

	"github.com/hiveot/hub/pkg/resolver"
)

// ResolverService implements the IResolverService interface
type ResolverService struct {
	// connected sessions by sessionID
	sessions map[net.Conn]*ResolverSession

	// mutex for updating sessions
	sessionMutex sync.RWMutex
}

// OnIncomingConnection notifies the service of a new incoming connection.
// This is invoked by the underlying protocol and returns a new session to use
// with the connection.
// If this connection closes then capabilites added in this session are removed.
func (svc *ResolverService) OnIncomingConnection(conn net.Conn) resolver.IResolverSession {
	_ = conn
	svc.sessionMutex.Lock()
	defer svc.sessionMutex.Unlock()
	newSession := NewResolverSession(svc)
	svc.sessions[conn] = newSession
	return newSession
}

// OnConnectionClosed is invoked if the connection with the client has closed.
// The service will remove the session.
func (svc *ResolverService) OnConnectionClosed(conn net.Conn, session resolver.IResolverSession) {
	_ = conn
	// remove service when connection closes
	svc.sessionMutex.Lock()
	defer svc.sessionMutex.Unlock()
	for id, s := range svc.sessions {
		if s == session {
			delete(svc.sessions, id)
			break
		}
	}
}

// GetCapability returns the capability with the given name, if available.
// This method will return a 'future' interface for the service providing the capability.
// This won't detect a broken connection to the provider until the capability is used.
func (svc *ResolverService) GetCapability(ctx context.Context,
	clientID, clientType, capabilityName string, args []string) (
	capability capnp.Client, err error) {

	svc.sessionMutex.RLock()

	// determine which session this belongs to
	var capInfo resolver.CapabilityInfo
	var session *ResolverSession
	found := false
	for _, session = range svc.sessions {
		capList, err2 := session.ListRegisteredCapabilities(ctx)
		if err2 == nil {
			for _, capInfo = range capList {
				if capInfo.CapabilityName == capabilityName {
					found = true
					break
				}
			}
		}
		if found {
			break
		}
	}
	svc.sessionMutex.RUnlock()

	// unknown capability
	if capInfo.CapabilityName != capabilityName {
		err = fmt.Errorf("unknown capability '%s' requested for client '%s'", capabilityName, clientID)
		logrus.Warning(err)
		return capability, err
	}

	capability, err = session.GetRegisteredCapability(ctx, clientID, clientType, capabilityName, args)

	return capability, err
}

// ListCapabilities returns list of capabilities of all connected services sorted by service and capability names
// This also verifies the connections and removes capabilities that are no longer valid.
func (svc *ResolverService) ListCapabilities(ctx context.Context) ([]resolver.CapabilityInfo, error) {

	capList := make([]resolver.CapabilityInfo, 0)

	svc.sessionMutex.RLock()
	defer svc.sessionMutex.RUnlock()

	for _, session := range svc.sessions {
		sessionCaps, err := session.ListRegisteredCapabilities(ctx)
		if err == nil {
			capList = append(capList, sessionCaps...)
		}
	}
	//logrus.Infof("listing '%d' capabilities from %d sessions", len(capList), len(svc.sessions))
	// sort by service ID + capability Name
	sort.Slice(capList, func(i, j int) bool {
		iName := capList[i].ServiceID + capList[i].CapabilityName
		jName := capList[j].ServiceID + capList[j].CapabilityName
		return iName < jName
	})
	return capList, nil
}

// Start currently has nothing to do as the capnpserver listens for incoming connections
func (svc *ResolverService) Start(_ context.Context) error {
	//logrus.Infof("Starting resolver service")
	return nil
}

// Stop closes all remaining sessions
func (svc *ResolverService) Stop() (err error) {

	svc.sessionMutex.RLock()
	sessionIDList := make([]net.Conn, 0, len(svc.sessions))
	for conn := range svc.sessions {
		sessionIDList = append(sessionIDList, conn)
	}
	svc.sessionMutex.RUnlock()

	logrus.Infof("Stopping resolver service. %d sessions remaining", len(sessionIDList))

	for _, sessionID := range sessionIDList {
		svc.sessionMutex.Lock()
		session := svc.sessions[sessionID]
		delete(svc.sessions, sessionID)
		session.Release()
		svc.sessionMutex.Unlock()
	}
	return err
}

// NewResolverService returns a new instance of the capability resolver
func NewResolverService() *ResolverService {
	svc := &ResolverService{
		sessions: make(map[net.Conn]*ResolverSession),
		//sessionMutex: sync.RWMutex{},
	}
	return svc
}
