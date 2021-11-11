package main

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/sirupsen/logrus"
	"github.com/wostzone/hub/auth/pkg/aclstore"
	"github.com/wostzone/hub/auth/pkg/authenticate"
	"github.com/wostzone/hub/auth/pkg/authorize"
	"github.com/wostzone/hub/auth/pkg/unpwstore"
	"github.com/wostzone/hub/lib/client/pkg/certs"
	"github.com/wostzone/hub/lib/client/pkg/config"
	"github.com/wostzone/hub/lib/client/pkg/signing"
	"github.com/wostzone/hub/lib/serve/pkg/certsetup"
	"github.com/wostzone/hub/lib/serve/pkg/hubnet"
)

//// Commandline commands
//const (
//	CmdCertBundle = "certbundle"
//	CmdClientcert = "clientcert"
//	CmdSetPasswd  = "setpasswd"
//	CmdSetRole    = "setrole"
//)
const Version = `0.2-alpha`

func main() {
	binFolder := path.Dir(os.Args[0])
	homeFolder := path.Dir(binFolder)
	ParseArgs(homeFolder, os.Args[1:])
}

// ParseArgs to handle commandline arguments
func ParseArgs(homeFolder string, args []string) {
	// var err error
	configFolder := path.Join(homeFolder, "config")
	certsFolder := path.Join(homeFolder, "certs")
	// configFolder := path.Join(homeFolder, "config")
	// ouRole := certsetup.OUClient
	// genKeys := false
	ifName, mac, ip := hubnet.GetOutboundInterface("")
	_ = ifName
	_ = mac
	sanName := ip.String()
	var optConf struct {
		// commands
		Certbundle  bool
		Clientcert  bool
		Devicecert  bool
		Setpassword bool
		Setrole     bool
		// arguments
		Loginid  string
		Deviceid string
		Groupid  string
		Role     string
		// options
		Aclfile string
		Config  string
		Certs   string
		Output  string
		Pubkey  string
		Iter    int
		Verbose bool
	}
	usage := `
Usage:
  auth certbundle [-v --certs=CertFolder]
  auth clientcert [-v --certs=CertFolder --pubkey=pubkeyfile] <loginID> 
  auth devicecert [-v --certs=CertFolder --pubkey=pubkeyfile] <deviceID> 
  auth setpassword [-v -c configFolder] [-i iterations] <loginID> 
  auth setrole [-v -c configFolder --aclfile=aclfile] <loginID> <groupID> <role>
  auth --help | --version

Commands:
  certbundle   Generate or refresh the Hub certificate bundle
  clientcert   Generate a signed client certificate, with pub/private keys if not given
  devicecert   Generate a signed device certificate, with pub/private keys if not given
  setpassword  Set user password
  setrole      Set user role in group

Arguments:
  loginID      used as the certificate CN, login name and certificate filename (loginID-cert.pem)
  groupID      group for access control 
  role         one of viewer, editor, manager, thing, or none to delete

Options:
  --aclfile=AclFile              use a different acl file instead of the default config/` + aclstore.DefaultAclFile + `
  -e --certs=CertFolder      location of Hub certificates [default: ` + certsFolder + `]
  -c --config=ConfigFolder   location of Hub config folder [default: ` + configFolder + `]
  -p --pubkey=PubKeyfile     use this public key file to generate certificate, instead of a new key pair
	-i --iter=iterations       Number of iterations for generating password [default: 10]
  -h --help                  show this help
	-v --verbose               show info logging
  --version                  show app version
`
	opts, err := docopt.ParseArgs(usage, args, Version)
	if err != nil {
		fmt.Printf("Parse Error: %s\n", err)
		os.Exit(1)
	}

	err = opts.Bind(&optConf)

	if optConf.Verbose {
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(logrus.WarnLevel)
	}

	if err != nil {
		fmt.Printf("Bind Error: %s\n", err)
		os.Exit(1)
	}
	_ = opts
	if optConf.Certbundle {
		fmt.Printf("Generating certificate Bundle. Certfolder=%s\n", optConf.Certs)
		err = HandleCreateCertbundle(optConf.Certs, sanName)
	} else if optConf.Clientcert {
		fmt.Printf("Generating Client certificate using CA from %s\n", optConf.Certs)
		err = HandleCreateClientCert(optConf.Certs, optConf.Loginid, optConf.Pubkey)
	} else if optConf.Devicecert {
		fmt.Printf("Generating Thing device certificate using CA from %s\n", optConf.Certs)
		err = HandleCreateDeviceCert(optConf.Certs, optConf.Deviceid, optConf.Pubkey)
	} else if optConf.Setpassword {
		fmt.Printf("Set user password\n")
		err = HandleSetPasswd(optConf.Config, optConf.Loginid, optConf.Iter)
	} else if optConf.Setrole {
		fmt.Printf("Set user role in group\n")
		err = HandleSetRole(optConf.Config, optConf.Loginid, optConf.Groupid, optConf.Role, optConf.Aclfile)
	} else {
		err = fmt.Errorf("invalid command")
	}
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(1)
	}
}

