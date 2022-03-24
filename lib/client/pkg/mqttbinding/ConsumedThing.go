// Package mqttbinding that implements the MQTT protocol binding for the ConsumedThing API
// Consumed Things are remote representations of Things used by consumers.
package mqttbinding

import (
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/wostzone/hub/lib/client/pkg/mqttclient"
	"github.com/wostzone/hub/lib/client/pkg/thing"
	"strings"
	"sync"
)

// MqttConsumedThing is the implementation of an ConsumedThing interface using the MQTT protocol binding.
// Thing consumers can use this API to subscribe to events and actions.
//
// The protocol binding schema is work in progress, as is the WoT MQTT protocol binding. For now this takes
// a best guess approach based on "https://w3c.github.io/wot-binding-templates"
//
// This closely follows the WoT scripting API for ConsumedThing as described at
// https://www.w3.org/TR/wot-scripting-api/#the-consumedthing-interface
type MqttConsumedThing struct {
	// Client for MQTT message bus access
	mqttClient *mqttclient.MqttClient
	// Internal slot with Thing Description document this Thing exposes
	td *thing.ThingTD
	// internal slot for Subscriptions by event name
	eventSubscriptions map[string]Subscription
	// mutex for async updating of subscriptions
	subscriptionMutex sync.Mutex
	// valueStore holds property values as received by events
	valueStore map[string]InteractionOutput
}

// GetThingDescription returns the TD document of this exposed Thing
// This returns the cached version of the TD
func (cThing *MqttConsumedThing) GetThingDescription() *thing.ThingTD {
	return cThing.td
}

// Handle incoming events.
// If the event name is a property name then update the property value store using the property data schema,
// otherwise use the data schema in the events section of the TD.
// Last invoke the subscriber to the event name, if any, or the default subscriber
//  address is the MQTT topic that the event is published on as: things/{id}/event/{eventName}
//  whereas message is the body of the event.
func (cThing *MqttConsumedThing) handleEvent(address string, message []byte) {
	var evData InteractionOutput

	logrus.Infof("MqttConsumedThing.handleEvent: received event on topic %s", address)

	// the event topic is "things/id/event/name"
	parts := strings.Split(address, "/")
	if len(parts) < 4 {
		logrus.Warningf("MqttConsumedThing.handleEvent: EventName is missing in topic %s", address)
		return
	}
	eventName := parts[3]

	// handle property events
	propAffordance := cThing.td.GetProperty(eventName)
	if propAffordance != nil {
		evData = NewInteractionOutputFromJson(message, &propAffordance.DataSchema)

		logrus.Infof("MqttConsumedThing.handleEvent: Event with topic %s is a property event", address)
		// TODO validate the data
	} else {
		// not a property event
		eventAffordance := cThing.td.GetEvent(eventName)
		if eventAffordance != nil {
			evData = NewInteractionOutputFromJson(message, &eventAffordance.Data)
			logrus.Infof("MqttConsumedThing.handleEvent: Event with topic %s is not a property event", address)
		} else {
			// unknown schema
			evData = NewInteractionOutputFromJson(message, nil)
		}
	}
	// property or event, it is stored in the valueStore
	cThing.valueStore[eventName] = evData

	// notify subscribers if any
	subscription, found := cThing.eventSubscriptions[eventName]
	if !found || subscription.Handler == nil {
		subscription, found = cThing.eventSubscriptions[""] // default subscriber
	}
	if subscription.Handler != nil {
		subscription.Handler(eventName, evData)
	}
}

// InvokeAction makes a request for invoking an Action and returns once the
// request is submitted.
//
// This will be posted on topic: "things/{thingID}/action/{actionName}" with data as payload
//
// Takes as arguments actionName, optionally action data as defined in the TD.
// Returns nil if the action request was submitted successfully or an error if failed
func (cThing *MqttConsumedThing) InvokeAction(actionName string, data interface{}) error {
	topic := strings.ReplaceAll(TopicAction, "{id}", cThing.td.ID) + "/" + actionName
	return cThing.mqttClient.PublishObject(topic, data)
}

// ObserveProperty makes a request for Property value change notifications.
// Takes as arguments propertyName, listener.
func (cThing *MqttConsumedThing) ObserveProperty(
	name string, handler func(name string, data InteractionOutput)) error {

	var err error = nil
	sub := Subscription{
		SubType: SubscriptionTypeProperty,
		Name:    name,
		Handler: handler,
	}
	cThing.eventSubscriptions[name] = sub
	return err
}

// ReadProperty reads a Property value from the local cache
// Returns the last known property value as a string or an error if
// the name is not a known property.
func (cThing *MqttConsumedThing) ReadProperty(name string) (InteractionOutput, error) {
	//return res, errors.New("'"+name + "' is not a known property" )
	value, found := cThing.valueStore[name]
	if !found {
		return value, errors.New("Property " + name + " does not exist on thing " + cThing.td.ID)
	}
	return value, nil
}

