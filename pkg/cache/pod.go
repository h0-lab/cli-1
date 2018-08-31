//
// Last.Backend LLC CONFIDENTIAL
// __________________
//
// [2014] - [2018] Last.Backend LLC
// All Rights Reserved.
//
// NOTICE:  All information contained herein is, and remains
// the property of Last.Backend LLC and its suppliers,
// if any.  The intellectual and technical concepts contained
// herein are proprietary to Last.Backend LLC
// and its suppliers and may be covered by Russian Federation and Foreign Patents,
// patents in process, and are protected by trade secret or copyright law.
// Dissemination of this information or reproduction of this material
// is strictly forbidden unless prior written permission is obtained
// from Last.Backend LLC.
//

package cache

import (
	"github.com/lastbackend/cli/pkg/distribution/types"
	"github.com/lastbackend/cli/pkg/log"
	"sync"
)

type PodCache struct {
	lock       sync.RWMutex
	stats      PodCacheStats
	containers map[string]*types.Container
	pods       map[string]*types.Pod
}

type PodCacheStats struct {
	pods       int
	containers int
}

func (ps *PodCache) GetPodsCount() int {
	log.V(logLevel).Debugf("Cache: PodCache: get pods count: %d", ps.stats.pods)
	return ps.stats.pods
}

func (ps *PodCache) GetContainersCount() int {
	log.V(logLevel).Debugf("Cache: PodCache: get containers count: %d", ps.stats.containers)
	return ps.stats.containers
}

func (ps *PodCache) GetPods() map[string]*types.Pod {
	log.V(logLevel).Debug("Cache: PodCache: get pods")
	return ps.pods
}

func (ps *PodCache) GetContainer(id string) *types.Container {
	log.V(logLevel).Debugf("Cache: PodCache: get container: %s", id)
	c, ok := ps.containers[id]
	if !ok {
		return nil
	}
	return c
}

func (ps *PodCache) AddContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: add container: %#v", c)
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.addContainer(c)

}

func (ps *PodCache) SetContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: set container: %#v", c)
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.setContainer(c)
}

func (ps *PodCache) DelContainer(c *types.Container) {
	log.V(logLevel).Debugf("Cache: PodCache: del container: %s", c.ID)
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.delContainer(c)
}

func (ps *PodCache) GetPod(id string) *types.Pod {
	log.V(logLevel).Debugf("Cache: PodCache: get pod: %s", id)
	ps.lock.Lock()
	defer ps.lock.Unlock()
	pod, ok := ps.pods[id]
	if !ok {
		return nil
	}
	return pod
}

func (ps *PodCache) AddPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: add pod: %#v", pod)
	ps.SetPod(pod)
}

func (ps *PodCache) SetPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: set pod: %#v", pod)
	ps.lock.Lock()
	defer ps.lock.Unlock()

	if p, ok := ps.pods[pod.Meta.Name]; ok {
		ps.delPod(p)
	}
	ps.addPod(pod)
}

func (ps *PodCache) DelPod(pod *types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: del pod: %#v", pod)
	ps.lock.Lock()
	defer ps.lock.Unlock()
	if p, ok := ps.pods[pod.Meta.Name]; ok {
		ps.delPod(p)
	}
}

func (ps *PodCache) SetPods(pods []*types.Pod) {
	log.V(logLevel).Debugf("Cache: PodCache: set pods: %#v", pods)
	for _, pod := range pods {
		ps.AddPod(pod)
	}
}

func (ps *PodCache) addPod(pod *types.Pod) {
	ps.pods[pod.Meta.Name] = pod
	ps.stats.pods++
}

func (ps *PodCache) delPod(pod *types.Pod) {
	delete(ps.pods, pod.Meta.Name)
	ps.stats.pods--
}

func (ps *PodCache) addContainer(c *types.Container) {
	if _, ok := ps.containers[c.ID]; !ok {
		ps.stats.containers++
	}
	ps.containers[c.ID] = c
}

func (ps *PodCache) setContainer(c *types.Container) {
	if _, ok := ps.containers[c.ID]; !ok {
		ps.stats.containers++
	}
	ps.containers[c.ID] = c
}

func (ps *PodCache) delContainer(c *types.Container) {
	if _, ok := ps.containers[c.ID]; ok {
		delete(ps.containers, c.ID)
		ps.stats.containers--
	}
}

func NewPodCache() *PodCache {
	log.V(logLevel).Debug("Cache: PodCache: initialization storage")

	return &PodCache{
		stats:      PodCacheStats{},
		containers: make(map[string]*types.Container),
		pods:       make(map[string]*types.Pod),
	}
}
