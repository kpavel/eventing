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

// Code generated by informer-gen. DO NOT EDIT.

package externalversions

import (
	"fmt"

	schema "k8s.io/apimachinery/pkg/runtime/schema"
	cache "k8s.io/client-go/tools/cache"
	v1alpha1 "knative.dev/eventing/pkg/apis/eventing/v1alpha1"
	messagingv1alpha1 "knative.dev/eventing/pkg/apis/messaging/v1alpha1"
	sourcesv1alpha1 "knative.dev/eventing/pkg/apis/sources/v1alpha1"
)

// GenericInformer is type of SharedIndexInformer which will locate and delegate to other
// sharedInformers based on type
type GenericInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() cache.GenericLister
}

type genericInformer struct {
	informer cache.SharedIndexInformer
	resource schema.GroupResource
}

// Informer returns the SharedIndexInformer.
func (f *genericInformer) Informer() cache.SharedIndexInformer {
	return f.informer
}

// Lister returns the GenericLister.
func (f *genericInformer) Lister() cache.GenericLister {
	return cache.NewGenericLister(f.Informer().GetIndexer(), f.resource)
}

// ForResource gives generic access to a shared informer of the matching type
// TODO extend this to unknown resources with a client pool
func (f *sharedInformerFactory) ForResource(resource schema.GroupVersionResource) (GenericInformer, error) {
	switch resource {
	// Group=eventing.knative.dev, Version=v1alpha1
	case v1alpha1.SchemeGroupVersion.WithResource("brokers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().Brokers().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("channels"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().Channels().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("clusterchannelprovisioners"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().ClusterChannelProvisioners().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("eventtypes"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().EventTypes().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("subscriptions"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().Subscriptions().Informer()}, nil
	case v1alpha1.SchemeGroupVersion.WithResource("triggers"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Eventing().V1alpha1().Triggers().Informer()}, nil

		// Group=messaging.knative.dev, Version=v1alpha1
	case messagingv1alpha1.SchemeGroupVersion.WithResource("channels"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Messaging().V1alpha1().Channels().Informer()}, nil
	case messagingv1alpha1.SchemeGroupVersion.WithResource("choices"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Messaging().V1alpha1().Choices().Informer()}, nil
	case messagingv1alpha1.SchemeGroupVersion.WithResource("inmemorychannels"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Messaging().V1alpha1().InMemoryChannels().Informer()}, nil
	case messagingv1alpha1.SchemeGroupVersion.WithResource("sequences"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Messaging().V1alpha1().Sequences().Informer()}, nil

		// Group=sources.eventing.knative.dev, Version=v1alpha1
	case sourcesv1alpha1.SchemeGroupVersion.WithResource("apiserversources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Sources().V1alpha1().ApiServerSources().Informer()}, nil
	case sourcesv1alpha1.SchemeGroupVersion.WithResource("containersources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Sources().V1alpha1().ContainerSources().Informer()}, nil
	case sourcesv1alpha1.SchemeGroupVersion.WithResource("cronjobsources"):
		return &genericInformer{resource: resource.GroupResource(), informer: f.Sources().V1alpha1().CronJobSources().Informer()}, nil

	}

	return nil, fmt.Errorf("no informer found for %v", resource)
}
