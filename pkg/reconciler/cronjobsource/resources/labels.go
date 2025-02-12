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

package resources

const (
	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "cronjob-source-controller"
)

// OldLabels are the pre-0.8 labels.
// TODO Delete after 0.8 is cut.
func OldLabels(name string) map[string]string {
	return map[string]string{
		"knative-eventing-source":      controllerAgentName,
		"knative-eventing-source-name": name,
	}
}

// Labels are the labels attached to all resources based on a CronJobSource.
func Labels(name string) map[string]string {
	return map[string]string{
		"sources.eventing.knative.dev/cronJobSource": name,
	}
}
