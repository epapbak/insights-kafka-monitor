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

package main_test

import (
	"context"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"

	main "github.com/RedHatInsights/insights-kafka-monitor"
)

// TestNewConsumerBadBroker function checks the consumer creation by
// using a non accessible Kafka broker.
func TestNewConsumerBadBroker(t *testing.T) {
	const expectedErr = "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)"

	// invalid broker configuration
	var brokerConfiguration = main.BrokerConfiguration{
		Address: "",
		Topic:   "whatever",
		Group:   "whatever",
		Enabled: true,
	}

	// try to construct new consumer
	mockConsumer, err := main.NewConsumer(brokerConfiguration, true)

	// check that error is really reported
	assert.EqualError(t, err, expectedErr)

	// test the return value
	assert.Equal(
		t,
		(*main.KafkaConsumer)(nil),
		mockConsumer,
		"consumer.New should return nil instead of Consumer implementation",
	)
}

// TestNewConsumerLocalBroker function checks the consumer creation by using a
// non accessible Kafka broker. This test assumes there is no local Kafka
// instance currently running
func TestNewConsumerLocalBroker(t *testing.T) {
	const expectedErr = "kafka: client has run out of available brokers to talk to (Is your cluster reachable?)"

	// valid broker configuration for local Kafka instance
	var brokerConfiguration = main.BrokerConfiguration{
		Address: "localhost:9092",
		Topic:   "platform.notifications.ingress",
		Group:   "",
		Enabled: true,
	}

	// try to construct new consumer
	mockConsumer, err := main.NewConsumer(brokerConfiguration, true)

	// check that error is really reported
	assert.EqualError(t, err, expectedErr)

	// test the return value
	assert.Equal(
		t,
		(*main.KafkaConsumer)(nil),
		mockConsumer,
		"consumer.New should return nil instead of Consumer implementation",
	)
}

// NewDummyConsumer function constructs new instance of (not running)
// KafkaConsumer.
func NewDummyConsumer() *main.KafkaConsumer {
	brokerCfg := main.BrokerConfiguration{
		Address: "localhost:1234",
		Topic:   "topic",
		Group:   "group",
	}
	return &main.KafkaConsumer{
		Configuration: brokerCfg,
		Verbose:       true,
		Ready:         make(chan bool),
	}
}

// TestProcessEmptyMessage check the behaviour of function ProcessMessage with
// empty message on input.
func TestProcessEmptyMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// prepare an empty message
	message := sarama.ConsumerMessage{}

	// try to process the message
	err := dummyConsumer.ProcessMessage(&message)

	// check for error - it should not be reported
	assert.Nil(t, err)

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())
}

// TestProcessCorrectMessage check the behaviour of function ProcessMessage for
// correct message.
func TestProcessCorrectMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// prepare a message
	message := sarama.ConsumerMessage{}

	// fill in a message payload
	ConsumerMessage := `{"foo": "bar"}`
	message.Value = []byte(ConsumerMessage)

	// message is correct
	err := dummyConsumer.ProcessMessage(&message)

	// check for error - it should not be reported
	assert.Nil(t, err)

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())
}

// TestConsumerSetup function checks the method KafkaConsumer.Setup().
func TestConsumerSetup(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// channel that needs to be closed in Setup
	dummyConsumer.Ready = make(chan bool)

	// try to setup the consumer (without consumer group)
	err := dummyConsumer.Setup(nil)

	// and check for any error
	assert.Nil(t, err)
}

// TestConsumerCleanup function checks the method KafkaConsumer.Cleanup().
func TestConsumerCleanup(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// try to cleanup the consumer (without consumer group)
	err := dummyConsumer.Cleanup(nil)

	// and check for any error
	assert.Nil(t, err)
}

// TestConsumerClose function checks the method KafkaConsumer.Close().
func TestConsumerClose(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// try to close the consumer (without consumer group)
	err := dummyConsumer.Close()

	// and check for any error
	assert.Nil(t, err)
}

// TestConsumerCloseCancel function checks the method KafkaConsumer.Close().
func TestConsumerCloseCancel(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// setup cancel hook
	_, cancel := context.WithCancel(context.Background())
	dummyConsumer.Cancel = cancel

	// try to close the consumer (without consumer group)
	err := dummyConsumer.Close()

	// and check for any error
	assert.Nil(t, err)
}

// TestHandleNilMessage function checks the method
// KafkaConsumer.HandleMessage() for nil input.
func TestHandleNilMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	// nil message
	dummyConsumer.HandleMessage(nil)
}

// TestHandleEmptyMessage function checks the method
// KafkaConsumer.HandleMessage() for empty message value.
func TestHandleEmptyMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	message := sarama.ConsumerMessage{}

	message.Value = []byte("")

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// empty message
	dummyConsumer.HandleMessage(&message)

	// counter checks
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())
}

// TestHandleCorrectMessage function checks the method
// KafkaConsumer.HandleMessage() for correct input.
func TestHandleCorrectMessage(t *testing.T) {
	// construct dummy consumer
	dummyConsumer := NewDummyConsumer()

	message := sarama.ConsumerMessage{}
	value := `[42`

	message.Value = []byte(value)

	// counter checks
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())

	// correct message
	dummyConsumer.HandleMessage(&message)

	// counter checks
	assert.Equal(t, uint64(1), dummyConsumer.GetNumberOfSuccessfullyConsumedMessages())
	assert.Equal(t, uint64(0), dummyConsumer.GetNumberOfErrorsConsumingMessages())
}
