/*
Copyright 2018 The Vulcan Authors.

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

package algorithm

import (

	"github.com/kubernetes-sigs/kube-batch/pkg/scheduler/framework"
	// Import default actions/plugins.
	_ "github.com/kubernetes-sigs/kube-batch/pkg/scheduler/actions"
	_ "github.com/kubernetes-sigs/kube-batch/pkg/scheduler/plugins"

	"volcano.sh/volcano/pkg/scheduler/algorithm/fairshare"
)

func init() {
	framework.RegisterPluginBuilder("fairshare", fairshare.New)
}

