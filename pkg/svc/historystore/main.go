// Package main with the history store
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v3"

	"github.com/hiveot/hub.grpc/go/svc"
	"github.com/hiveot/hub/internal/folders"
	"github.com/hiveot/hub/internal/listener"
	"github.com/hiveot/hub/pkg/svc/historystore/config"
	"github.com/hiveot/hub/pkg/svc/historystore/mongohs"
)

// DefaultConfigFile is the default configuration file with database settings
const DefaultConfigFile = "historystore.yaml"

// Start the history store service using gRPC
func main() {
	svcConfig := config.NewHistoryStoreConfig()
	configFile := path.Join(folders.GetFolders("").Config, DefaultConfigFile)

	// Add commandline option '-c configFile which holds service connection info
	flag.StringVar(&configFile, "c", configFile, "Service configuration with database connection info")
	lis := listener.CreateServiceListener(config.ServiceName)

	configData, err := ioutil.ReadFile(configFile)
	if err == nil {
		err = yaml.Unmarshal(configData, &svcConfig)
	}
	if err != nil {
		logrus.Fatalf("Error reading service configuration file '%s': %v", configFile, err)
	}

	// For now only mongodb is supported
	// This service needs the storage location and name
	service := mongohs.NewMongoHistoryStoreServer(svcConfig)
	s := grpc.NewServer()
	svc.RegisterHistoryStoreServer(s, service)

	// exit the service when signal is received and close the listener
	listener.ExitOnSignal(lis, func() {
		logrus.Infof("Shutting down '%s'", config.ServiceName)
	})

	// Start listening
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Service '%s; exited: %v", config.ServiceName, err)
	}
}