// ReadMultipleProperties reads multiple Property values with one request.
// propertyNames is an array with names of properties to return
// Returns a PropertyMap object that maps keys from propertyNames to InteractionOutput of that property.
func (cThing *MqttConsumedThing) ReadMultipleProperties(names []string) map[string]InteractionOutput {
	res := make(map[string]InteractionOutput, 0)
	for _, name := range names {
		output, _ := cThing.ReadProperty(name)
		res[name] = output
	}
	return res
}

// ReadAllProperties reads all properties of the Thing with one request.
// Returns a PropertyMap object that maps keys from all Property names to InteractionOutput
// of the properties.
func (cThing *MqttConsumedThing) ReadAllProperties() map[string]InteractionOutput {
	res := make(map[string]InteractionOutput, 0)

	for name := range cThing.td.Properties {
		output, _ := cThing.ReadProperty(name)
		res[name] = output
	}
	return res
}

// Stop delivering notifications for event subscriptions
func (cThing *MqttConsumedThing) Stop() {
	topic := strings.ReplaceAll(TopicThingEvent, "{thingID}", cThing.td.ID)
	cThing.mqttClient.Unsubscribe(topic)
}

// SubscribeEvent makes a request for subscribing to Event or property change notifications.
// Takes as arguments eventName, listener and optionally onerror and option
// When eventName is a propertyName the event is a property value update.
// The engine already subscribes to events for updating properties, use this subscription to be notified of
// a change to a particular property or a specific event.
// Returns nil if subscription is successful or error if failed
func (cThing *MqttConsumedThing) SubscribeEvent(
	eventName string, handler func(eventName string, data InteractionOutput)) error {
	sub := Subscription{
		SubType: SubscriptionTypeEvent, // what is the purpose of capturing this?
		Name:    eventName,
		Handler: handler,
	}
	cThing.eventSubscriptions[eventName] = sub
	return nil
}

// WriteProperty submit a request to change a property value.
// Takes as arguments propertyName and value, and sends a property update to the exposedThing that in turn
// updates the actual device.
// This does not update the property immediately. It is up to the exposedThing to perform necessary validation
// and notify subscribers with an event after the change has been applied.
//
// There is no error feedback in case the request cannot be handled. The requester will only receive a
// property change event when the request has completed successfully. Failure to complete the request can be caused
// by an invalid value or if the IoT device is not in a state to accept changes.
//
// TBD: if there is a need to be notified of failure then a future update can add a write-property failed event.
//
// This will be published on topic "things/{thingID}/property/{name}"
//
// It returns an error if the property update could not be sent and nil if it is successfully
//  published. Final confirmation is obtained if an event is received with the updated property value.
func (cThing *MqttConsumedThing) WriteProperty(propName string, value interface{}) error {
	topic := strings.ReplaceAll(TopicThingProperty, "{id}", cThing.td.ID) + "/" + propName
	err := cThing.mqttClient.PublishObject(topic, value)
	if err != nil {
		logrus.Errorf("MqttConsumedThing:WriteProperty: Failed publishing update request on topic %s: %s", topic, err)
	}
	return err
}

// WriteMultipleProperties writes multiple property values.
// Takes as arguments properties - as a map keys being Property names and values as Property values.
//
// This will be posted as individual update requests:
//
// It returns an error if the action could not be sent and nil if the action is successfully
//  published. Final success is achieved if the property value will be updated through an event.
func (cThing *MqttConsumedThing) WriteMultipleProperties(properties map[string]interface{}) error {
	var err error
	for propName, value := range properties {
		err = cThing.WriteProperty(propName, value)
	}
	return err
}

// Consume constructs a consumed thing instance from a TD.
// A consumed Thing is a remote instance of a thing for the purpose of interaction by remote consumers.
// This subscribes to events from the remote thing.
//
// tdDoc is a Thing Description document of the Thing to expose.
// mqttClient client for binding to the MQTT protocol
func Consume(tdDoc *thing.ThingTD, mqttClient *mqttclient.MqttClient) *MqttConsumedThing {
	cThing := MqttConsumedThing{
		mqttClient:         mqttClient,
		td:                 tdDoc,
		eventSubscriptions: make(map[string]Subscription),
		subscriptionMutex:  sync.Mutex{},
		valueStore:         make(map[string]InteractionOutput),
	}
	// in order to keep props up to date, subscribe to all events
	topic := strings.ReplaceAll(TopicThingEvent, "{id}", cThing.td.ID) + "/#"
	cThing.mqttClient.Subscribe(topic, cThing.handleEvent)

	return &cThing
}
