// Code generated by MockGen. DO NOT EDIT.
// Source: cmd/list_connectors.go

// Package mock_cmd is a generated GoMock package.
package mock_cmd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	meroxa "github.com/meroxa/meroxa-go"
)

// MockListConnectorsClient is a mock of ListConnectorsClient interface.
type MockListConnectorsClient struct {
	ctrl     *gomock.Controller
	recorder *MockListConnectorsClientMockRecorder
}

// MockListConnectorsClientMockRecorder is the mock recorder for MockListConnectorsClient.
type MockListConnectorsClientMockRecorder struct {
	mock *MockListConnectorsClient
}

// NewMockListConnectorsClient creates a new mock instance.
func NewMockListConnectorsClient(ctrl *gomock.Controller) *MockListConnectorsClient {
	mock := &MockListConnectorsClient{ctrl: ctrl}
	mock.recorder = &MockListConnectorsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockListConnectorsClient) EXPECT() *MockListConnectorsClientMockRecorder {
	return m.recorder
}

// GetPipelineByName mocks base method.
func (m *MockListConnectorsClient) GetPipelineByName(ctx context.Context, name string) (*meroxa.Pipeline, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPipelineByName", ctx, name)
	ret0, _ := ret[0].(*meroxa.Pipeline)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPipelineByName indicates an expected call of GetPipelineByName.
func (mr *MockListConnectorsClientMockRecorder) GetPipelineByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPipelineByName", reflect.TypeOf((*MockListConnectorsClient)(nil).GetPipelineByName), ctx, name)
}

// ListConnectors mocks base method.
func (m *MockListConnectorsClient) ListConnectors(ctx context.Context) ([]*meroxa.Connector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListConnectors", ctx)
	ret0, _ := ret[0].([]*meroxa.Connector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListConnectors indicates an expected call of ListConnectors.
func (mr *MockListConnectorsClientMockRecorder) ListConnectors(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListConnectors", reflect.TypeOf((*MockListConnectorsClient)(nil).ListConnectors), ctx)
}

// ListPipelineConnectors mocks base method.
func (m *MockListConnectorsClient) ListPipelineConnectors(ctx context.Context, pipelineID int) ([]*meroxa.Connector, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListPipelineConnectors", ctx, pipelineID)
	ret0, _ := ret[0].([]*meroxa.Connector)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListPipelineConnectors indicates an expected call of ListPipelineConnectors.
func (mr *MockListConnectorsClientMockRecorder) ListPipelineConnectors(ctx, pipelineID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListPipelineConnectors", reflect.TypeOf((*MockListConnectorsClient)(nil).ListPipelineConnectors), ctx, pipelineID)
}
