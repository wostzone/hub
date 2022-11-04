package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/hiveot/hub.go/pkg/logging"
	"github.com/hiveot/hub/internal/listener"
	"github.com/hiveot/hub/internal/svcconfig"
	"github.com/hiveot/hub/pkg/authz"
	"github.com/hiveot/hub/pkg/authz/capnpserver"
	"github.com/hiveot/hub/pkg/authz/service"
)

const aclStoreFile = "authz.acl"

// main entry point to start the authorization service
func main() {
	logging.SetLogging("info", "")

	f := svcconfig.LoadServiceConfig(authz.ServiceName, false, nil)
	aclStoreFolder := filepath.Join(f.Stores, authz.ServiceName)
	aclStorePath := filepath.Join(aclStoreFolder, aclStoreFile)
	os.Mkdir(aclStoreFolder, 0700)

	// parse commandline and create server listening socket
	srvListener := listener.CreateServiceListener(f.Run, authz.ServiceName)
	ctx := context.Background()

	svc := service.NewAuthzService(ctx, aclStorePath)
	err := svc.Start(ctx)
	if err == nil {
		defer svc.Stop()
	}
	if err == nil {
		logrus.Infof("AuthzCapnpServer starting on: %s", srvListener.Addr())
		err = capnpserver.StartAuthzCapnpServer(ctx, srvListener, svc)
	}
	if err != nil {
		msg := fmt.Sprintf("ERROR: Service '%s' failed to start: %s\n", authz.ServiceName, err)
		logrus.Fatal(msg)
	}
	logrus.Warningf("Authz service ended gracefully")

	os.Exit(0)
}
