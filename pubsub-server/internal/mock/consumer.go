// Code generated by MockGen. DO NOT EDIT.
// Source: pubsub-server/internal/consumer.go

// Package mock_internal is a generated GoMock package.
package mock_internal

import (
	reflect "reflect"

	internal "github.com/getjump/pubsub/pubsub-server/internal"
	gomock "github.com/golang/mock/gomock"
)

// MockConsumer is a mock of Consumer interface.
type MockConsumer struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMockRecorder
}

// MockConsumerMockRecorder is the mock recorder for MockConsumer.
type MockConsumerMockRecorder struct {
	mock *MockConsumer
}

// NewMockConsumer creates a new mock instance.
func NewMockConsumer(ctrl *gomock.Controller) *MockConsumer {
	mock := &MockConsumer{ctrl: ctrl}
	mock.recorder = &MockConsumerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumer) EXPECT() *MockConsumerMockRecorder {
	return m.recorder
}

// GetID mocks base method.
func (m *MockConsumer) GetID() uint64 {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(uint64)
	return ret0
}

// GetID indicates an expected call of GetID.
func (mr *MockConsumerMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockConsumer)(nil).GetID))
}

// GetMessageChan mocks base method.
func (m *MockConsumer) GetMessageChan() chan internal.ConsumerMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageChan")
	ret0, _ := ret[0].(chan internal.ConsumerMessage)
	return ret0
}

// GetMessageChan indicates an expected call of GetMessageChan.
func (mr *MockConsumerMockRecorder) GetMessageChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageChan", reflect.TypeOf((*MockConsumer)(nil).GetMessageChan))
}

// ShutdownChan mocks base method.
func (m *MockConsumer) ShutdownChan() chan struct{} {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShutdownChan")
	ret0, _ := ret[0].(chan struct{})
	return ret0
}

// ShutdownChan indicates an expected call of ShutdownChan.
func (mr *MockConsumerMockRecorder) ShutdownChan() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShutdownChan", reflect.TypeOf((*MockConsumer)(nil).ShutdownChan))
}

// MockConsumerFactory is a mock of ConsumerFactory interface.
type MockConsumerFactory struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerFactoryMockRecorder
}

// MockConsumerFactoryMockRecorder is the mock recorder for MockConsumerFactory.
type MockConsumerFactoryMockRecorder struct {
	mock *MockConsumerFactory
}

// NewMockConsumerFactory creates a new mock instance.
func NewMockConsumerFactory(ctrl *gomock.Controller) *MockConsumerFactory {
	mock := &MockConsumerFactory{ctrl: ctrl}
	mock.recorder = &MockConsumerFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerFactory) EXPECT() *MockConsumerFactoryMockRecorder {
	return m.recorder
}

// NewConsumer mocks base method.
func (m *MockConsumerFactory) NewConsumer() internal.Consumer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConsumer")
	ret0, _ := ret[0].(internal.Consumer)
	return ret0
}

// NewConsumer indicates an expected call of NewConsumer.
func (mr *MockConsumerFactoryMockRecorder) NewConsumer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConsumer", reflect.TypeOf((*MockConsumerFactory)(nil).NewConsumer))
}

// MockConsumerMatcher is a mock of ConsumerMatcher interface.
type MockConsumerMatcher struct {
	ctrl     *gomock.Controller
	recorder *MockConsumerMatcherMockRecorder
}

// MockConsumerMatcherMockRecorder is the mock recorder for MockConsumerMatcher.
type MockConsumerMatcherMockRecorder struct {
	mock *MockConsumerMatcher
}

// NewMockConsumerMatcher creates a new mock instance.
func NewMockConsumerMatcher(ctrl *gomock.Controller) *MockConsumerMatcher {
	mock := &MockConsumerMatcher{ctrl: ctrl}
	mock.recorder = &MockConsumerMatcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConsumerMatcher) EXPECT() *MockConsumerMatcherMockRecorder {
	return m.recorder
}

// AddConsumers mocks base method.
func (m *MockConsumerMatcher) AddConsumers(topic string, consumers ...internal.Consumer) {
	m.ctrl.T.Helper()
	varargs := []interface{}{topic}
	for _, a := range consumers {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "AddConsumers", varargs...)
}

// AddConsumers indicates an expected call of AddConsumers.
func (mr *MockConsumerMatcherMockRecorder) AddConsumers(topic interface{}, consumers ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{topic}, consumers...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConsumers", reflect.TypeOf((*MockConsumerMatcher)(nil).AddConsumers), varargs...)
}

// MatchTopic mocks base method.
func (m *MockConsumerMatcher) MatchTopic(topic string) []internal.Consumer {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MatchTopic", topic)
	ret0, _ := ret[0].([]internal.Consumer)
	return ret0
}

// MatchTopic indicates an expected call of MatchTopic.
func (mr *MockConsumerMatcherMockRecorder) MatchTopic(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MatchTopic", reflect.TypeOf((*MockConsumerMatcher)(nil).MatchTopic), topic)
}

// RemoveConsumerFromTopic mocks base method.
func (m *MockConsumerMatcher) RemoveConsumerFromTopic(topic string, consumer internal.Consumer) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveConsumerFromTopic", topic, consumer)
}

// RemoveConsumerFromTopic indicates an expected call of RemoveConsumerFromTopic.
func (mr *MockConsumerMatcherMockRecorder) RemoveConsumerFromTopic(topic, consumer interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveConsumerFromTopic", reflect.TypeOf((*MockConsumerMatcher)(nil).RemoveConsumerFromTopic), topic, consumer)
}
