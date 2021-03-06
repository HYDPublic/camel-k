/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package trait

import (
	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/builder/kaniko"
	"github.com/apache/camel-k/pkg/builder/s2i"
	"github.com/apache/camel-k/pkg/platform"
)

// TODO: we should add a way to label a trait as platform so it cannot be disabled/removed
type builderTrait struct {
	BaseTrait `property:",squash"`
}

func newBuilderTrait() *builderTrait {
	return &builderTrait{
		BaseTrait: newBaseTrait("builder"),
	}
}

func (*builderTrait) appliesTo(e *Environment) bool {
	return e.Context != nil && e.Context.Status.Phase == v1alpha1.IntegrationContextPhaseBuilding
}

func (*builderTrait) apply(e *Environment) error {
	if platform.SupportsS2iPublishStrategy(e.Platform) {
		e.Steps = s2i.DefaultSteps
	} else if platform.SupportsKanikoPublishStrategy(e.Platform) {
		e.Steps = kaniko.DefaultSteps
	}

	return nil
}
