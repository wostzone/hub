// Package mqttbinding that implements the ExposedThing API using the MQTT protocol binding
// Exposed Things are used by IoT device implementers to provide access to the device.
package mqttbinding

import (
	"errors"
	"github.com/wostzone/hub/lib/client/pkg/mqttclient"
	"github.com/wostzone/hub/lib/client/pkg/thing"
	"strings"
	"sync"
)

// MqttExposedThing is the implementation of an ExposedThing interface using the MQTT protocol binding.
// Thing implementers can use this API to subscribe to actions and publish TDs and events.
//
// This closely follows the WoT scripting API for ExposedThing as described at
// https://www.w3.org/TR/wot-scripting-api/#the-exposedthing-interface
//
// Differences with the WoT scripting API:
//  1. The WoT scripting API uses ECMAScript with promises for asynchronous results.
//     This implementation uses channels to return async results, which is golang idiomatic.
//  2. The WoT scripting API method names are in 'lowerCase' format.
//     In golang lowerCase makes things private. This implementation uses 'UpperCase' name format.
//  3. Most methods are synchronous instead of asynchronous as the MQTT client is synchronous.
//     The result of actions indicates that it was submitted successfully. Actions do not have
//     a return value (in WoST) as they are not remote procedure calls. If the effect of an
//     action is needed then consumers should subscribe to property changes that are submitted
//     as the action is executed. The results of actions by others then be handled in the same way.
//  4. Actions are only handled by devices that are not asleep as the message bus does not
//     yet support queued actions. This is a limitation of the message bus. Future implementations
//     of the message bus can add queuing to support intermittent connected devices.
//     Use of the 'retain' flag is not recommended for actions on devices that also have a manual input.
//  5. If an action is not allowed then no error is returned. In most cases the MQTT bus won't accept
//     the request in the first place.
//
// Example of properties in a TD with forms for mqtt protocol binding.
// The forms will likely be provided through a @context link to a standardized WoST model, once the semantics are figured out.
// {
//   "properties": {
//        "onoff": {
//            "@type": "iot:SwitchOnOff",
//            "title": "Switch on or off status"
//            "description": "More elaborate description of the onoff property"
//            "observable": true,    // form must provide a observeproperty binding
//            "type": "boolean",
//            "unit": "binary",
//            "readOnly": true,
//            "writeOnly": true,     // property is writable. form must provide a writeproperty binding
//            // protocol binding for the property
//        }
//        // These forms apply to all writable properties
//        "forms": [{
//          	"href": "mqtt://${broker}/things/${thingID}/property/${propertyName}",
//          	"op": ["writeproperty"],
//				"mqv:controlPacketValue": "PUBLISH",
//              "contentType": "application/json"
//          }, {
//              // TBD. MQTT topic. How to parameterize in a generic schema?
//          	"href": "mqtt://${broker}/things/${thingID}/property/${propertyName}",
//              "op": ["observeproperty"],
//				"mqv:controlPacketValue": "SUBSCRIBE",
//              "contentType": "application/json"
//         }],
//       }
//    }
// }

type MqttExposedThing struct {
	// ExposedThing extends a ConsumedThing
	MqttConsumedThing

	// handlers for executing a requested action. name is the action name. "" for default handler
	actionHandlers map[string]func(actionName string, value InteractionOutput)

	// handler for writing to properties when a write request is received
	// to set the default write handler use name ""
	writeHandlers map[string]func(name string, value InteractionOutput) error
}

// Destroy the exposed thing. This stops serving external requests
func (eThing *MqttExposedThing) Destroy() {
	topic := strings.ReplaceAll(TopicThingEvent, "{id}", eThing.td.ID) + "/#"
	eThing.mqttClient.Unsubscribe(topic)

	eThing.eventSubscriptions = nil
	eThing.actionHandlers = nil
	eThing.writeHandlers = nil
}

// EmitPropertyChange publishes a property change event, which in turn will notify all
// observers (subscribers) of the property.
//
// The topic will be things/{thingID}/event/{propName} and payload will be the new property value
// propName is the name of the property in the TD.
// newValue is the new value of the property, used as the event payload
// Returns an error if the property value cannot be published
func (eThing *MqttExposedThing) EmitPropertyChange(propName string, newValue InteractionOutput) error {
	// update the cached value in ConsumedThing base class
	eThing.valueStore[propName] = newValue

	// submit the property change as an event with the property name
	topic := strings.ReplaceAll(TopicThingEvent, "{id}", eThing.td.ID) + "/" + propName
	err := eThing.mqttClient.PublishObject(topic, newValue)
	return err
}

