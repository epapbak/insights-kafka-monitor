/*
Copyright © 2022 Red Hat, Inc.

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

package main

// Messages
const (
	versionMessage                 = "CCX Notification Writer version 1.0"
	authorsMessage                 = "Pavel Tisnovsky, Red Hat Inc."
	connectionToBrokerMessage      = "Connection to broker"
	operationFailedMessage         = "Operation failed"
	notConnectedToBrokerMessage    = "Not connected to broker"
	brokerConnectionSuccessMessage = "Broker connection OK"
)

// Configuration-related constants
const (
	configFileEnvVariableName = "INSIGHTS_KAFKA_MONITOR_CONFIG_FILE"
	defaultConfigFileName     = "config"
)

// Exit codes
const (
	// ExitStatusOK means that the tool finished with success
	ExitStatusOK = iota
	// ExitStatusError is a general error code
	ExitStatusError
	// ExitStatusConsumerError is returned in case of any consumer-related error
	ExitStatusConsumerError
	// ExitStatusKafkaError is returned in case of any Kafka-related error
	ExitStatusKafkaError
)