// CreateKeyPair generate a key pair in PEM format and save it to the cert folder
// as <clientID>-pub.pem and <clientID>-key.pem
// Returns the public key PEM content
func CreateKeyPair(clientID string, certFolder string) (privKey *ecdsa.PrivateKey, err error) {
	privKey = signing.CreateECDSAKeys()
	privKeyFile := path.Join(certFolder, clientID+"-priv.pem")
	pubKeyFile := path.Join(certFolder, clientID+"-pub.pem")
	err = certs.SaveKeysToPEM(privKey, privKeyFile)
	if err == nil {
		pubKeyPem, _ := certs.PublicKeyToPEM(&privKey.PublicKey)
		err = ioutil.WriteFile(pubKeyFile, []byte(pubKeyPem), 0644)
	}
	if err != nil {
		fmt.Printf("Failed saving keys: %s\n", err)
	}
	if err == nil {
		fmt.Printf("Generated public and private key pair as: %s and %s\n", pubKeyFile, privKeyFile)
	}
	return privKey, err
}

// HandleSetPasswd sets the login name and password for a consumer
func HandleSetPasswd(configFolder string, username string, iterations int) error {
	var pwHash string
	var err error
	var passwd string
	reader := bufio.NewReader(os.Stdin)
	unpwFilePath := path.Join(configFolder, unpwstore.DefaultPasswordFile)
	unpwStore := unpwstore.NewPasswordFileStore(unpwFilePath, "auth.main.HandleSetPasswd")
	err = unpwStore.Open()
	if err == nil {
		fmt.Printf("\nNew Password: ")
		passwd, err = reader.ReadString('\n')
		passwd = strings.Replace(passwd, "\n", "", -1)
		if err != nil {
			return err
		}
		if passwd == "" {
			return fmt.Errorf("missing password")
		}
		// pwHash, err = authen.CreatePasswordHash(passwd, authen.PWHASH_ARGON2id, uint(iterations))
		pwHash, err = authenticate.CreatePasswordHash(passwd, authenticate.PWHASH_ARGON2id, uint(iterations))
	}
	if err == nil {
		err = unpwStore.SetPasswordHash(username, pwHash)
	}
	if err == nil {
		unpwStore.Close()
		fmt.Printf("Password updated for user %s\n", username)
	}
	return err
}

// HandleCreateCertbundle generates the hub certificate bundle CA, Hub and Plugin keys
// and certificates.
//  If the CA certificate already exist it is NOT updated
//  If the Hub and Plugin certificates already exist, they are renewed
func HandleCreateCertbundle(certsFolder string, sanName string) error {
	err := certsetup.CreateCertificateBundle([]string{sanName}, certsFolder)
	if err != nil {
		return err
	}
	fmt.Printf("Server and Plugin certificates generated in %s\n", certsFolder)
	return nil
}

