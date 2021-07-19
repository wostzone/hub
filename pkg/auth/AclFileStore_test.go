package auth_test

import (
	"os"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/wostzone/hub/pkg/auth"
	"github.com/wostzone/wostlib-go/pkg/hubconfig"
)

const aclFileName = "acl-test.yaml"

var configFolder string
var aclFile string
var aclStore *auth.AclFileStore

// Use the config folder to store the acl files
func TestMain(m *testing.M) {
	hubconfig.SetLogging("info", "")
	cwd, _ := os.Getwd()
	homeFolder := path.Join(cwd, "../../test")
	configFolder = path.Join(homeFolder, "config")

	// Make sure an ACL file exist
	aclFile = path.Join(configFolder, aclFileName)
	fp, _ := os.Create(aclFile)
	fp.Close()
	aclStore = auth.NewAclFileStore(aclFile)

	res := m.Run()
	os.Exit(res)
}

func TestOpenCloseAclStore(t *testing.T) {
	// as := auth.NewAclStoreFile(aclFile)
	err := aclStore.Open()
	assert.NoError(t, err)
	time.Sleep(time.Second * 1)
	assert.NoError(t, err)
	aclStore.Close()
}

func TestSetRole(t *testing.T) {
	// as := auth.NewAclStoreFile(aclFile)
	user1 := "user1"
	role1 := auth.GroupRoleManager
	group1 := "group1"
	err := aclStore.Open()
	assert.NoError(t, err)

	err = aclStore.SetRole(user1, group1, role1)
	assert.NoError(t, err)

	// time to reload
	time.Sleep(time.Second)

	groups := aclStore.GetGroups(user1)
	assert.GreaterOrEqual(t, len(groups), 1)

	role := aclStore.GetRole(user1, groups)
	assert.Equal(t, auth.GroupRoleManager, role)

	aclStore.Close()
}

func TestWriteAclToTempFail(t *testing.T) {

	err := aclStore.Open()
	assert.NoError(t, err)
	_, err = aclStore.WriteToTemp("/badfolder")
	assert.Error(t, err)
	aclStore.Close()
}

func TestCompareRoles(t *testing.T) {
	ge := auth.IsRoleGreaterEqual(auth.GroupRoleViewer, auth.GroupRoleNone)
	assert.True(t, ge)
	ge = auth.IsRoleGreaterEqual(auth.GroupRoleNone, auth.GroupRoleViewer)
	assert.False(t, ge)

	ge = auth.IsRoleGreaterEqual(auth.GroupRoleEditor, auth.GroupRoleViewer)
	assert.True(t, ge)
	ge = auth.IsRoleGreaterEqual(auth.GroupRoleViewer, auth.GroupRoleEditor)
	assert.False(t, ge)

	ge = auth.IsRoleGreaterEqual(auth.GroupRoleManager, auth.GroupRoleEditor)
	assert.True(t, ge)
	ge = auth.IsRoleGreaterEqual(auth.GroupRoleEditor, auth.GroupRoleManager)
	assert.False(t, ge)

}

func TestMissingAclFile(t *testing.T) {
	as := auth.NewAclFileStore("missingaclfile")
	err := as.Open()
	assert.Error(t, err)
	as.Close()

}

func TestBadAclFile(t *testing.T) {
	// loading the hub-bad.yaml should fail as it isn't a valid yaml file
	as := auth.NewAclFileStore(path.Join(configFolder, "mosquitto.conf.template"))
	err := as.Open()
	assert.Error(t, err)
	as.Close()
}

func TestFailWriteFile(t *testing.T) {
	as := auth.NewAclFileStore("/root/nopermissions")

	err := as.Open()
	assert.Error(t, err)

	// err = os.Chmod(aclFile, 0400)
	// assert.NoError(t, err)

	// err = aclStore.SetRole("user1", "group1", "somerole")
	// assert.Error(t, err)
	// os.Remove(aclFile)
	aclStore.Close()
}
