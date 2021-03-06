/**
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package org.apache.camel.component.knative;

import org.apache.camel.RuntimeCamelException;
import org.apache.camel.cloud.ServiceDefinition;
import org.apache.camel.spi.Metadata;
import org.apache.camel.spi.UriParam;

import java.util.Arrays;

import static org.apache.camel.util.CollectionHelper.mapOf;

public class KnativeConfiguration implements Cloneable {
    @UriParam
    @Metadata(required = "true")
    private KnativeEnvironment environment;

    public KnativeConfiguration() {
    }

    // ************************
    //
    // Properties
    //
    // ************************

    public KnativeEnvironment getEnvironment() {
        return environment;
    }

    /**
     * The environment
     */
    public void setEnvironment(KnativeEnvironment environment) {
        this.environment = environment;
    }

    // ************************
    //
    // Cloneable
    //
    // ************************

    public KnativeConfiguration copy() {
        try {
            return (KnativeConfiguration)super.clone();
        } catch (CloneNotSupportedException e) {
            throw new RuntimeCamelException(e);
        }
    }
}
