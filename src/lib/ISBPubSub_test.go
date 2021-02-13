package lib

import (
	"io/ioutil"
	"net/http"
	"path"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func startLittleServer(t *testing.T, hostPort string) *http.Server {
	var upgrader websocket.Upgrader = websocket.Upgrader{}

	router := mux.NewRouter()
	router.HandleFunc("/channel/Chan1/{stage}", func(resp http.ResponseWriter, req *http.Request) {
		pubOrSub := mux.Vars(req)["stage"]
		logrus.Infof("TestNewPubSub: calling socket upgrade to websocket: %s", pubOrSub)
		upgrader.Upgrade(resp, req, nil)
	})

	httpServer := &http.Server{
		Addr:    hostPort,
		Handler: router,
	}
	go func() {
		// cs.updateMutex.Unlock()
		err := httpServer.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			assert.NoError(t, err, "ListenAndServe failed: %s", err)
		}
	}()
	return httpServer
}

func TestNewPubSub(t *testing.T) {
	const channel1 = "Chan1"
	const hostPort = "localhost:9678"
	const client1ID = "cid1"
	const certFolder = "../../test"
	var err error
	httpServer := startLittleServer(t, hostPort)
	time.Sleep(100 * time.Millisecond)

	pubConn, err := NewPublisher(hostPort, client1ID, channel1)
	assert.NoError(t, err)

	subConn, err := NewSubscriber(hostPort, client1ID, channel1, func(channel string, msg []byte) {
	})
	assert.NoError(t, err)

	err = SendMessage(pubConn, []byte("Hello world"))
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)
	httpServer.Close()
	subConn.Close()
	pubConn.Close()
	time.Sleep(100 * time.Millisecond)

}

func TestNewPubSubErrors(t *testing.T) {
	const channel1 = "Chan1"
	const hostPort = "localhost:9678"
	const client1ID = "cid1"
	const certFolder = "../../test"

	serverCertPath := path.Join(certFolder, ServerCertFile)
	clientCertPath := path.Join(certFolder, ClientCertFile)
	clientKeyPath := path.Join(certFolder, ClientKeyFile)

	serverCertPEM, _ := ioutil.ReadFile(serverCertPath)
	clientCertPEM, _ := ioutil.ReadFile(clientCertPath)
	clientKeyPEM, _ := ioutil.ReadFile(clientKeyPath)

	logrus.Infof("Testing authentication on channel %s", channel1)
	// cs, err := servicebus.StartServiceBus(hostPort)
	// require.NoError(t, err)
	time.Sleep(time.Second)

	_, err := NewPublisher(hostPort, client1ID, channel1)
	assert.Error(t, err)

	_, err = NewTLSPublisher(hostPort, client1ID, channel1, clientCertPEM, clientKeyPEM, serverCertPEM)
	assert.Error(t, err)

	_, err = NewSubscriber(hostPort, client1ID, channel1, func(channel string, msg []byte) {})
	assert.Error(t, err, "Error creating subscriber")

	_, err = NewTLSSubscriber(hostPort, client1ID, channel1, clientCertPEM, clientKeyPEM, serverCertPEM, func(channel string, msg []byte) {})
	assert.Error(t, err, "Error creating subscriber")

	// c := &websocket.Conn{}
	// err = SendMessage(c, []byte("no c error"))
	// require.Error(t, err, "Error creating subscriber")

	// cs.Stop()
}
