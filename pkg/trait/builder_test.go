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
	"testing"

	"github.com/apache/camel-k/pkg/builder"

	"github.com/apache/camel-k/pkg/util/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestBuilderTraitNotAppliedBecauseOfNilContext(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I),
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko),
	}

	for _, e := range environments {
		e.Context = nil

		t.Run(string(e.Platform.Spec.Cluster), func(t *testing.T) {
			err := NewCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.NotContains(t, e.ExecutedTraits, ID("builder"))
			assert.Empty(t, e.Steps)
		})
	}
}

func TestBuilderTraitNotAppliedBecauseOfNilPhase(t *testing.T) {
	environments := []*Environment{
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I),
		createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko),
	}

	for _, e := range environments {
		e.Context.Status.Phase = ""

		t.Run(string(e.Platform.Spec.Cluster), func(t *testing.T) {
			err := NewCatalog().apply(e)

			assert.Nil(t, err)
			assert.NotEmpty(t, e.ExecutedTraits)
			assert.NotContains(t, e.ExecutedTraits, ID("builder"))
			assert.Empty(t, e.Steps)
		})
	}
}

func TestS2IBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterOpenShift, v1alpha1.IntegrationPlatformBuildPublishStrategyS2I)
	err := NewCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.Contains(t, env.ExecutedTraits, ID("builder"))
	assert.NotEmpty(t, env.Steps)
	assert.Len(t, env.Steps, 5)
	assert.Condition(t, func() bool {
		for _, s := range env.Steps {
			if s.ID() == "publisher/s2i" && s.Phase() == builder.ApplicationPublishPhase {
				return true
			}
		}

		return false
	})
}

func TestKanikoBuilderTrait(t *testing.T) {
	env := createBuilderTestEnv(v1alpha1.IntegrationPlatformClusterKubernetes, v1alpha1.IntegrationPlatformBuildPublishStrategyKaniko)
	err := NewCatalog().apply(env)

	assert.Nil(t, err)
	assert.NotEmpty(t, env.ExecutedTraits)
	assert.Contains(t, env.ExecutedTraits, ID("builder"))
	assert.NotEmpty(t, env.Steps)
	assert.Len(t, env.Steps, 5)
	assert.Condition(t, func() bool {
		for _, s := range env.Steps {
			if s.ID() == "publisher/kaniko" && s.Phase() == builder.ApplicationPublishPhase {
				return true
			}
		}

		return false
	})
}

func createBuilderTestEnv(cluster v1alpha1.IntegrationPlatformCluster, strategy v1alpha1.IntegrationPlatformBuildPublishStrategy) *Environment {
	return &Environment{
		Integration: &v1alpha1.Integration{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "test",
				Namespace: "ns",
			},
			Status: v1alpha1.IntegrationStatus{
				Phase: v1alpha1.IntegrationPhaseDeploying,
			},
		},
		Context: &v1alpha1.IntegrationContext{
			Status: v1alpha1.IntegrationContextStatus{
				Phase: v1alpha1.IntegrationContextPhaseBuilding,
			},
		},
		Platform: &v1alpha1.IntegrationPlatform{
			Spec: v1alpha1.IntegrationPlatformSpec{
				Cluster: cluster,
				Build: v1alpha1.IntegrationPlatformBuildSpec{
					PublishStrategy: strategy,
					Registry:        "registry",
				},
			},
		},
		ExecutedTraits: make([]ID, 0),
		Resources:      kubernetes.NewCollection(),
	}
}
