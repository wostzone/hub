package clientconfigstore_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/wostzone/hub/authn/pkg/clientconfigstore"
	"github.com/wostzone/wost-go/pkg/logging"
	"os"
	"path"
	"testing"
)

var storeFolder = ""

// TestMain determines the store location
// Used for all test cases in this package
func TestMain(m *testing.M) {
	logging.SetLogging("info", "")
	cwd, _ := os.Getwd()
	homeFolder := path.Join(cwd, "../../test")
	storeFolder = path.Join(homeFolder, "configStore")

	res := m.Run()
	os.Exit(res)
}

func TestOpenClose(t *testing.T) {
	cfgStore := clientconfigstore.NewClientConfigStore(storeFolder)
	err := cfgStore.Open()
	assert.NoError(t, err)
	cfgStore.Close()
}

func TestWriteConfig(t *testing.T) {
	const user1ID = "user1"
	const app1ID = "app1"
	const configData = "Hello world"
	cfgStore := clientconfigstore.NewClientConfigStore(storeFolder)
	err := cfgStore.Open()
	assert.NoError(t, err)

	err = cfgStore.Put(user1ID, app1ID, configData)
	assert.NoError(t, err)

	rxData := cfgStore.Get(user1ID, app1ID)
	assert.Equal(t, configData, rxData)
	cfgStore.Close()
}

func TestReadMissingConfig(t *testing.T) {
	const user1ID = "user1"
	const app2ID = "app2"
	cfgStore := clientconfigstore.NewClientConfigStore(storeFolder)
	err := cfgStore.Open()
	assert.NoError(t, err)

	rxData := cfgStore.Get(user1ID, app2ID)
	assert.Empty(t, rxData)
	cfgStore.Close()
}
