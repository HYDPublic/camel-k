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
	"reflect"
	"strings"

	"github.com/apache/camel-k/pkg/apis/camel/v1alpha1"
	"github.com/apache/camel-k/pkg/platform"
	"github.com/fatih/structs"
)

// Catalog collects all information about traits in one place
type Catalog struct {
	tDependencies Trait
	tDeployment   Trait
	tKnative      Trait
	tService      Trait
	tRoute        Trait
	tIngress      Trait
	tOwner        Trait
	tBuilder      Trait
}

// NewCatalog creates a new trait Catalog
func NewCatalog() *Catalog {
	return &Catalog{
		tDependencies: newDependenciesTrait(),
		tDeployment:   newDeploymentTrait(),
		tKnative:      newKnativeTrait(),
		tService:      newServiceTrait(),
		tRoute:        newRouteTrait(),
		tIngress:      newIngressTrait(),
		tOwner:        newOwnerTrait(),
		tBuilder:      newBuilderTrait(),
	}
}

func (c *Catalog) allTraits() []Trait {
	return []Trait{
		c.tDependencies,
		c.tDeployment,
		c.tKnative,
		c.tService,
		c.tRoute,
		c.tIngress,
		c.tOwner,
		c.tBuilder,
	}
}

func (c *Catalog) traitsFor(environment *Environment) []Trait {
	profile := platform.GetProfile(environment.Platform)
	if environment.Context != nil && environment.Context.Spec.Profile != "" {
		profile = environment.Context.Spec.Profile
	}
	if environment.Integration != nil && environment.Integration.Spec.Profile != "" {
		profile = environment.Integration.Spec.Profile
	}

	switch profile {
	case v1alpha1.TraitProfileOpenShift:
		return []Trait{
			c.tDependencies,
			c.tDeployment,
			c.tService,
			c.tRoute,
			c.tOwner,
			c.tBuilder,
		}
	case v1alpha1.TraitProfileKubernetes:
		return []Trait{
			c.tDependencies,
			c.tDeployment,
			c.tService,
			c.tIngress,
			c.tOwner,
			c.tBuilder,
		}
	case v1alpha1.TraitProfileKnative:
		return []Trait{
			c.tDependencies,
			c.tKnative,
			c.tOwner,
			c.tBuilder,
		}
	}

	return nil
}

func (c *Catalog) apply(environment *Environment) error {
	c.configure(environment)
	traits := c.traitsFor(environment)
	for _, trait := range traits {
		if !trait.appliesTo(environment) {
			continue
		}

		if trait.IsAuto() {
			if err := trait.autoconfigure(environment); err != nil {
				return err
			}
		}
		if trait.IsEnabled() {
			if err := trait.apply(environment); err != nil {
				return err
			}
			environment.ExecutedTraits = append(environment.ExecutedTraits, trait.ID())
		}
	}
	return nil
}

// GetTrait returns the trait with the given ID
func (c *Catalog) GetTrait(id string) Trait {
	for _, t := range c.allTraits() {
		if t.ID() == ID(id) {
			return t
		}
	}
	return nil
}

func (c *Catalog) configure(env *Environment) {
	if env.Integration == nil || env.Integration.Spec.Traits == nil {
		return
	}
	for id, traitSpec := range env.Integration.Spec.Traits {
		catTrait := c.GetTrait(id)
		if catTrait != nil {
			traitSpec.Decode(catTrait)
		}
	}
}

// ComputeTraitsProperties returns all key/value configuration properties that can be used to configure traits
func (c *Catalog) ComputeTraitsProperties() []string {
	results := make([]string, 0)
	for _, trait := range c.allTraits() {
		c.processFields(structs.Fields(trait), func(name string) {
			results = append(results, string(trait.ID())+"."+name)
		})
	}

	return results
}

func (c *Catalog) processFields(fields []*structs.Field, processor func(string)) {
	for _, f := range fields {
		if f.IsEmbedded() && f.IsExported() && f.Kind() == reflect.Struct {
			c.processFields(f.Fields(), processor)
		}

		if f.IsEmbedded() {
			continue
		}

		property := f.Tag("property")

		if property != "" {
			items := strings.Split(property, ",")
			processor(items[0])
		}
	}
}