// HandleCreateClientCert creates a consumer client certificate and optionally private/public keypair
//  certFolder where to find the CA certificate and key used to sign the client certificate
//  clientID for the CN of the client certificate. Used to identify the consumer.
//  pubKeyFile with path to the client's public key of the certificate
func HandleCreateClientCert(certFolder string, clientID string, pubKeyFile string) error {
	var pubKey *ecdsa.PublicKey
	ou := certsetup.OUClient
	pemPath := path.Join(certFolder, config.DefaultCaCertFile)
	caCert, err := certs.LoadX509CertFromPEM(pemPath)
	if err != nil {
		return err
	}
	pemPath = path.Join(certFolder, config.DefaultCaKeyFile)
	caKey, err := certs.LoadKeysFromPEM(pemPath)
	if err != nil {
		return err
	}
	// If a public key file is given, use it, otherwise generate a pair
	if pubKeyFile != "" {
		fmt.Printf("Using public key file: %s\n", pubKeyFile)
		pubKey, err = certs.LoadPublicKeyFromPEM(pubKeyFile)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No public key file was provided. Creating a key pair ") // no newline
		privKey, err := CreateKeyPair(clientID, "")
		pubKey = &privKey.PublicKey
		if err != nil {
			return err
		}
	}
	durationDays := certsetup.DefaultCertDurationDays
	cert, err := certsetup.CreateHubClientCert(
		clientID, ou, pubKey, caCert, caKey, time.Now(), durationDays)
	if err != nil {
		return err
	}
	pemPath = path.Join(".", clientID+"-cert.pem")
	certs.SaveX509CertToPEM(cert, pemPath)

	fmt.Printf("Client certificate saved at %s\n", pemPath)
	return nil
}

// HandleCreateDeviceCert creates a device client certificate for a device and save it in the certFolder
// This is similar to creating a consumer certificate
func HandleCreateDeviceCert(certFolder string, deviceID string, pubKeyFile string) error {
	const deviceCertValidityDays = 30
	var pubKey *ecdsa.PublicKey
	pemPath := path.Join(certFolder, config.DefaultCaCertFile)
	caCert, err := certs.LoadX509CertFromPEM(pemPath)
	if err != nil {
		return err
	}
	pemPath = path.Join(certFolder, config.DefaultCaKeyFile)
	caKey, err := certs.LoadKeysFromPEM(pemPath)
	if err != nil {
		return err
	}
	// If a public key file is given, use it, otherwise generate a pair
	if pubKeyFile != "" {
		fmt.Printf("Using public key file: %s\n", pubKeyFile)
		pubKey, err = certs.LoadPublicKeyFromPEM(pubKeyFile)
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("No public key file was provided. Creating a new key pair ") // no newline
		privKey, err := CreateKeyPair(deviceID, "")
		pubKey = &privKey.PublicKey
		if err != nil {
			return err
		}
	}

	certPEM, err := certsetup.CreateHubClientCert(
		deviceID, certsetup.OUIoTDevice, pubKey,
		caCert, caKey,
		time.Now(), deviceCertValidityDays)
	if err != nil {
		return err
	}
	// save the new certificate
	pemPath = path.Join(".", deviceID+"-cert.pem")
	err = certs.SaveX509CertToPEM(certPEM, pemPath)

	fmt.Printf("Device certificate saved at %s\n", pemPath)
	return err
}

// HandleSetRole sets the role of a client in a group.
func HandleSetRole(configFolder string, clientID string, groupID string, role string, aclFile string) error {
	if role != authorize.GroupRoleEditor && role != authorize.GroupRoleViewer &&
		role != authorize.GroupRoleManager && role != authorize.GroupRoleThing && role != authorize.GroupRoleNone {
		err := fmt.Errorf("invalid role '%s'", role)
		return err
	}
	aclFilePath := path.Join(configFolder, aclstore.DefaultAclFile)
	if aclFile != "" {
		// option to specify an acl file wrt home
		aclFilePath = path.Join(path.Dir(configFolder), aclFile)
	}
	aclStore := aclstore.NewAclFileStore(aclFilePath, "author.main.HandleSetRole")
	err := aclStore.Open()
	if err == nil {
		err = aclStore.SetRole(clientID, groupID, role)
	}
	if err == nil {
		fmt.Printf("Client '%s' role set to '%s' for group '%s'\n", clientID, role, groupID)
	}
	aclStore.Close()
	return err
}
