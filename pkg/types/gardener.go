package types

// import (
// 	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
// )

// // ShootList is a list of Shoot objects.
// type ShootList struct {
// 	metav1.TypeMeta `json:",inline"`
// 	// Standard list object metadata.
// 	// +optional
// 	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
// 	// Items is the list of Shoots.
// 	Items []Shoot `json:"items" protobuf:"bytes,2,rep,name=items"`
// }

// // Shoot represents a Shoot cluster created and managed by Gardener.
// type Shoot struct {
// 	metav1.TypeMeta `json:",inline"`
// 	// Standard object metadata.
// 	// +optional
// 	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
// 	// Specification of the Shoot cluster.
// 	// If the object's deletion timestamp is set, this field is immutable.
// 	// +optional
// 	Spec ShootSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
// 	// Most recently observed status of the Shoot cluster.
// 	// +optional
// 	Status ShootStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
// }

// // ShootSpec is the specification of a Shoot.
// type ShootSpec struct {
// 	// Addons contains information about enabled/disabled addons and their configuration.
// 	// +optional
// 	Addons *Addons `json:"addons,omitempty" protobuf:"bytes,1,opt,name=addons"`
// 	// CloudProfileName is a name of a CloudProfile object. This field is immutable.
// 	CloudProfileName string `json:"cloudProfileName" protobuf:"bytes,2,opt,name=cloudProfileName"`
// 	// DNS contains information about the DNS settings of the Shoot.
// 	// +optional
// 	DNS *DNS `json:"dns,omitempty" protobuf:"bytes,3,opt,name=dns"`
// 	// Extensions contain type and provider information for Shoot extensions.
// 	// +optional
// 	Extensions []Extension `json:"extensions,omitempty" protobuf:"bytes,4,rep,name=extensions"`
// 	// Hibernation contains information whether the Shoot is suspended or not.
// 	// +optional
// 	Hibernation *Hibernation `json:"hibernation,omitempty" protobuf:"bytes,5,opt,name=hibernation"`
// 	// Kubernetes contains the version and configuration settings of the control plane components.
// 	Kubernetes Kubernetes `json:"kubernetes" protobuf:"bytes,6,opt,name=kubernetes"`
// 	// Networking contains information about cluster networking such as CNI Plugin type, CIDRs, ...etc.
// 	Networking Networking `json:"networking" protobuf:"bytes,7,opt,name=networking"`
// 	// Maintenance contains information about the time window for maintenance operations and which
// 	// operations should be performed.
// 	// +optional
// 	Maintenance *Maintenance `json:"maintenance,omitempty" protobuf:"bytes,8,opt,name=maintenance"`
// 	// Monitoring contains information about custom monitoring configurations for the shoot.
// 	// +optional
// 	Monitoring *Monitoring `json:"monitoring,omitempty" protobuf:"bytes,9,opt,name=monitoring"`
// 	// Provider contains all provider-specific and provider-relevant information.
// 	Provider Provider `json:"provider" protobuf:"bytes,10,opt,name=provider"`
// 	// Purpose is the purpose class for this cluster.
// 	// +optional
// 	Purpose *ShootPurpose `json:"purpose,omitempty" protobuf:"bytes,11,opt,name=purpose,casttype=ShootPurpose"`
// 	// Region is a name of a region. This field is immutable.
// 	Region string `json:"region" protobuf:"bytes,12,opt,name=region"`
// 	// SecretBindingName is the name of the a SecretBinding that has a reference to the provider secret.
// 	// The credentials inside the provider secret will be used to create the shoot in the respective account.
// 	// This field is immutable.
// 	SecretBindingName string `json:"secretBindingName" protobuf:"bytes,13,opt,name=secretBindingName"`
// 	// SeedName is the name of the seed cluster that runs the control plane of the Shoot.
// 	// This field is immutable when the SeedChange feature gate is disabled.
// 	// +optional
// 	SeedName *string `json:"seedName,omitempty" protobuf:"bytes,14,opt,name=seedName"`
// 	// SeedSelector is an optional selector which must match a seed's labels for the shoot to be scheduled on that seed.
// 	// +optional
// 	SeedSelector *SeedSelector `json:"seedSelector,omitempty" protobuf:"bytes,15,opt,name=seedSelector"`
// 	// Resources holds a list of named resource references that can be referred to in extension configs by their names.
// 	// +optional
// 	Resources []NamedResourceReference `json:"resources,omitempty" protobuf:"bytes,16,rep,name=resources"`
// 	// Tolerations contains the tolerations for taints on seed clusters.
// 	// +patchMergeKey=key
// 	// +patchStrategy=merge
// 	// +optional
// 	Tolerations []Toleration `json:"tolerations,omitempty" patchStrategy:"merge" patchMergeKey:"key" protobuf:"bytes,17,rep,name=tolerations"`
// 	// ExposureClassName is the optional name of an exposure class to apply a control plane endpoint exposure strategy.
// 	// This field is immutable.
// 	// +optional
// 	ExposureClassName *string `json:"exposureClassName,omitempty" protobuf:"bytes,18,opt,name=exposureClassName"`
// 	// SystemComponents contains the settings of system components in the control or data plane of the Shoot cluster.
// 	// +optional
// 	SystemComponents *SystemComponents `json:"systemComponents" protobuf:"bytes,19,opt,name=systemComponents"`
// }
