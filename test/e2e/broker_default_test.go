// +build e2e

/*
Copyright 2019 The Knative Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"knative.dev/eventing/pkg/apis/eventing/v1alpha1"
	pkgResources "knative.dev/eventing/pkg/reconciler/namespace/resources"
	"knative.dev/eventing/test/base/resources"
	"knative.dev/eventing/test/common"

	"k8s.io/apimachinery/pkg/util/uuid"
	"knative.dev/pkg/test/logging"
)

const (
	waitForFilterPodRunning = 30 * time.Second
	selectorKey             = "end2end-test-broker-trigger"

	defaultBrokerName = pkgResources.DefaultBrokerName
	any               = v1alpha1.TriggerAnyFilter
	eventType1        = "type1"
	eventType2        = "type2"
	eventSource1      = "source1"
	eventSource2      = "source2"
)

// eventTypeAndSource specifies the type and source of an Event.
type eventTypeAndSource struct {
	Type   string
	Source string
}

// Helper struct to tie the type and sources of the events we expect to receive
// in subscribers with the selectors we use when creating their pods.
type eventReceiver struct {
	typeAndSource eventTypeAndSource
	selector      map[string]string
}

// This test annotates the testing namespace so that a default broker is created.
// It then binds many triggers with different filtering patterns to that default broker,
// and sends different events to the broker's address. Finally, it verifies that only
// the appropriate events are routed to the subscribers.
func TestDefaultBrokerWithManyTriggers(t *testing.T) {
	client := setup(t, true)
	defer tearDown(client)

	// Label namespace so that it creates the default broker.
	if err := client.LabelNamespace(map[string]string{"knative-eventing-injection": "enabled"}); err != nil {
		t.Fatalf("Error annotating namespace: %v", err)
	}

	// Wait for default broker ready.
	if err := client.WaitForResourceReady(defaultBrokerName, common.BrokerTypeMeta); err != nil {
		t.Fatalf("Error waiting for default broker to become ready: %v", err)
	}

	// These are the event types and sources that triggers will listen to, as well as the selectors
	// to set  in the subscriber and services pods.
	eventsToReceive := []eventReceiver{
		{eventTypeAndSource{Type: any, Source: any}, newSelector()},
		{eventTypeAndSource{Type: eventType1, Source: any}, newSelector()},
		{eventTypeAndSource{Type: any, Source: eventSource1}, newSelector()},
		{eventTypeAndSource{Type: eventType1, Source: eventSource1}, newSelector()},
	}

	// Create subscribers.
	for _, event := range eventsToReceive {
		subscriberName := name("dumper", event.typeAndSource.Type, event.typeAndSource.Source)
		pod := resources.EventLoggerPod(subscriberName)
		client.CreatePodOrFail(pod, common.WithService(subscriberName))
	}

	// Create triggers.
	for _, event := range eventsToReceive {
		triggerName := name("trigger", event.typeAndSource.Type, event.typeAndSource.Source)
		subscriberName := name("dumper", event.typeAndSource.Type, event.typeAndSource.Source)
		client.CreateTriggerOrFail(triggerName,
			resources.WithSubscriberRefForTrigger(subscriberName),
			resources.WithTriggerFilter(event.typeAndSource.Source, event.typeAndSource.Type),
		)
	}

	// Wait for all test resources to become ready before sending the events.
	if err := client.WaitForAllTestResourcesReady(); err != nil {
		t.Fatalf("Failed to get all test resources ready: %v", err)
	}

	// These are the event types and sources that will be send.
	eventsToSend := []eventTypeAndSource{
		{eventType1, eventSource1},
		{eventType1, eventSource2},
		{eventType2, eventSource1},
		{eventType2, eventSource2},
	}
	// Map to save the expected events per dumper so that we can verify the delivery.
	expectedEvents := make(map[string][]string)
	// Map to save the unexpected events per dumper so that we can verify that they weren't delivered.
	unexpectedEvents := make(map[string][]string)
	for _, eventToSend := range eventsToSend {
		// Create cloud event.
		// Using event type and source as part of the body for easier debugging.
		body := fmt.Sprintf("Body-%s-%s", eventToSend.Type, eventToSend.Source)
		cloudEvent := &resources.CloudEvent{
			Source: eventToSend.Source,
			Type:   eventToSend.Type,
			Data:   fmt.Sprintf(`{"msg":%q}`, body),
		}
		// Create sender pod.
		senderPodName := name("sender", eventToSend.Type, eventToSend.Source)
		if err := client.SendFakeEventToAddressable(senderPodName, defaultBrokerName, common.BrokerTypeMeta, cloudEvent); err != nil {
			t.Fatalf("Error send cloud event to broker: %v", err)
		}

		// Check on every dumper whether we should expect this event or not, and add its body
		// to the expectedEvents/unexpectedEvents maps.
		for _, eventToReceive := range eventsToReceive {
			subscriberName := name("dumper", eventToReceive.typeAndSource.Type, eventToReceive.typeAndSource.Source)
			if shouldExpectEvent(&eventToSend, &eventToReceive, t.Logf) {
				expectedEvents[subscriberName] = append(expectedEvents[subscriberName], body)
			} else {
				unexpectedEvents[subscriberName] = append(unexpectedEvents[subscriberName], body)
			}
		}
	}

	for _, event := range eventsToReceive {
		subscriberName := name("dumper", event.typeAndSource.Type, event.typeAndSource.Source)
		if err := client.CheckLog(subscriberName, common.CheckerContainsAll(expectedEvents[subscriberName])); err != nil {
			t.Fatalf("Event(s) not found in logs of subscriber pod %q: %v", subscriberName, err)
		}
		// At this point all the events should have been received in the pod.
		// We check whether we find unexpected events. If so, then we fail.
		found, err := client.FindAnyLogContents(subscriberName, unexpectedEvents[subscriberName])
		if err != nil {
			t.Fatalf("Failed querying to find log contents in pod %q: %v", subscriberName, err)
		}
		if found {
			t.Fatalf("Unexpected event(s) found in logs of subscriber pod %q", subscriberName)
		}
	}
}

// Helper function to create names for different objects (e.g., triggers, services, etc.).
func name(obj, eventType, eventSource string) string {
	// Pod names need to be lowercase. We might have an eventType as Any, that is why we lowercase them.
	if eventType == "" {
		eventType = "testany"
	}
	if eventSource == "" {
		eventSource = "testany"
	}
	return strings.ToLower(fmt.Sprintf("%s-%s-%s", obj, eventType, eventSource))
}

// Returns a new selector with a random uuid.
func newSelector() map[string]string {
	return map[string]string{selectorKey: string(uuid.NewUUID())}
}

// Checks whether we should expect to receive 'eventToSend' in 'eventReceiver' based on its type and source pattern.
func shouldExpectEvent(eventToSend *eventTypeAndSource, receiver *eventReceiver, logf logging.FormatLogger) bool {
	if receiver.typeAndSource.Type != any && receiver.typeAndSource.Type != eventToSend.Type {
		return false
	}
	if receiver.typeAndSource.Source != any && receiver.typeAndSource.Source != eventToSend.Source {
		return false
	}
	return true
}
