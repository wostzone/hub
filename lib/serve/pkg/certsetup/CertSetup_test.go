package certsetup_test

import (
	"crypto/x509"
	"os"
	"os/exec"
	"path"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/wostzone/hub/lib/client/pkg/certs"
	"github.com/wostzone/hub/lib/client/pkg/config"
	"github.com/wostzone/hub/lib/serve/pkg/certsetup"
)

var homeFolder string
var certFolder string

// removeCerts easy cleanup for existing device certificate
func removeServerCerts() {
	_, _ = exec.Command("sh", "-c", "rm -f "+path.Join(certFolder, "*.pem")).Output()
}

// TestMain clears the certs folder for clean testing
func TestMain(m *testing.M) {
	cwd, _ := os.Getwd()
	homeFolder = path.Join(cwd, "../../test")
	certFolder = path.Join(homeFolder, "certs")
	config.SetLogging("info", "")
	removeServerCerts()

	res := m.Run()
	os.Exit(res)
}

// func TestLoadCreateCertKeyBadFile(t *testing.T) {
// 	removeServerCerts()
// 	_, err := certsetup.LoadOrCreateCertKey("/root/nopermission.pem")
// 	assert.Error(t, err)
// }
// func TestTLSCertificateGeneration(t *testing.T) {
// 	hostnames := []string{"127.0.0.1"}
// 	clientID := "3rdparty-client"

// 	// test creating ca a

func TestCreateCA(t *testing.T) {
	// test creating hub CA certificate
	caCert, caKeys := certsetup.CreateHubCA()
	require.NotNil(t, caCert)
	require.NotNil(t, caKeys)
}

func TestClientCertBadCA(t *testing.T) {
	clientID := "client1"
	ou := certsetup.OUClient
	caCert, caKey := certsetup.CreateHubCA()
	keys := certs.CreateECDSAKeys()

	clientCert, err := certsetup.CreateHubClientCert(clientID, ou,
		&keys.PublicKey, nil, caKey, time.Now(), certsetup.TempCertDurationDays)
	assert.Error(t, err)
	assert.Empty(t, clientCert)

	clientCert, err = certsetup.CreateHubClientCert(clientID, ou,
		&keys.PublicKey, caCert, nil, time.Now(), certsetup.TempCertDurationDays)
	assert.Error(t, err)
	assert.Empty(t, clientCert)
}

func TestCreateServerCert(t *testing.T) {
	// test creating hub certificate
	names := []string{"127.0.0.1", "localhost"}
	caCert, caKey := certsetup.CreateHubCA()
	cert, err := certsetup.CreateHubServerCert(names, caCert, caKey)
	require.NoError(t, err)
	require.NotNil(t, cert)
	require.NotNil(t, cert.PrivateKey)

	// todo, verify names in certificate
}

func TestServerCertBadCA(t *testing.T) {
	hostnames := []string{"127.0.0.1"}
	caCert, caKey := certsetup.CreateHubCA()
	//
	hubCert, err := certsetup.CreateHubServerCert(hostnames, caCert, nil)
	require.Error(t, err)
	require.Empty(t, hubCert)

	hubCert, err = certsetup.CreateHubServerCert(hostnames, nil, caKey)
	require.Error(t, err)
	require.Empty(t, hubCert)

	badCa := x509.Certificate{}
	hubCert, err = certsetup.CreateHubServerCert(hostnames, &badCa, caKey)
	require.Error(t, err)
	require.Empty(t, hubCert)
}
func TestCreateClientCert(t *testing.T) {
	clientID := "plugin1"
	ou := certsetup.OUPlugin
	// test creating hub certificate
	caCert, caKeys := certsetup.CreateHubCA()
	keys := certs.CreateECDSAKeys()

	hubCert, err := certsetup.CreateHubClientCert(clientID, ou,
		&keys.PublicKey, caCert, caKeys, time.Now(), 1)
	require.NoErrorf(t, err, "TestServiceCert: Failed creating server certificate")
	require.NotNil(t, hubCert)
}
func TestCreateDeviceCert(t *testing.T) {
	deviceID := "device1"
	ou := certsetup.OUIoTDevice
	// test creating hub certificate
	caCert, caKeys := certsetup.CreateHubCA()
	keys := certs.CreateECDSAKeys()

	hubCert, err := certsetup.CreateHubClientCert(deviceID, ou,
		&keys.PublicKey, caCert, caKeys, time.Now(), 1)
	require.NoErrorf(t, err, "TestServiceCert: Failed creating server certificate")
	require.NotNil(t, hubCert)
}

func TestCreateBundle(t *testing.T) {
	hostnames := []string{"127.0.0.1"}

	// test creating hub CA certificate
	err := certsetup.CreateCertificateBundle(hostnames, certFolder)
	require.NoError(t, err)
}

func TestCreateBundleBadFolder(t *testing.T) {
	hostnames := []string{"127.0.0.1"}

	// test creating hub CA certificate
	err := certsetup.CreateCertificateBundle(hostnames, "/not/a/valid/folder")
	require.Error(t, err)
}

func TestCreateBundleBadNames(t *testing.T) {
	// test creating hub CA certificate
	err := certsetup.CreateCertificateBundle(nil, certFolder)
	require.Error(t, err)
}
