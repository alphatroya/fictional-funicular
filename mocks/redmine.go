// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/alphatroya/redmine-helper-bot/redmine (interfaces: Client)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	redmine "github.com/alphatroya/redmine-helper-bot/redmine"
	gomock "github.com/golang/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Activities mocks base method.
func (m *MockClient) Activities() ([]*redmine.Activities, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Activities")
	ret0, _ := ret[0].([]*redmine.Activities)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Activities indicates an expected call of Activities.
func (mr *MockClientMockRecorder) Activities() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Activities", reflect.TypeOf((*MockClient)(nil).Activities))
}

// AddComment mocks base method.
func (m *MockClient) AddComment(arg0, arg1 string, arg2 int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddComment", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddComment indicates an expected call of AddComment.
func (mr *MockClientMockRecorder) AddComment(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddComment", reflect.TypeOf((*MockClient)(nil).AddComment), arg0, arg1, arg2)
}

// AssignedIssues mocks base method.
func (m *MockClient) AssignedIssues() ([]*redmine.Issue, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AssignedIssues")
	ret0, _ := ret[0].([]*redmine.Issue)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AssignedIssues indicates an expected call of AssignedIssues.
func (mr *MockClientMockRecorder) AssignedIssues() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AssignedIssues", reflect.TypeOf((*MockClient)(nil).AssignedIssues))
}

// FillHoursRequest mocks base method.
func (m *MockClient) FillHoursRequest(arg0, arg1, arg2, arg3 string) (*redmine.TimeEntryBodyResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FillHoursRequest", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*redmine.TimeEntryBodyResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FillHoursRequest indicates an expected call of FillHoursRequest.
func (mr *MockClientMockRecorder) FillHoursRequest(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FillHoursRequest", reflect.TypeOf((*MockClient)(nil).FillHoursRequest), arg0, arg1, arg2, arg3)
}

// Issue mocks base method.
func (m *MockClient) Issue(arg0 string) (*redmine.IssueContainer, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Issue", arg0)
	ret0, _ := ret[0].(*redmine.IssueContainer)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Issue indicates an expected call of Issue.
func (mr *MockClientMockRecorder) Issue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Issue", reflect.TypeOf((*MockClient)(nil).Issue), arg0)
}

// TodayTimeEntries mocks base method.
func (m *MockClient) TodayTimeEntries() ([]*redmine.TimeEntryResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TodayTimeEntries")
	ret0, _ := ret[0].([]*redmine.TimeEntryResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TodayTimeEntries indicates an expected call of TodayTimeEntries.
func (mr *MockClientMockRecorder) TodayTimeEntries() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TodayTimeEntries", reflect.TypeOf((*MockClient)(nil).TodayTimeEntries))
}
