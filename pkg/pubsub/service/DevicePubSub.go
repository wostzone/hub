package service

import (
	"context"
	"encoding/json"

	"github.com/sirupsen/logrus"

	"github.com/hiveot/hub.capnp/go/vocab"
	"github.com/hiveot/hub/lib/caphelp"

	"github.com/hiveot/hub/lib/thing"
	"github.com/hiveot/hub/pkg/pubsub"
	"github.com/hiveot/hub/pkg/pubsub/core"
)

// DevicePubSub provides pub/sub capability to IoT devices.
// The IoT device is a gateway for the Things it manages, hence it has a gateway ID that is also
// its ThingID.
type DevicePubSub struct {
	// the publisherID is the thingID of the IoT device or service that publishes
	publisherID string
	// core is the pubsub engine
	core *core.PubSubCore
	// subscriptionIDs from the core
	subscriptionIDs []string
}

// PubEvent publishes the given thing event. The payload is an event value as per TD.
func (dps *DevicePubSub) PubEvent(
	_ context.Context, thingID, name string, value []byte) (err error) {

	thingAddr := thing.MakeThingAddr(dps.publisherID, thingID)
	tv := thing.NewThingValue(thingAddr, name, caphelp.Clone(value))
	// note that marshal will copy the value so its buffer can be reused by capnp
	tvSerialized, _ := json.Marshal(tv)
	topic := MakeThingTopic(thingAddr, pubsub.MessageTypeEvent, name)
	go dps.core.Publish(topic, tvSerialized)
	return
}

// PubProperties publishes an event with the given properties.
// The props is a map of property name-value pairs.
func (dps *DevicePubSub) PubProperties(
	_ context.Context, thingID string, props map[string][]byte) (err error) {

	thingAddr := thing.MakeThingAddr(dps.publisherID, thingID)
	propsValue, _ := json.Marshal(props)
	tv := thing.NewThingValue(thingAddr, vocab.WoTProperties, propsValue)
	// note that marshal will copy the props map so its buffer can be reused by capnp
	tvSerialized, _ := json.Marshal(tv)
	topic := MakeThingTopic(thingAddr, pubsub.MessageTypeEvent, vocab.WoTProperties)
	dps.core.Publish(topic, tvSerialized)
	return
}

// PubTD publishes the given thing TD as an event. The payload is a TD document.
// The event MUST be from the same device.
func (dps *DevicePubSub) PubTD(_ context.Context,
	thingID string, deviceType string, td []byte) (err error) {

	thingAddr := thing.MakeThingAddr(dps.publisherID, thingID)
	tv := thing.NewThingValue(thingAddr, pubsub.MessageTypeTD, td)
	// note that marshal will copy the TD so its buffer can be reused by capnp
	tvSerialized, _ := json.Marshal(tv)
	topic := MakeThingTopic(thingAddr, pubsub.MessageTypeTD, deviceType)
	dps.core.Publish(topic, tvSerialized)
	return
}

// SubAction subscribes to messages for the given thingID and action name
//
//	thingID and actionName are optional. Use "" to receive actions for all things or names.
func (dps *DevicePubSub) SubAction(
	_ context.Context, thingID string, actionName string,
	handler func(actionValue *thing.ThingValue)) (err error) {
	thingAddr := thing.MakeThingAddr(dps.publisherID, thingID)

	topic := MakeThingTopic(thingAddr, pubsub.MessageTypeAction, actionName)
	subscriptionID, err := dps.core.Subscribe(topic,
		func(topic string, message []byte) {

			msgValue := &thing.ThingValue{}
			err = json.Unmarshal(message, msgValue)
			if err != nil {
				logrus.Error(err)
			}
			handler(msgValue)
		})
	if err == nil {
		dps.subscriptionIDs = append(dps.subscriptionIDs, subscriptionID)
	}
	return err
}

// Release the capability and end subscriptions
func (dps *DevicePubSub) Release() {
	err := dps.core.Unsubscribe(dps.subscriptionIDs)

	if err != nil {
		logrus.Errorf("IoT device %s unsubscribe failed: %s", dps.publisherID, err)
	}
	dps.subscriptionIDs = nil
}

// NewDevicePubSub provides the capability for a device to publish actions and subscribe to events
//
//	publisherID is the thingID of the IoT device doing the publishing
//	core is the core pubsub that is used for publishing and subscribing
func NewDevicePubSub(gatewayID string, core *core.PubSubCore) *DevicePubSub {
	deviceCap := &DevicePubSub{
		publisherID:     gatewayID,
		core:            core,
		subscriptionIDs: make([]string, 0),
	}
	return deviceCap
}
