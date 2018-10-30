/*
Copyright 2018 the Heptio Ark contributors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/heptio/ark/pkg/apis/ark/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// VolumeSnapshotLocationLister helps list VolumeSnapshotLocations.
type VolumeSnapshotLocationLister interface {
	// List lists all VolumeSnapshotLocations in the indexer.
	List(selector labels.Selector) (ret []*v1.VolumeSnapshotLocation, err error)
	// VolumeSnapshotLocations returns an object that can list and get VolumeSnapshotLocations.
	VolumeSnapshotLocations(namespace string) VolumeSnapshotLocationNamespaceLister
	VolumeSnapshotLocationListerExpansion
}

// volumeSnapshotLocationLister implements the VolumeSnapshotLocationLister interface.
type volumeSnapshotLocationLister struct {
	indexer cache.Indexer
}

// NewVolumeSnapshotLocationLister returns a new VolumeSnapshotLocationLister.
func NewVolumeSnapshotLocationLister(indexer cache.Indexer) VolumeSnapshotLocationLister {
	return &volumeSnapshotLocationLister{indexer: indexer}
}

// List lists all VolumeSnapshotLocations in the indexer.
func (s *volumeSnapshotLocationLister) List(selector labels.Selector) (ret []*v1.VolumeSnapshotLocation, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.VolumeSnapshotLocation))
	})
	return ret, err
}

// VolumeSnapshotLocations returns an object that can list and get VolumeSnapshotLocations.
func (s *volumeSnapshotLocationLister) VolumeSnapshotLocations(namespace string) VolumeSnapshotLocationNamespaceLister {
	return volumeSnapshotLocationNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// VolumeSnapshotLocationNamespaceLister helps list and get VolumeSnapshotLocations.
type VolumeSnapshotLocationNamespaceLister interface {
	// List lists all VolumeSnapshotLocations in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.VolumeSnapshotLocation, err error)
	// Get retrieves the VolumeSnapshotLocation from the indexer for a given namespace and name.
	Get(name string) (*v1.VolumeSnapshotLocation, error)
	VolumeSnapshotLocationNamespaceListerExpansion
}

// volumeSnapshotLocationNamespaceLister implements the VolumeSnapshotLocationNamespaceLister
// interface.
type volumeSnapshotLocationNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all VolumeSnapshotLocations in the indexer for a given namespace.
func (s volumeSnapshotLocationNamespaceLister) List(selector labels.Selector) (ret []*v1.VolumeSnapshotLocation, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.VolumeSnapshotLocation))
	})
	return ret, err
}

// Get retrieves the VolumeSnapshotLocation from the indexer for a given namespace and name.
func (s volumeSnapshotLocationNamespaceLister) Get(name string) (*v1.VolumeSnapshotLocation, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("volumesnapshotlocation"), name)
	}
	return obj.(*v1.VolumeSnapshotLocation), nil
}
