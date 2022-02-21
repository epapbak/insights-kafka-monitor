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

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Messages
const (
	versionMessage                 = "Insights Kafka monitor version 1.0"
	authorsMessage                 = "Pavel Tisnovsky, Red Hat Inc."
	connectionToBrokerMessage      = "Connection to broker"
	operationFailedMessage         = "Operation failed"
	notConnectedToBrokerMessage    = "Not connected to broker"
	brokerConnectionSuccessMessage = "Broker connection OK"
	brokerAddress                  = "Broker address"
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

// showVersion function displays version information.
func showVersion() {
	fmt.Println(versionMessage)
}

// showAuthors function displays information about authors.
func showAuthors() {
	fmt.Println(authorsMessage)
}

// showConfiguration function displays actual configuration.
func showConfiguration(config ConfigStruct) {
	brokerConfig := GetBrokerConfiguration(config)
	log.Info().
		Str(brokerAddress, brokerConfig.Address).
		Str("Topic", brokerConfig.Topic).
		Str("Group", brokerConfig.Group).
		Bool("Enabled", brokerConfig.Enabled).
		Msg("Broker configuration")

	loggingConfig := GetLoggingConfiguration(config)
	log.Info().
		Str("Level", loggingConfig.LogLevel).
		Bool("Pretty colored debug logging", loggingConfig.Debug).
		Msg("Logging configuration")
}

// doSelectedOperation function perform operation selected on command line.
// When no operation is specified, the Insights Kafka monitor service is
// started instead.
func doSelectedOperation(configuration ConfigStruct, cliFlags CliFlags) (int, error) {
	switch {
	case cliFlags.ShowVersion:
		showVersion()
		return ExitStatusOK, nil
	case cliFlags.ShowAuthors:
		showAuthors()
		return ExitStatusOK, nil
	case cliFlags.ShowConfiguration:
		showConfiguration(configuration)
		return ExitStatusOK, nil
	case cliFlags.CheckConnectionToKafka:
		// return tryToConnectToKafka(configuration)
	default:
		// exitCode, err := startService(configuration)
		// return exitCode, err
	}
	// this can not happen: return ExitStatusOK, nil
	return ExitStatusOK, nil
}

// main function is entry point to the Kafka monitor service.
func main() {
	var cliFlags CliFlags

	// define and parse all command line options
	flag.BoolVar(&cliFlags.ShowVersion, "version", false, "show version")
	flag.BoolVar(&cliFlags.ShowAuthors, "authors", false, "show authors")
	flag.BoolVar(&cliFlags.ShowConfiguration, "show-configuration", false, "show configuration")
	flag.BoolVar(&cliFlags.CheckConnectionToKafka, "check-kafka", false, "check connection to Kafka")
	flag.Parse()

	// config has exactly the same structure as *.toml file
	config, err := LoadConfiguration(configFileEnvVariableName, defaultConfigFileName)
	if err != nil {
		log.Err(err).Msg("Load configuration")
	}

	if config.Logging.Debug {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	log.Debug().Msg("Started")

	// perform selected operation
	exitStatus, err := doSelectedOperation(config, cliFlags)
	if err != nil {
		log.Err(err).Msg("Do selected operation")
		os.Exit(exitStatus)
		return
	}

	log.Debug().Msg("Finished")
}
