//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Condition) DeepCopyInto(out *Condition) {
	*out = *in
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Condition.
func (in *Condition) DeepCopy() *Condition {
	if in == nil {
		return nil
	}
	out := new(Condition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DashboardSpec) DeepCopyInto(out *DashboardSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DashboardSpec.
func (in *DashboardSpec) DeepCopy() *DashboardSpec {
	if in == nil {
		return nil
	}
	out := new(DashboardSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorConfig) DeepCopyInto(out *EstimatorConfig) {
	*out = *in
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(ModelSelectorSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorConfig.
func (in *EstimatorConfig) DeepCopy() *EstimatorConfig {
	if in == nil {
		return nil
	}
	out := new(EstimatorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorGroup) DeepCopyInto(out *EstimatorGroup) {
	*out = *in
	if in.Total != nil {
		in, out := &in.Total, &out.Total
		*out = new(EstimatorConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.Components != nil {
		in, out := &in.Components, &out.Components
		*out = new(EstimatorConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorGroup.
func (in *EstimatorGroup) DeepCopy() *EstimatorGroup {
	if in == nil {
		return nil
	}
	out := new(EstimatorGroup)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EstimatorStatus) DeepCopyInto(out *EstimatorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EstimatorStatus.
func (in *EstimatorStatus) DeepCopy() *EstimatorStatus {
	if in == nil {
		return nil
	}
	out := new(EstimatorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExporterDeploymentSpec) DeepCopyInto(out *ExporterDeploymentSpec) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExporterDeploymentSpec.
func (in *ExporterDeploymentSpec) DeepCopy() *ExporterDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(ExporterDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExporterSpec) DeepCopyInto(out *ExporterSpec) {
	*out = *in
	in.Deployment.DeepCopyInto(&out.Deployment)
	if in.Redfish != nil {
		in, out := &in.Redfish, &out.Redfish
		*out = new(RedfishSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExporterSpec.
func (in *ExporterSpec) DeepCopy() *ExporterSpec {
	if in == nil {
		return nil
	}
	out := new(ExporterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExporterStatus) DeepCopyInto(out *ExporterStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExporterStatus.
func (in *ExporterStatus) DeepCopy() *ExporterStatus {
	if in == nil {
		return nil
	}
	out := new(ExporterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalEstimatorSpec) DeepCopyInto(out *InternalEstimatorSpec) {
	*out = *in
	in.Node.DeepCopyInto(&out.Node)
	in.Container.DeepCopyInto(&out.Container)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalEstimatorSpec.
func (in *InternalEstimatorSpec) DeepCopy() *InternalEstimatorSpec {
	if in == nil {
		return nil
	}
	out := new(InternalEstimatorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalExporterDeploymentSpec) DeepCopyInto(out *InternalExporterDeploymentSpec) {
	*out = *in
	in.ExporterDeploymentSpec.DeepCopyInto(&out.ExporterDeploymentSpec)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalExporterDeploymentSpec.
func (in *InternalExporterDeploymentSpec) DeepCopy() *InternalExporterDeploymentSpec {
	if in == nil {
		return nil
	}
	out := new(InternalExporterDeploymentSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalExporterSpec) DeepCopyInto(out *InternalExporterSpec) {
	*out = *in
	in.Deployment.DeepCopyInto(&out.Deployment)
	if in.Redfish != nil {
		in, out := &in.Redfish, &out.Redfish
		*out = new(RedfishSpec)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalExporterSpec.
func (in *InternalExporterSpec) DeepCopy() *InternalExporterSpec {
	if in == nil {
		return nil
	}
	out := new(InternalExporterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InternalModelServerSpec) DeepCopyInto(out *InternalModelServerSpec) {
	*out = *in
	in.Storage.DeepCopyInto(&out.Storage)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InternalModelServerSpec.
func (in *InternalModelServerSpec) DeepCopy() *InternalModelServerSpec {
	if in == nil {
		return nil
	}
	out := new(InternalModelServerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Kepler) DeepCopyInto(out *Kepler) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Kepler.
func (in *Kepler) DeepCopy() *Kepler {
	if in == nil {
		return nil
	}
	out := new(Kepler)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Kepler) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerInternal) DeepCopyInto(out *KeplerInternal) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerInternal.
func (in *KeplerInternal) DeepCopy() *KeplerInternal {
	if in == nil {
		return nil
	}
	out := new(KeplerInternal)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeplerInternal) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerInternalList) DeepCopyInto(out *KeplerInternalList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]KeplerInternal, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerInternalList.
func (in *KeplerInternalList) DeepCopy() *KeplerInternalList {
	if in == nil {
		return nil
	}
	out := new(KeplerInternalList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeplerInternalList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerInternalSpec) DeepCopyInto(out *KeplerInternalSpec) {
	*out = *in
	in.Exporter.DeepCopyInto(&out.Exporter)
	if in.Estimator != nil {
		in, out := &in.Estimator, &out.Estimator
		*out = new(InternalEstimatorSpec)
		(*in).DeepCopyInto(*out)
	}
	if in.ModelServer != nil {
		in, out := &in.ModelServer, &out.ModelServer
		*out = new(InternalModelServerSpec)
		(*in).DeepCopyInto(*out)
	}
	out.OpenShift = in.OpenShift
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerInternalSpec.
func (in *KeplerInternalSpec) DeepCopy() *KeplerInternalSpec {
	if in == nil {
		return nil
	}
	out := new(KeplerInternalSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerInternalStatus) DeepCopyInto(out *KeplerInternalStatus) {
	*out = *in
	in.Exporter.DeepCopyInto(&out.Exporter)
	out.Estimator = in.Estimator
	out.ModelServer = in.ModelServer
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerInternalStatus.
func (in *KeplerInternalStatus) DeepCopy() *KeplerInternalStatus {
	if in == nil {
		return nil
	}
	out := new(KeplerInternalStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerList) DeepCopyInto(out *KeplerList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Kepler, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerList.
func (in *KeplerList) DeepCopy() *KeplerList {
	if in == nil {
		return nil
	}
	out := new(KeplerList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KeplerList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerSpec) DeepCopyInto(out *KeplerSpec) {
	*out = *in
	in.Exporter.DeepCopyInto(&out.Exporter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerSpec.
func (in *KeplerSpec) DeepCopy() *KeplerSpec {
	if in == nil {
		return nil
	}
	out := new(KeplerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KeplerStatus) DeepCopyInto(out *KeplerStatus) {
	*out = *in
	in.Exporter.DeepCopyInto(&out.Exporter)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KeplerStatus.
func (in *KeplerStatus) DeepCopy() *KeplerStatus {
	if in == nil {
		return nil
	}
	out := new(KeplerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelSelectorSpec) DeepCopyInto(out *ModelSelectorSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelSelectorSpec.
func (in *ModelSelectorSpec) DeepCopy() *ModelSelectorSpec {
	if in == nil {
		return nil
	}
	out := new(ModelSelectorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServerStatus) DeepCopyInto(out *ModelServerStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServerStatus.
func (in *ModelServerStatus) DeepCopy() *ModelServerStatus {
	if in == nil {
		return nil
	}
	out := new(ModelServerStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelServerStorageSpec) DeepCopyInto(out *ModelServerStorageSpec) {
	*out = *in
	if in.PersistentVolumeClaim != nil {
		in, out := &in.PersistentVolumeClaim, &out.PersistentVolumeClaim
		*out = new(v1.PersistentVolumeClaimSpec)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelServerStorageSpec.
func (in *ModelServerStorageSpec) DeepCopy() *ModelServerStorageSpec {
	if in == nil {
		return nil
	}
	out := new(ModelServerStorageSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OpenShiftSpec) DeepCopyInto(out *OpenShiftSpec) {
	*out = *in
	out.Dashboard = in.Dashboard
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OpenShiftSpec.
func (in *OpenShiftSpec) DeepCopy() *OpenShiftSpec {
	if in == nil {
		return nil
	}
	out := new(OpenShiftSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RedfishSpec) DeepCopyInto(out *RedfishSpec) {
	*out = *in
	out.ProbeInterval = in.ProbeInterval
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RedfishSpec.
func (in *RedfishSpec) DeepCopy() *RedfishSpec {
	if in == nil {
		return nil
	}
	out := new(RedfishSpec)
	in.DeepCopyInto(out)
	return out
}
