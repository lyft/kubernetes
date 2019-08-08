/*
Copyright 2017 The Kubernetes Authors.

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

package stats

import (
	"fmt"
	"path"
	"sort"
	"strings"

	cadvisorapiv2 "github.com/google/cadvisor/info/v2"
	"k8s.io/klog"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	statsapi "k8s.io/kubernetes/pkg/kubelet/apis/stats/v1alpha1"
	"k8s.io/kubernetes/pkg/kubelet/cadvisor"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/leaky"
	"k8s.io/kubernetes/pkg/kubelet/server/stats"
	"k8s.io/kubernetes/pkg/kubelet/status"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// cadvisorStatsProvider implements the containerStatsProvider interface by
// getting the container stats from cAdvisor. This is needed by docker and rkt
// integrations since they do not provide stats from CRI.
type cadvisorStatsProvider struct {
	// cadvisor is used to get the stats of the cgroup for the containers that
	// are managed by pods.
	cadvisor cadvisor.Interface
	// resourceAnalyzer is used to get the volume stats of the pods.
	resourceAnalyzer stats.ResourceAnalyzer
	// imageService is used to get the stats of the image filesystem.
	imageService kubecontainer.ImageService
	// statusProvider is used to get pod metadata
	statusProvider status.PodStatusProvider
}

// newCadvisorStatsProvider returns a containerStatsProvider that provides
// container stats from cAdvisor.
func newCadvisorStatsProvider(
	cadvisor cadvisor.Interface,
	resourceAnalyzer stats.ResourceAnalyzer,
	imageService kubecontainer.ImageService,
	statusProvider status.PodStatusProvider,
) containerStatsProvider {
	return &cadvisorStatsProvider{
		cadvisor:         cadvisor,
		resourceAnalyzer: resourceAnalyzer,
		imageService:     imageService,
		statusProvider:   statusProvider,
	}
}

// ListPodStats returns the stats of all the pod-managed containers.
func (p *cadvisorStatsProvider) ListPodStats() ([]statsapi.PodStats, error) {
	// Gets node root filesystem information and image filesystem stats, which
	// will be used to populate the available and capacity bytes/inodes in
	// container stats.

	klog.Warningf("in cadvisor_stats_provider.ListPodStats: %+v", p)

	rootFsInfo, err := p.cadvisor.RootFsInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get rootFs info: %v", err)
	}
	imageFsInfo, err := p.cadvisor.ImagesFsInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get imageFs info: %v", err)
	}
	infos, err := getCadvisorContainerInfo(p.cadvisor)
	if err != nil {
		return nil, fmt.Errorf("failed to get container info from cadvisor: %v", err)
	}
	// removeTerminatedContainerInfo will also remove pod level cgroups, so save the infos into allInfos first
	allInfos := infos
	klog.Warningf("in cadvisor_stats_provider.ListPodStats allInfos: %+v", allInfos)

	infos = removeTerminatedContainerInfo(infos)

	klog.Warningf("in cadvisor_stats_provider.ListPodStats (after removeterminatedcontainerinfo) infos: %+v", infos)

	// Map each container to a pod and update the PodStats with container data.
	podToStats := map[statsapi.PodReference]*statsapi.PodStats{}
	for key, cinfo := range infos {
		// On systemd using devicemapper each mount into the container has an
		// associated cgroup. We ignore them to ensure we do not get duplicate
		// entries in our summary. For details on .mount units:
		// http://man7.org/linux/man-pages/man5/systemd.mount.5.html
		if strings.HasSuffix(key, ".mount") {
			continue
		}
		// Build the Pod key if this container is managed by a Pod
		if !isPodManagedContainer(&cinfo) {
			continue
		}
		ref := buildPodRef(cinfo.Spec.Labels)

		// Lookup the PodStats for the pod using the PodRef. If none exists,
		// initialize a new entry.
		podStats, found := podToStats[ref]
		if !found {
			podStats = &statsapi.PodStats{PodRef: ref}
			podToStats[ref] = podStats
		}

		// Update the PodStats entry with the stats from the container by
		// adding it to podStats.Containers.
		containerName := kubetypes.GetContainerName(cinfo.Spec.Labels)
		klog.Warningf("in cadvisor_stats_provider.ListPodStats looking up stats for containerName: %+v", containerName)

		if containerName == leaky.PodInfraContainerName {
			// Special case for infrastructure container which is hidden from
			// the user and has network stats.
			klog.Warningf("in cadvisor_stats_provider.ListPodStats not appending container: %+v", containerName)

			podStats.Network = cadvisorInfoToNetworkStats("pod:"+ref.Namespace+"_"+ref.Name, &cinfo)
		} else {
			klog.Warningf("in cadvisor_stats_provider.ListPodStats appending container: %+v", containerName)

			podStats.Containers = append(podStats.Containers, *cadvisorInfoToContainerStats(containerName, &cinfo, &rootFsInfo, &imageFsInfo))
		}
	}

	// Add each PodStats to the result.
	result := make([]statsapi.PodStats, 0, len(podToStats))
	for _, podStats := range podToStats {
		// Lookup the volume stats for each pod.
		podUID := types.UID(podStats.PodRef.UID)
		var ephemeralStats []statsapi.VolumeStats
		if vstats, found := p.resourceAnalyzer.GetPodVolumeStats(podUID); found {
			ephemeralStats = make([]statsapi.VolumeStats, len(vstats.EphemeralVolumes))
			copy(ephemeralStats, vstats.EphemeralVolumes)
			podStats.VolumeStats = append(vstats.EphemeralVolumes, vstats.PersistentVolumes...)
		}
		podStats.EphemeralStorage = calcEphemeralStorage(podStats.Containers, ephemeralStats, &rootFsInfo, nil, false)
		// Lookup the pod-level cgroup's CPU and memory stats
		podInfo := getCadvisorPodInfoFromPodUID(podUID, allInfos)
		if podInfo != nil {
			cpu, memory := cadvisorInfoToCPUandMemoryStats(podInfo)
			podStats.CPU = cpu
			podStats.Memory = memory
		}

		status, found := p.statusProvider.GetPodStatus(podUID)
		if found && status.StartTime != nil && !status.StartTime.IsZero() {
			podStats.StartTime = *status.StartTime
			// only append stats if we were able to get the start time of the pod
			result = append(result, *podStats)
		}
	}

	klog.Warningf("in cadvisor_stats_provider.ListPodStats result: %+v", result)

	return result, nil
}

// ListPodStatsAndUpdateCPUNanoCoreUsage updates the cpu nano core usage for
// the containers and returns the stats for all the pod-managed containers.
// For cadvisor, cpu nano core usages are pre-computed and cached, so this
// function simply calls ListPodStats.
func (p *cadvisorStatsProvider) ListPodStatsAndUpdateCPUNanoCoreUsage() ([]statsapi.PodStats, error) {
	klog.Warningf("in cadvisor_stats_provider.ListPodStatsAndUpdateCPUNanoCoreUsage: %+v", p)

	return p.ListPodStats()
}

// ListPodCPUAndMemoryStats returns the cpu and memory stats of all the pod-managed containers.
func (p *cadvisorStatsProvider) ListPodCPUAndMemoryStats() ([]statsapi.PodStats, error) {
	klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats: %+v", p)

	infos, err := getCadvisorContainerInfo(p.cadvisor)
	klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats (before remove terminated) infos: %+v", infos)

	if err != nil {
		return nil, fmt.Errorf("failed to get container info from cadvisor: %v", err)
	}
	// removeTerminatedContainerInfo will also remove pod level cgroups, so save the infos into allInfos first
	allInfos := infos
	infos = removeTerminatedContainerInfo(infos)
	klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats (after remove terminated) infos: %+v", infos)

	// Map each container to a pod and update the PodStats with container data.
	podToStats := map[statsapi.PodReference]*statsapi.PodStats{}
	for key, cinfo := range infos {
		// On systemd using devicemapper each mount into the container has an
		// associated cgroup. We ignore them to ensure we do not get duplicate
		// entries in our summary. For details on .mount units:
		// http://man7.org/linux/man-pages/man5/systemd.mount.5.html
		if strings.HasSuffix(key, ".mount") {
			continue
		}
		// Build the Pod key if this container is managed by a Pod
		if !isPodManagedContainer(&cinfo) {
			continue
		}
		ref := buildPodRef(cinfo.Spec.Labels)

		// Lookup the PodStats for the pod using the PodRef. If none exists,
		// initialize a new entry.
		podStats, found := podToStats[ref]
		if !found {
			podStats = &statsapi.PodStats{PodRef: ref}
			podToStats[ref] = podStats
		}

		// Update the PodStats entry with the stats from the container by
		// adding it to podStats.Containers.
		containerName := kubetypes.GetContainerName(cinfo.Spec.Labels)
		klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats looking up containerName: %+v", containerName)

		if containerName == leaky.PodInfraContainerName {
			// Special case for infrastructure container which is hidden from
			// the user and has network stats.
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats not appending containerName: %+v", containerName)

			podStats.StartTime = metav1.NewTime(cinfo.Spec.CreationTime)
		} else {
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats appending ContainerName: %+v", containerName)

			podStats.Containers = append(podStats.Containers, *cadvisorInfoToContainerCPUAndMemoryStats(containerName, &cinfo))
		}
	}

	// Add each PodStats to the result.
	result := make([]statsapi.PodStats, 0, len(podToStats))
	for _, podStats := range podToStats {
		podUID := types.UID(podStats.PodRef.UID)
		// Lookup the pod-level cgroup's CPU and memory stats
		podInfo := getCadvisorPodInfoFromPodUID(podUID, allInfos)
		if podInfo != nil {
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats podInfo (not nil): %+v", podInfo)

			cpu, memory := cadvisorInfoToCPUandMemoryStats(podInfo)
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats: %+v", p)

			podStats.CPU = cpu
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats podStats.CPU: %+v", cpu)

			podStats.Memory = memory
			klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats podStats.Memory: %+v", memory)

		}
		result = append(result, *podStats)
	}
	klog.Warningf("in cadvisor_stats_provider.ListPodCPUAndMemoryStats result: %+v", result)


	return result, nil
}

// ImageFsStats returns the stats of the filesystem for storing images.
func (p *cadvisorStatsProvider) ImageFsStats() (*statsapi.FsStats, error) {
	imageFsInfo, err := p.cadvisor.ImagesFsInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to get imageFs info: %v", err)
	}
	imageStats, err := p.imageService.ImageStats()
	if err != nil || imageStats == nil {
		return nil, fmt.Errorf("failed to get image stats: %v", err)
	}

	var imageFsInodesUsed *uint64
	if imageFsInfo.Inodes != nil && imageFsInfo.InodesFree != nil {
		imageFsIU := *imageFsInfo.Inodes - *imageFsInfo.InodesFree
		imageFsInodesUsed = &imageFsIU
	}

	return &statsapi.FsStats{
		Time:           metav1.NewTime(imageFsInfo.Timestamp),
		AvailableBytes: &imageFsInfo.Available,
		CapacityBytes:  &imageFsInfo.Capacity,
		UsedBytes:      &imageStats.TotalStorageBytes,
		InodesFree:     imageFsInfo.InodesFree,
		Inodes:         imageFsInfo.Inodes,
		InodesUsed:     imageFsInodesUsed,
	}, nil
}

// ImageFsDevice returns name of the device where the image filesystem locates,
// e.g. /dev/sda1.
func (p *cadvisorStatsProvider) ImageFsDevice() (string, error) {
	imageFsInfo, err := p.cadvisor.ImagesFsInfo()
	if err != nil {
		return "", err
	}
	return imageFsInfo.Device, nil
}

// buildPodRef returns a PodReference that identifies the Pod managing cinfo
func buildPodRef(containerLabels map[string]string) statsapi.PodReference {
	podName := kubetypes.GetPodName(containerLabels)
	podNamespace := kubetypes.GetPodNamespace(containerLabels)
	podUID := kubetypes.GetPodUID(containerLabels)
	return statsapi.PodReference{Name: podName, Namespace: podNamespace, UID: podUID}
}

// isPodManagedContainer returns true if the cinfo container is managed by a Pod
func isPodManagedContainer(cinfo *cadvisorapiv2.ContainerInfo) bool {
	podName := kubetypes.GetPodName(cinfo.Spec.Labels)
	podNamespace := kubetypes.GetPodNamespace(cinfo.Spec.Labels)
	managed := podName != "" && podNamespace != ""
	if !managed && podName != podNamespace {
		klog.Warningf(
			"Expect container to have either both podName (%s) and podNamespace (%s) labels, or neither.",
			podName, podNamespace)
	}
	return managed
}

// getCadvisorPodInfoFromPodUID returns a pod cgroup information by matching the podUID with its CgroupName identifier base name
func getCadvisorPodInfoFromPodUID(podUID types.UID, infos map[string]cadvisorapiv2.ContainerInfo) *cadvisorapiv2.ContainerInfo {
	for key, info := range infos {
		if cm.IsSystemdStyleName(key) {
			// Convert to internal cgroup name and take the last component only.
			internalCgroupName := cm.ParseSystemdToCgroupName(key)
			key = internalCgroupName[len(internalCgroupName)-1]
		} else {
			// Take last component only.
			key = path.Base(key)
		}
		if cm.GetPodCgroupNameSuffix(podUID) == key {
			return &info
		}
	}
	return nil
}

// removeTerminatedContainerInfo returns the specified containerInfo but with
// the stats of the terminated containers removed.
//
// A ContainerInfo is considered to be of a terminated container if it has an
// older CreationTime and zero CPU instantaneous and memory RSS usage.
func removeTerminatedContainerInfo(containerInfo map[string]cadvisorapiv2.ContainerInfo) map[string]cadvisorapiv2.ContainerInfo {
	klog.Warningf("in cadvisor_stats_provider.removeTerminatedContainerInfo containerinfo passed in: %+v", containerInfo)

	cinfoMap := make(map[containerID][]containerInfoWithCgroup)
	for key, cinfo := range containerInfo {
		if !isPodManagedContainer(&cinfo) {
			continue
		}
		cinfoID := containerID{
			podRef:        buildPodRef(cinfo.Spec.Labels),
			containerName: kubetypes.GetContainerName(cinfo.Spec.Labels),
		}
		klog.Warningf("in cadvisor_stats_provider.removeTerminatedContainerInfo appending key %+v: %+v", key, cinfo)

		cinfoMap[cinfoID] = append(cinfoMap[cinfoID], containerInfoWithCgroup{
			cinfo:  cinfo,
			cgroup: key,
		})
	}
	result := make(map[string]cadvisorapiv2.ContainerInfo)
	for _, refs := range cinfoMap {
		if len(refs) == 1 {
			result[refs[0].cgroup] = refs[0].cinfo
			continue
		}
		sort.Sort(ByCreationTime(refs))
		i := 0
		klog.Warningf("in cadvisor_stats_provider.removeTerminatedContainerInfo refs: %+v", refs)

		for ; i < len(refs); i++ {
			if hasMemoryAndCPUInstUsage(&refs[i].cinfo) {
				// Stops removing when we first see an info with non-zero
				// CPU/Memory usage.
				klog.Warningf("in cadvisor_stats_provider.removeTerminatedContainerInfo found nonzero info: %+v", refs[i].cinfo)

				break
			}
		}
		for ; i < len(refs); i++ {
			result[refs[i].cgroup] = refs[i].cinfo
		}
	}
	klog.Warningf("in cadvisor_stats_provider.removeTerminatedContainerInfo result: %+v", result)

	return result
}

// ByCreationTime implements sort.Interface for []containerInfoWithCgroup based
// on the cinfo.Spec.CreationTime field.
type ByCreationTime []containerInfoWithCgroup

func (a ByCreationTime) Len() int      { return len(a) }
func (a ByCreationTime) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByCreationTime) Less(i, j int) bool {
	if a[i].cinfo.Spec.CreationTime.Equal(a[j].cinfo.Spec.CreationTime) {
		// There shouldn't be two containers with the same name and/or the same
		// creation time. However, to make the logic here robust, we break the
		// tie by moving the one without CPU instantaneous or memory RSS usage
		// to the beginning.
		return hasMemoryAndCPUInstUsage(&a[j].cinfo)
	}
	return a[i].cinfo.Spec.CreationTime.Before(a[j].cinfo.Spec.CreationTime)
}

// containerID is the identity of a container in a pod.
type containerID struct {
	podRef        statsapi.PodReference
	containerName string
}

// containerInfoWithCgroup contains the ContainerInfo and its cgroup name.
type containerInfoWithCgroup struct {
	cinfo  cadvisorapiv2.ContainerInfo
	cgroup string
}

// hasMemoryAndCPUInstUsage returns true if the specified container info has
// both non-zero CPU instantaneous usage and non-zero memory RSS usage, and
// false otherwise.
func hasMemoryAndCPUInstUsage(info *cadvisorapiv2.ContainerInfo) bool {
	klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage: %+v", info)

	if !info.Spec.HasCpu || !info.Spec.HasMemory {
		klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage, either no memory or cpu usage found on info.Spec, hasCPU %+v , hasMemory %+v", info.Spec.HasCpu, info.Spec.HasMemory)
		return false
	}
	cstat, found := latestContainerStats(info)
	klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage cstat: %+v", cstat)
	if !found {
		klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage !found, returning false: %+v", found)
		return false
	}
	if cstat.CpuInst == nil {
		klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage cpuInst is nil, returing false: %+v", cstat.CpuInst)
		return false
	}
	klog.Warningf("in cadvisor_stats_provider.hasMemoryAndCPUInstUsage cstat.CpuInst.Usage.Total: %+v, cstat.Memory.Rss: %+v", cstat.CpuInst.Usage.Total, cstat.Memory.RSS)
	return cstat.CpuInst.Usage.Total != 0 && cstat.Memory.RSS != 0
}

func getCadvisorContainerInfo(ca cadvisor.Interface) (map[string]cadvisorapiv2.ContainerInfo, error) {
	klog.Warningf("in cadvisor_stats_provider.getCadvisorContainerInfo: %+v", ca)
	infos, err := ca.ContainerInfoV2("/", cadvisorapiv2.RequestOptions{
		IdType:    cadvisorapiv2.TypeName,
		Count:     2, // 2 samples are needed to compute "instantaneous" CPU
		Recursive: true,
	})
	klog.Warningf("in cadvisor_stats_provider.getCadvisorContainerInfo infos: %+v", infos)

	if err != nil {
		if _, ok := infos["/"]; ok {
			// If the failure is partial, log it and return a best-effort
			// response.
			klog.Errorf("Partial failure issuing cadvisor.ContainerInfoV2: %v", err)
		} else {
			return nil, fmt.Errorf("failed to get root cgroup stats: %v", err)
		}
	}
	return infos, nil
}