// EmitEvent publishes an event to subscribers.
// The topic will be things/{thingID}/event/{name} and payload will be the event data.
// If the event cannot be published, for example when disconnected, an error is returned.
// For important events this can be used to retry them when connection is restored.
//
// name is the name of the event as described in the TD, or one of the general purpose events.
// data is the event value as defined in the TD events schema and used as the payload
// Returns an error if the event is not found or cannot be published
func (eThing *MqttExposedThing) EmitEvent(name string, data interface{}) error {
	_, found := eThing.td.Events[name]
	if !found {
		return errors.New("EmitEvent: event name '" + name + "' is not defined in the TD document")
	}

	topic := strings.ReplaceAll(TopicThingEvent, "{id}", eThing.td.ID) + "/" + name
	err := eThing.mqttClient.PublishObject(topic, data)
	return err
}

// Expose starts serving external requests for the Thing, so that WoT Interactions using Properties, Actions and Events
// will be possible. This also publishes the TD document of this Thing.
func (eThing *MqttExposedThing) Expose() error {
	topic := strings.ReplaceAll(TopicThingEvent, "{id}", eThing.td.ID) + "/#"
	eThing.mqttClient.Subscribe(topic, func(address string, message []byte) {})

	// Also publish this Thing's TD document
	topic = strings.ReplaceAll(TopicThingTD, "{id}", eThing.td.ID)
	err := eThing.mqttClient.PublishObject(topic, eThing.td)
	return err
}

//func (eThing *MqttExposedThing) SetPropertyReadHandler(func(name string) string) error {
//	return errors.New("not implemented")
//}

// SetActionHandler sets the handler for handling an action for the IoT device.
//  Only a single handler is active. If a handler is set when a previous handler was already set then the
//  latest handler will be used.
//
// The device code should implement this handler to updated configuration of the device.
//
// actionName is the action name this handler is for. If a single handler can take care of most actions
//  then use "" as the name to indicate it is the default handler.
//
// The handler should return nil if the write is accepted or an error if not accepted. The property value
// in the TD will be updated after the property has changed through the change notification handler.
func (eThing *MqttExposedThing) SetActionHandler(
	actionName string, actionHandler func(actionName string, value InteractionOutput)) {

	eThing.actionHandlers[actionName] = actionHandler
}

// SetPropertyObserveHandler sets the handler for subscribing to properties
// Not implemented as subscriptions are handled by the MQTT message bus
func (eThing *MqttExposedThing) SetPropertyObserveHandler(handler func(name string) InteractionOutput) error {
	_ = handler
	return errors.New("not implemented")
}

// SetPropertyUnobserveHandler sets the handler for unsubscribing to properties
// Not implemented as subscriptions are handled by the MQTT message bus
func (eThing *MqttExposedThing) SetPropertyUnobserveHandler(handler func(name string) InteractionOutput) error {
	_ = handler
	return errors.New("not implemented")
}

// SetPropertyReadHandler sets the handler for reading a property of the IoT device
// Not implemented as property values are updated with events and not requested.
// The latest property value can be found with the TD properties.
func (eThing *MqttExposedThing) SetPropertyReadHandler(handler func(name string) string) error {
	_ = handler
	return errors.New("not implemented")
}

// SetPropertyWriteHandler sets the handler for writing a property of the IoT device. This typically applies
// to device configuration. If the property is read-only the handler must return an error.
//  Only a single handler is active. If a handler is set when a previous handler was already set then the
//  latest handler will be used.
//
// The device code should implement this handler to updated configuration of the device.
//
// propName is the property name this handler is for. If a single handler can take care of writing most properties
//  then use "" as the name to indicate it is the default handler.
//
// The handler should return nil if the write is accepted or an error if not accepted. The property value
// in the TD will be updated after the property has changed through the change notification handler.
func (eThing *MqttExposedThing) SetPropertyWriteHandler(
	propName string, writeHandler func(propName string, value InteractionOutput) error) {

	eThing.writeHandlers[propName] = writeHandler
}

// Produce constructs an exposed thing from a TD.
// An exposed Thing is a local instance of a thing for the purpose of interaction
// with remote consumers.
//
// tdDoc is a Thing Description document of the Thing to expose.
// mqttClient client for binding to the MQTT protocol
func Produce(tdDoc *thing.ThingTD, mqttClient *mqttclient.MqttClient) *MqttExposedThing {
	eThing := &MqttExposedThing{
		MqttConsumedThing: MqttConsumedThing{
			mqttClient:         mqttClient,
			td:                 tdDoc,
			eventSubscriptions: make(map[string]Subscription),
			subscriptionMutex:  sync.Mutex{},
		},
		actionHandlers: make(map[string]func(actionName string, value InteractionOutput)),
		writeHandlers:  make(map[string]func(name string, value InteractionOutput) error),
	}
	return eThing
}
