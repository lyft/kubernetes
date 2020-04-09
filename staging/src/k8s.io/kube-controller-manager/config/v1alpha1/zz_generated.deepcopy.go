// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AttachDetachControllerConfiguration) DeepCopyInto(out *AttachDetachControllerConfiguration) {
	*out = *in
	out.ReconcilerSyncLoopPeriod = in.ReconcilerSyncLoopPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AttachDetachControllerConfiguration.
func (in *AttachDetachControllerConfiguration) DeepCopy() *AttachDetachControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(AttachDetachControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CSRSigningControllerConfiguration) DeepCopyInto(out *CSRSigningControllerConfiguration) {
	*out = *in
	out.ClusterSigningDuration = in.ClusterSigningDuration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CSRSigningControllerConfiguration.
func (in *CSRSigningControllerConfiguration) DeepCopy() *CSRSigningControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(CSRSigningControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CloudProviderConfiguration) DeepCopyInto(out *CloudProviderConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CloudProviderConfiguration.
func (in *CloudProviderConfiguration) DeepCopy() *CloudProviderConfiguration {
	if in == nil {
		return nil
	}
	out := new(CloudProviderConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CronJobControllerConfiguration) DeepCopyInto(out *CronJobControllerConfiguration) {
	*out = *in
	out.CronJobControllerSyncPeriod = in.CronJobControllerSyncPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CronJobControllerConfiguration.
func (in *CronJobControllerConfiguration) DeepCopy() *CronJobControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(CronJobControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DaemonSetControllerConfiguration) DeepCopyInto(out *DaemonSetControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DaemonSetControllerConfiguration.
func (in *DaemonSetControllerConfiguration) DeepCopy() *DaemonSetControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(DaemonSetControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeploymentControllerConfiguration) DeepCopyInto(out *DeploymentControllerConfiguration) {
	*out = *in
	out.DeploymentControllerSyncPeriod = in.DeploymentControllerSyncPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeploymentControllerConfiguration.
func (in *DeploymentControllerConfiguration) DeepCopy() *DeploymentControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeploymentControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DeprecatedControllerConfiguration) DeepCopyInto(out *DeprecatedControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DeprecatedControllerConfiguration.
func (in *DeprecatedControllerConfiguration) DeepCopy() *DeprecatedControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(DeprecatedControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EndpointControllerConfiguration) DeepCopyInto(out *EndpointControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EndpointControllerConfiguration.
func (in *EndpointControllerConfiguration) DeepCopy() *EndpointControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(EndpointControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GarbageCollectorControllerConfiguration) DeepCopyInto(out *GarbageCollectorControllerConfiguration) {
	*out = *in
	if in.EnableGarbageCollector != nil {
		in, out := &in.EnableGarbageCollector, &out.EnableGarbageCollector
		*out = new(bool)
		**out = **in
	}
	if in.GCIgnoredResources != nil {
		in, out := &in.GCIgnoredResources, &out.GCIgnoredResources
		*out = make([]GroupResource, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GarbageCollectorControllerConfiguration.
func (in *GarbageCollectorControllerConfiguration) DeepCopy() *GarbageCollectorControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(GarbageCollectorControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GenericControllerManagerConfiguration) DeepCopyInto(out *GenericControllerManagerConfiguration) {
	*out = *in
	out.MinResyncPeriod = in.MinResyncPeriod
	out.ClientConnection = in.ClientConnection
	out.ControllerStartInterval = in.ControllerStartInterval
	in.LeaderElection.DeepCopyInto(&out.LeaderElection)
	if in.Controllers != nil {
		in, out := &in.Controllers, &out.Controllers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	out.Debugging = in.Debugging
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GenericControllerManagerConfiguration.
func (in *GenericControllerManagerConfiguration) DeepCopy() *GenericControllerManagerConfiguration {
	if in == nil {
		return nil
	}
	out := new(GenericControllerManagerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GroupResource) DeepCopyInto(out *GroupResource) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GroupResource.
func (in *GroupResource) DeepCopy() *GroupResource {
	if in == nil {
		return nil
	}
	out := new(GroupResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *HPAControllerConfiguration) DeepCopyInto(out *HPAControllerConfiguration) {
	*out = *in
	out.HorizontalPodAutoscalerSyncPeriod = in.HorizontalPodAutoscalerSyncPeriod
	out.HorizontalPodAutoscalerUpscaleForbiddenWindow = in.HorizontalPodAutoscalerUpscaleForbiddenWindow
	out.HorizontalPodAutoscalerDownscaleStabilizationWindow = in.HorizontalPodAutoscalerDownscaleStabilizationWindow
	out.HorizontalPodAutoscalerDownscaleForbiddenWindow = in.HorizontalPodAutoscalerDownscaleForbiddenWindow
	if in.HorizontalPodAutoscalerUseRESTClients != nil {
		in, out := &in.HorizontalPodAutoscalerUseRESTClients, &out.HorizontalPodAutoscalerUseRESTClients
		*out = new(bool)
		**out = **in
	}
	out.HorizontalPodAutoscalerCPUInitializationPeriod = in.HorizontalPodAutoscalerCPUInitializationPeriod
	out.HorizontalPodAutoscalerInitialReadinessDelay = in.HorizontalPodAutoscalerInitialReadinessDelay
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new HPAControllerConfiguration.
func (in *HPAControllerConfiguration) DeepCopy() *HPAControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(HPAControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *JobControllerConfiguration) DeepCopyInto(out *JobControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new JobControllerConfiguration.
func (in *JobControllerConfiguration) DeepCopy() *JobControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(JobControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeCloudSharedConfiguration) DeepCopyInto(out *KubeCloudSharedConfiguration) {
	*out = *in
	out.CloudProvider = in.CloudProvider
	out.RouteReconciliationPeriod = in.RouteReconciliationPeriod
	out.NodeMonitorPeriod = in.NodeMonitorPeriod
	if in.ConfigureCloudRoutes != nil {
		in, out := &in.ConfigureCloudRoutes, &out.ConfigureCloudRoutes
		*out = new(bool)
		**out = **in
	}
	out.NodeSyncPeriod = in.NodeSyncPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeCloudSharedConfiguration.
func (in *KubeCloudSharedConfiguration) DeepCopy() *KubeCloudSharedConfiguration {
	if in == nil {
		return nil
	}
	out := new(KubeCloudSharedConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KubeControllerManagerConfiguration) DeepCopyInto(out *KubeControllerManagerConfiguration) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.Generic.DeepCopyInto(&out.Generic)
	in.KubeCloudShared.DeepCopyInto(&out.KubeCloudShared)
	out.AttachDetachController = in.AttachDetachController
	out.CSRSigningController = in.CSRSigningController
	out.CronJobController = in.CronJobController
	out.DaemonSetController = in.DaemonSetController
	out.DeploymentController = in.DeploymentController
	out.DeprecatedController = in.DeprecatedController
	out.EndpointController = in.EndpointController
	in.GarbageCollectorController.DeepCopyInto(&out.GarbageCollectorController)
	in.HPAController.DeepCopyInto(&out.HPAController)
	out.JobController = in.JobController
	out.NamespaceController = in.NamespaceController
	out.NodeIPAMController = in.NodeIPAMController
	in.NodeLifecycleController.DeepCopyInto(&out.NodeLifecycleController)
	in.PersistentVolumeBinderController.DeepCopyInto(&out.PersistentVolumeBinderController)
	out.PodGCController = in.PodGCController
	out.ReplicaSetController = in.ReplicaSetController
	out.ReplicationController = in.ReplicationController
	out.ResourceQuotaController = in.ResourceQuotaController
	out.SAController = in.SAController
	out.ServiceController = in.ServiceController
	out.TTLAfterFinishedController = in.TTLAfterFinishedController
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KubeControllerManagerConfiguration.
func (in *KubeControllerManagerConfiguration) DeepCopy() *KubeControllerManagerConfiguration {
	if in == nil {
		return nil
	}
	out := new(KubeControllerManagerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *KubeControllerManagerConfiguration) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NamespaceControllerConfiguration) DeepCopyInto(out *NamespaceControllerConfiguration) {
	*out = *in
	out.NamespaceSyncPeriod = in.NamespaceSyncPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NamespaceControllerConfiguration.
func (in *NamespaceControllerConfiguration) DeepCopy() *NamespaceControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(NamespaceControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeIPAMControllerConfiguration) DeepCopyInto(out *NodeIPAMControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeIPAMControllerConfiguration.
func (in *NodeIPAMControllerConfiguration) DeepCopy() *NodeIPAMControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(NodeIPAMControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeLifecycleControllerConfiguration) DeepCopyInto(out *NodeLifecycleControllerConfiguration) {
	*out = *in
	if in.EnableTaintManager != nil {
		in, out := &in.EnableTaintManager, &out.EnableTaintManager
		*out = new(bool)
		**out = **in
	}
	out.NodeStartupGracePeriod = in.NodeStartupGracePeriod
	out.NodeMonitorGracePeriod = in.NodeMonitorGracePeriod
	out.PodEvictionTimeout = in.PodEvictionTimeout
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeLifecycleControllerConfiguration.
func (in *NodeLifecycleControllerConfiguration) DeepCopy() *NodeLifecycleControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(NodeLifecycleControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeBinderControllerConfiguration) DeepCopyInto(out *PersistentVolumeBinderControllerConfiguration) {
	*out = *in
	out.PVClaimBinderSyncPeriod = in.PVClaimBinderSyncPeriod
	in.VolumeConfiguration.DeepCopyInto(&out.VolumeConfiguration)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeBinderControllerConfiguration.
func (in *PersistentVolumeBinderControllerConfiguration) DeepCopy() *PersistentVolumeBinderControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeBinderControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PersistentVolumeRecyclerConfiguration) DeepCopyInto(out *PersistentVolumeRecyclerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PersistentVolumeRecyclerConfiguration.
func (in *PersistentVolumeRecyclerConfiguration) DeepCopy() *PersistentVolumeRecyclerConfiguration {
	if in == nil {
		return nil
	}
	out := new(PersistentVolumeRecyclerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PodGCControllerConfiguration) DeepCopyInto(out *PodGCControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PodGCControllerConfiguration.
func (in *PodGCControllerConfiguration) DeepCopy() *PodGCControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(PodGCControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaSetControllerConfiguration) DeepCopyInto(out *ReplicaSetControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaSetControllerConfiguration.
func (in *ReplicaSetControllerConfiguration) DeepCopy() *ReplicaSetControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ReplicaSetControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicationControllerConfiguration) DeepCopyInto(out *ReplicationControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicationControllerConfiguration.
func (in *ReplicationControllerConfiguration) DeepCopy() *ReplicationControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ReplicationControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ResourceQuotaControllerConfiguration) DeepCopyInto(out *ResourceQuotaControllerConfiguration) {
	*out = *in
	out.ResourceQuotaSyncPeriod = in.ResourceQuotaSyncPeriod
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ResourceQuotaControllerConfiguration.
func (in *ResourceQuotaControllerConfiguration) DeepCopy() *ResourceQuotaControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ResourceQuotaControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SAControllerConfiguration) DeepCopyInto(out *SAControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SAControllerConfiguration.
func (in *SAControllerConfiguration) DeepCopy() *SAControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(SAControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServiceControllerConfiguration) DeepCopyInto(out *ServiceControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServiceControllerConfiguration.
func (in *ServiceControllerConfiguration) DeepCopy() *ServiceControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(ServiceControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *TTLAfterFinishedControllerConfiguration) DeepCopyInto(out *TTLAfterFinishedControllerConfiguration) {
	*out = *in
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new TTLAfterFinishedControllerConfiguration.
func (in *TTLAfterFinishedControllerConfiguration) DeepCopy() *TTLAfterFinishedControllerConfiguration {
	if in == nil {
		return nil
	}
	out := new(TTLAfterFinishedControllerConfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeConfiguration) DeepCopyInto(out *VolumeConfiguration) {
	*out = *in
	if in.EnableHostPathProvisioning != nil {
		in, out := &in.EnableHostPathProvisioning, &out.EnableHostPathProvisioning
		*out = new(bool)
		**out = **in
	}
	if in.EnableDynamicProvisioning != nil {
		in, out := &in.EnableDynamicProvisioning, &out.EnableDynamicProvisioning
		*out = new(bool)
		**out = **in
	}
	out.PersistentVolumeRecyclerConfiguration = in.PersistentVolumeRecyclerConfiguration
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeConfiguration.
func (in *VolumeConfiguration) DeepCopy() *VolumeConfiguration {
	if in == nil {
		return nil
	}
	out := new(VolumeConfiguration)
	in.DeepCopyInto(out)
	return out
}
