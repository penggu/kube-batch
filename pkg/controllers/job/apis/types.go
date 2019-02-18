/*
Copyright 2019 The Volcano Authors.

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

package apis

import (
	"fmt"

	"k8s.io/api/core/v1"

	"hpw.cloud/volcano/pkg/apis/batch/v1alpha1"
)

type JobInfo struct {
	Namespace string
	Name      string

	Job  *v1alpha1.Job
	Pods map[string]map[string]*v1.Pod
}

func (ji *JobInfo) Clone() *JobInfo {
	job := &JobInfo{
		Namespace: ji.Namespace,
		Name:      ji.Name,
		Job:       ji.Job,

		Pods: make(map[string]map[string]*v1.Pod),
	}

	for key, pods := range ji.Pods {
		job.Pods[key] = make(map[string]*v1.Pod)
		for pn, pod := range pods {
			job.Pods[key][pn] = pod
		}
	}

	return job
}

func (ji *JobInfo) SetJob(job *v1alpha1.Job) {
	ji.Name = job.Name
	ji.Namespace = job.Namespace
	ji.Job = job
}

func (ji *JobInfo) AddPod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if _, found := ji.Pods[taskName]; !found {
		ji.Pods[taskName] = make(map[string]*v1.Pod)
	}
	if _, found := ji.Pods[taskName][pod.Name]; found {
		return fmt.Errorf("duplicated pod")
	}
	ji.Pods[taskName][pod.Name] = pod

	return nil
}

func (ji *JobInfo) UpdatePod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if _, found := ji.Pods[taskName]; !found {
		return fmt.Errorf("can not find task %s in cache", taskName)
	}
	if _, found := ji.Pods[taskName][pod.Name]; !found {
		return fmt.Errorf("can not find pod <%s/%s> in cache",
			pod.Namespace, pod.Name)
	}
	ji.Pods[taskName][pod.Name] = pod

	return nil
}

func (ji *JobInfo) DeletePod(pod *v1.Pod) error {
	taskName, found := pod.Annotations[v1alpha1.TaskSpecKey]
	if !found {
		return fmt.Errorf("failed to taskName of Pod <%s/%s>",
			pod.Namespace, pod.Name)
	}

	if pods, found := ji.Pods[taskName]; found {
		delete(pods, pod.Name)
		if len(pods) == 0 {
			delete(ji.Pods, taskName)
		}
	}

	return nil
}

type Request struct {
	Namespace string
	JobName   string
	TaskName  string

	Event  v1alpha1.Event
	Action v1alpha1.Action
}
